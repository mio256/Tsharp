#include "include/parse.h"
#include "include/stack.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

parser_T* init_parser(lexer_T* lexer)
{
    parser_T* parser = calloc(1, sizeof(struct PARSER_STRUCT));
    parser->lexer = lexer;
    parser->current_token = lexer_get_next_token(lexer);
    token_T* prev_token = parser->current_token;
    parser->stack = init_stack();
    return parser;
}

void parser_eat(parser_T* parser, int token_type)
{
    if (parser->current_token->type == token_type)
    {
        parser->prev_token = parser->current_token;
        parser->current_token = lexer_get_next_token(parser->lexer);
    }
    else
    {
        printf(
            "SyntaxError: unexpected token '%s' (line %d)\n",
            parser->current_token->value,
            parser->lexer->line_n
        );
        exit(1);
    }
}

AST_T* parser_parse(parser_T* parser, stack_T* stack)
{
    return parser_parse_statements(parser, stack);
}

AST_T* parser_parse_statement(parser_T* parser, stack_T* stack)
{
    switch (parser->current_token->type)
    {
        case TOKEN_ID: return parser_parse_id(parser, stack);
        default: return 0;
    }

    return init_ast(AST_NOOP);
}

AST_T* parser_parse_statements(parser_T* parser, stack_T* stack)
{
    AST_T* compound = init_ast(AST_COMPOUND);
    compound->compound_value = calloc(1, sizeof(struct AST_STRUCT*));
    AST_T* ast_statement = parser_parse_statement(parser, stack);
    compound->compound_value[0] = ast_statement;
    compound->compound_size += 1;
    while(parser->current_token->type == TOKEN_COLON)
    {
        if (parser->current_token->type == TOKEN_COLON)
            parser_eat(parser, TOKEN_COLON);
        
        AST_T* ast_statement = parser_parse_statement(parser, stack);
        if (ast_statement)
        {
            compound->compound_size += 1;
            compound->compound_value = realloc(
                compound->compound_value,
                compound->compound_size * sizeof(struct AST_STRUCT)
            );
            compound->compound_value[compound->compound_size-1] = ast_statement;
        }
    }
    compound->stack = stack;
    return compound;
}

AST_T* parser_parse_expr(parser_T* parser, stack_T* stack)
{
    switch (parser->current_token->type)
    {
        case TOKEN_STRING: return parser_parse_string(parser, stack);
        case TOKEN_INT: return parser_parse_int(parser, stack);
        case TOKEN_ID: return parser_parse_id(parser, stack);
        default: return 0;
    }
    return init_ast(AST_NOOP);
}

AST_T* parser_parse_push(parser_T* parser, stack_T* stack)
{
    AST_T* ast = init_ast(AST_PUSH);
    parser_eat(parser, TOKEN_ID);
    parser_eat(parser, TOKEN_COLON);
    parser_eat(parser, TOKEN_COLON);
    ast->push_value = parser_parse_expr(parser, stack);

    ast->stack = stack;
    return ast;
}

AST_T* parser_parse_print(parser_T* parser, stack_T* stack)
{
    AST_T* ast = init_ast(AST_PRINT);
    parser_eat(parser, TOKEN_ID);
    ast->stack = stack;
    return ast;
}

AST_T* parser_parse_dup(parser_T* parser, stack_T* stack)
{
    AST_T* ast = init_ast(AST_DUP);
    parser_eat(parser, TOKEN_ID);
    ast->stack = stack;
    return ast;
}

AST_T* parser_parse_string(parser_T* parser, stack_T* stack)
{
    AST_T* ast_string = init_ast(AST_STRING);
    ast_string->string_value = parser->current_token->value;
    parser_eat(parser, TOKEN_STRING);
    ast_string->stack = stack;
    return ast_string;
}

AST_T* parser_parse_int(parser_T* parser, stack_T* stack)
{
    char* endPtr;
    long int int_value = strtol(parser->current_token->value, &endPtr, 10);
    parser_eat(parser, TOKEN_INT);
    AST_T* ast_int = init_ast(AST_INT);
    ast_int->int_value = int_value;
    ast_int->stack = stack;
    return ast_int;
}

AST_T* parser_parse_id(parser_T* parser, stack_T* stack)
{
    if (strcmp(parser->current_token->value, "push") == 0)
        return parser_parse_push(parser, stack);

    if (strcmp(parser->current_token->value, "print") == 0)
        return parser_parse_print(parser, stack);

    if (strcmp(parser->current_token->value, "dup") == 0)
        return parser_parse_dup(parser, stack);

    printf("SyntaxError: unkown name '%s'\n", parser->current_token->value);
    exit(1);
}
