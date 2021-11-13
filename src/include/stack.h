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

#endif