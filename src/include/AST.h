#ifndef AST_H
#define AST_H
#include <stdio.h>

typedef struct AST_STRUCT
{
    enum {
        AST_STRING,
        AST_INT,
        AST_PUSH,
        AST_COMPOUND,
        AST_NOOP,
    } type;

    struct STACK_STRUCT* stack;

    //AST_PUSH
    struct AST_STRUCT* push_value;

    // AST_STRING
    char* string_value;

    // AST_INT
    long int int_value;

    // AST_COMPOUND
    struct AST_STRUCT** compound_value;
    size_t compound_size;
} AST_T;

AST_T* init_ast(int type);
#endif