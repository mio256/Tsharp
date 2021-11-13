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

AST_T* stack_get_first_value(stack_T* stack)
{
    if (stack->stack_size == 0) {
        printf("PrintError: stack is empty [ ]\n");
        exit(1);
    }
    AST_T* stackv = stack->stack[stack->stack_size-1];
    return stackv;
}

AST_T* stack_drop_first_value(stack_T* stack)
{
    if (stack->stack_size == 0) {
        printf("DropError: stack is empty [ ]\n");
        exit(1);
    }
    AST_T* ast = stack->stack[stack->stack_size-1];
    ast_free(ast);
    //free(stack->stack[stack->stack_size-1]);
    stack->stack_size -= 1;
    return (void*) 0;
}
