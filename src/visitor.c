#include "include/visitor.h"
#include "include/scope.h"
#include "include/AST.h"
#include "include/token.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <ctype.h>
#include <sys/stat.h>

int write_to_file(char* contents) {
    FILE *fp; fp = fopen("tshc/tshc.c", "a+");
    if (fp) { fprintf(fp, "%s", contents); fclose(fp); return 0; }
    return 0;
}

visitor_T* init_visitor()
{
    visitor_T* visitor = calloc(1, sizeof(struct VISITOR_STRUCT));
    mkdir("tshc", S_IRWXU);
    FILE *fp; fp = fopen("tshc/tshc.c", "w"); fclose(fp);
    write_to_file(
        "// Generated C Programming Language\n"
        "#include <stdio.h>\n"
        "#include <stdlib.h>\n"
    );
    return visitor;
}

AST_T* visitor_visit(visitor_T* visitor, AST_T* node, char* op)
{
    switch (node->type)
    {
        case AST_FUNCTION_DEFINITION: return visitor_visit_function_definition(visitor, node); break;
        case AST_VARIABLE_DEFINITION: return visitor_visit_variable_definition(visitor, node); break;
        case AST_FUNCTION_CALL: return visitor_visit_function_call(visitor, node, op); break;
        case AST_IF: return visitor_visit_if(visitor, node, op); break;
        case AST_VARIABLE: return visitor_visit_variable(visitor, node, op); break;
        case AST_STRING: return visitor_visit_string(visitor, node, op); break;
        case AST_INT: return visitor_visit_int(visitor, node, op); break;
        case AST_COMPOUND: return visitor_visit_compound(visitor, node); break;
        case AST_NOOP: return node; break;
        default: printf("Error: uncaught statement of type '%d'\n", node->type); break;
    }
    exit(1);

    return init_ast(AST_NOOP);
}

static AST_T* builtin_function_print(visitor_T* visitor, AST_T** args, int args_size)
{
    write_to_file("    printf(");
    if (args_size != 1) {
        printf("Error: function print() expecting most and least one argument");
        exit(1);
    }
    AST_T* visited_ast = visitor_visit(visitor, args[0], "print");
    write_to_file(");\n");

    return init_ast(AST_NOOP);
}

AST_T* visitor_visit_function_call(visitor_T* visitor, AST_T* node, char* op)
{
    if (strcmp(node->function_call_name, "print") == 0)
        return builtin_function_print(visitor, node->function_call_args, node->function_call_args_size);


    if (strcmp(node->function_call_name, "exit") == 0) {
        if (strcmp(op, "print") == 0) {printf("Error: function exit() has no return value\n");exit(1);}
        if (node->function_call_args_size != 0) {printf("Error: function exit() expected at most 0 argument got %zu\n", node->function_call_args_size);exit(1);}
        write_to_file("exit(1);\n");
        return node;
    }

    AST_T* fdef = scope_get_func_definition(node->scope, node->function_call_name);

    if (fdef == (void*) 0) {printf("Error: undefined function '%s'\n", node->function_call_name);exit(1);}

    visitor_visit(visitor, fdef->function_definition_body, "");

    return node;
}

AST_T* visitor_visit_variable_definition(visitor_T* visitor, AST_T* node)
{
    AST_T* vdef = scope_get_variable_definition(node->scope, node->variable_definition_v_name, node->variable_definition_f_name);
    AST_T* visited_ast = visitor_visit(visitor, node->variable_definition_value, "");
    if (vdef == (void*) 0) {
        scope_add_variable_definition(node->scope, node);
        if (visited_ast->type == AST_STRING) {
            write_to_file("    char* ");
        } else if (visited_ast->type == AST_INT) {
            write_to_file("    int ");
        }
    }
    else
    {
        write_to_file("    ");
        AST_T* visited_scope_value = visitor_visit(visitor, vdef->variable_definition_value, "");
        if (visited_scope_value->type != visited_ast->type) {
            printf("TypeError: '%s' is type ", node->variable_definition_v_name);
            switch (visited_scope_value->type)
            {
                case AST_STRING: printf("<string>"); break;
                case AST_INT: printf("<int>"); break;
                default: printf("<undefined type>"); break;
            }
            printf(" cannot use ");
            switch (visited_ast->type)
            {
                case AST_STRING: printf("'%s'\n", visited_ast->string_value); break;
                case AST_INT: printf("'%s'\n", visited_ast->int_value); break;
                default: printf("<undefined type>"); break;
            }
            exit(1);
        }
    }

    write_to_file(node->variable_definition_v_name);
    write_to_file(" = ");
    visitor_visit(visitor, node->variable_definition_value, "print");
    write_to_file(";\n");

    return node;
}

AST_T* visitor_visit_function_definition(visitor_T* visitor, AST_T* node)
{
    scope_add_func_definition(node->scope, node);

    write_to_file("\nint ");
    write_to_file(node->function_definition_name);
    write_to_file("() {\n");

    visitor_visit(visitor, node->function_definition_body, "");

    write_to_file("}\n");

    return node;
}

AST_T* visitor_visit_if(visitor_T* visitor, AST_T* node, char* op)
{
    write_to_file("    if (");
    visitor_visit(visitor, node->if_op, "print");
    write_to_file(") {\n");
    visitor_visit(visitor, node->if_body, "");
    write_to_file("    }");
    if (node->else_body != (void*) 0) {
        write_to_file(" else {\n");
        visitor_visit(visitor, node->else_body, "");
        write_to_file("    }\n");
    } else {
        write_to_file("\n");
    }
    return node;
}

AST_T* visitor_visit_string(visitor_T* visitor, AST_T* node, char* op)
{
    if (strcmp(op, "print") == 0)
    {
        write_to_file("\"");
        write_to_file(node->string_value);
        write_to_file("\"");
    }
    return node;
}

AST_T* visitor_visit_int(visitor_T* visitor, AST_T* node, char* op)
{
    if (strcmp(op, "print") == 0)
        write_to_file(node->int_value);

    return node;
}

AST_T* visitor_visit_variable(visitor_T* visitor, AST_T* node, char* op)
{
    AST_T* vdef = scope_get_variable_definition(node->scope, node->variable_name, node->variable_f_name);

    if (vdef != (void*) 0) {
        AST_T* visited_ast = visitor_visit(visitor, vdef->variable_definition_value, "");
        if (strcmp(op, "print") == 0) {
            if (visited_ast->type == AST_STRING)
            {
                write_to_file("\"%s\", ");
            }
            else
            if (visited_ast->type == AST_INT)
            {
                write_to_file("\"%d\", ");
            }
            write_to_file(node->variable_name);
        }
        if (strcmp(op, "compare") == 0) {
            write_to_file(node->variable_name);
        }
        return visitor_visit(visitor, vdef->variable_definition_value, "");
    }

    printf("Error: undefined variable '%s'\n", node->variable_name);
    exit(1);
}

AST_T* visitor_visit_compound(visitor_T* visitor, AST_T* node)
{
    for (int i = 0; i < node->compound_size; i++)
        visitor_visit(visitor, node->compound_value[i], "");

    return init_ast(AST_NOOP);
}