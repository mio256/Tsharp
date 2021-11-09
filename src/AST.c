#include "include/AST.h"
#include <stdlib.h>

AST_T* init_ast(int type)
{
    AST_T* ast = calloc(1, sizeof(struct AST_STRUCT));
    ast->type = type;

    ast->scope = (void*) 0;

    // AST_VARIABLE_DEFINITION
    ast->variable_definition_v_name = (void*) 0;
    ast->variable_definition_f_name = (void*) 0;
    ast->variable_definition_value = (void*) 0;

    // AST_FUNCTION_DEFINITION
    ast->function_definition_name = (void*) 0;
    ast->function_definition_body = (void*) 0;
    ast->function_definition_args = (void*) 0;
    ast->function_definition_args_size = 0;
    
    // AST_VARIABLE
    ast->variable_name = (void*) 0;
    ast->variable_f_name = (void*) 0;

    // AST_FUNCTION_CALL
    ast->function_call_name = (void*) 0;
    ast->function_call_args = (void*) 0;
    ast->function_call_args_size = 0;

    // AST_IF
    ast->if_op = (void*) 0;
    ast->if_body = (void*) 0;
    ast->else_body = (void*) 0;

    // AST_STRING
    ast->string_value = (void*) 0;

    // AST_INT
    ast->int_value = (void*) 0;

    // AST_COMPOUND
    ast->compound_value = (void*) 0;
    ast->compound_size = 0;

    return ast;
}