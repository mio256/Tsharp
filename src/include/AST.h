#ifndef AST_H
#define AST_H
#include <stdio.h>

typedef struct AST_STRUCT
{
    enum {
        AST_VARIABLE_DEFINITION,
        AST_FUNCTION_DEFINITION,
        AST_VARIABLE,
        AST_FUNCTION_CALL,
        AST_COMPARE,
        AST_IF,
        AST_STRING,
        AST_INT,
        AST_COMPOUND,
        AST_NOOP,
    } type;

    struct SCOPE_STRUCT* scope;

    // AST_VARIABLE_DEFINITION
    char* variable_definition_v_name;
    char* variable_definition_f_name;
    struct AST_STRUCT* variable_definition_value;

    // AST_FUNCTION_DEFINITION
    char* function_definition_name;
    struct AST_STRUCT* function_definition_body;
    struct AST_STRUCT** function_definition_args;
    size_t function_definition_args_size;
    
    // AST_VARIABLE
    char* variable_name;
    char* variable_f_name;

    // AST_FUNCTION_CALL
    char* function_call_name;
    struct AST_STRUCT** function_call_args;
    size_t function_call_args_size;

    // AST_IF
    struct AST_STRUCT* if_op;
    struct AST_STRUCT* if_body;
    struct AST_STRUCT* else_body;

    // AST_STRING
    char* string_value;

    // AST_INT
    char* int_value;

    // AST_COMPOUND
    struct AST_STRUCT** compound_value;
    size_t compound_size;
} AST_T;

AST_T* init_ast(int type);
#endif