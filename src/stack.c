#include "include/stack.h"
#include <string.h>
#include <stdio.h>
#include <stdlib.h>

stack_T* init_stack()
{
    stack_T* stack = calloc(1, sizeof(struct STACK_STRUCT));

    stack->stack = (void*) 0;
    stack->stack_size = 0;

    return stack;
}

AST_T* stack_push_value(stack_T* stack, AST_T* pushv)
{
    stack->stack_size += 1;
    if (stack->stack == (void*) 0) {
        stack->stack = calloc(1, sizeof(struct AST_STACK*));
    } else {
        stack->stack = realloc(stack->stack, stack->stack_size * sizeof(struct AST_STACK*)); 
    }
    stack->stack[stack->stack_size-1] = pushv;
    return pushv;
}
