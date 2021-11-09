#include "include/scope.h"
#include <string.h>
#include <stdio.h>
#include <stdlib.h>

scope_T* init_scope()
{
    scope_T* scope = calloc(1, sizeof(struct SCOPE_STRUCT));
    
    scope->variable_definitions = (void*) 0;
    scope->variable_definitions_size = 0;

    scope->func_definitions = (void*) 0;
    scope->func_definitions_size = 0;

    return scope;
}

AST_T* scope_add_func_definition(scope_T* scope, AST_T* fdef)
{
    scope->func_definitions_size += 1;
    if (scope->func_definitions == (void*) 0) {
        scope->func_definitions = calloc(1, sizeof(struct AST_STRUCT*));
    } else {
        scope->func_definitions = realloc(scope->func_definitions, scope->func_definitions_size * sizeof(struct AST_STRUCT*)); 
    }
    scope->func_definitions[scope->func_definitions_size-1] = fdef;
    return fdef;
}

AST_T* scope_get_func_definition(scope_T* scope, const char* fname)
{
    for (int i = 0; i < scope->func_definitions_size; i++) {
        AST_T* fdef = scope->func_definitions[i];
        if (strcmp(fdef->function_definition_name, fname) == 0) {
            return fdef;
        }
    }
    return (void*)0;
}

AST_T* scope_find_func(scope_T* scope, const char* name)
{
    for (int i = 0; i < scope->func_definitions_size; i++)
    {
        AST_T* ldef = scope->func_definitions[i];
        if (strcmp(ldef->function_definition_name, name) == 0)
        {
            return ldef;
        }
    }

    if (strcmp(name, "main") == 0) {
        printf("Error: function main is undeclared\n");
        exit(1);
    }

    printf("Error: undifined function '%s'\n", name);
    exit(1);
}