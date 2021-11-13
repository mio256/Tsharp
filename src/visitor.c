#include "include/visitor.h"
#include "include/stack.h"
#include "include/AST.h"
#include "include/token.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

visitor_T* init_visitor()
{
    visitor_T* visitor = calloc(1, sizeof(struct VISITOR_STRUCT));
    return visitor;
}

AST_T* visitor_visit(visitor_T* visitor, AST_T* node)
{
    switch (node->type)
    {
        case AST_PRINT: return visitor_visit_print(visitor, node); break;
        case AST_PUSH: return visitor_visit_push(visitor, node); break;
        case AST_STRING: return visitor_visit_string(visitor, node); break;
        case AST_INT: return visitor_visit_int(visitor, node); break;
        case AST_COMPOUND: return visitor_visit_compound(visitor, node); break;
        case AST_NOOP: return node; break;
    }

    printf("Error: uncaught statement of type '%d'\n", node->type);
    exit(1);

    return init_ast(AST_NOOP);
}

static AST_T* builtin_function_print(visitor_T* visitor, AST_T* arg)
{
    AST_T* visited_ast = visitor_visit(visitor, arg);

    switch (visited_ast->type) {
        case AST_STRING: printf("%s", visited_ast->string_value); break;
        case AST_INT: printf("%ld", visited_ast->int_value); break;
    }
    printf("\n");
}

AST_T* visitor_visit_push(visitor_T* visitor, AST_T* node)
{
    stack_push_value(node->stack, node->push_value);
    return node;
}

AST_T* visitor_visit_print(visitor_T* visitor, AST_T* node)
{
    AST_T* stackv = stack_get_first_value(node->stack);
    builtin_function_print(visitor, stackv);
    //stack_drop_first_value(node->stack);
    return node;
}

AST_T* visitor_visit_string(visitor_T* visitor, AST_T* node)
{
    return node;
}

AST_T* visitor_visit_int(visitor_T* visitor, AST_T* node)
{
    return node;
}

AST_T* visitor_visit_compound(visitor_T* visitor, AST_T* node)
{
    for (int i = 0; i < node->compound_size; i++)
        visitor_visit(visitor, node->compound_value[i]);

    return init_ast(AST_NOOP);
}