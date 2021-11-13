#ifndef STACK_H
#define STACK_H
#include "AST.h"

typedef struct STACK_STRUCT
{
    AST_T** stack;
    size_t stack_size;
} stack_T;

stack_T* init_stack();

AST_T* stack_push_value(stack_T* stack, AST_T* pushv);

AST_T* stack_get_first_value(stack_T* stack);

AST_T* stack_drop_first_value(stack_T* stack);

AST_T* stack_free(stack_T* stack);

#endif