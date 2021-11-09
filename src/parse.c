#include "include/parse.h"
#include "include/scope.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

parser_T* init_parser(lexer_T* lexer)
{
    parser_T* parser = calloc(1, sizeof(struct PARSER_STRUCT));
    parser->lexer = lexer;
    parser->current_token = lexer_get_next_token(lexer);
    token_T* prev_token = parser->current_token;
    parser->scope = init_scope();
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

AST_T* parser_parse(parser_T* parser, scope_T* scope)
{
    return parser_parse_statements(parser, scope, "");
}

AST_T* parser_parse_statement(parser_T* parser, scope_T* scope, char* func_name)
{
    switch (parser->current_token->type)
    {
        case TOKEN_ID: return parser_parse_id(parser, scope, func_name);
        default: return 0;
    }

    return init_ast(AST_NOOP);
}

AST_T* parser_parse_statements(parser_T* parser, scope_T* scope, char* func_name)
{
    AST_T* compound = init_ast(AST_COMPOUND);
    compound->compound_value = calloc(1, sizeof(struct AST_STRUCT*));
    AST_T* ast_statement = parser_parse_statement(parser, scope, func_name);
    compound->compound_value[0] = ast_statement;
    compound->compound_size += 1;
    while(parser->current_token->type == TOKEN_SEMI)
    {
        if (parser->current_token->type == TOKEN_SEMI)
            parser_eat(parser, TOKEN_SEMI);
        
        AST_T* ast_statement = parser_parse_statement(parser, scope, func_name);
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
    compound->scope = scope;
    return compound;
}

AST_T* parser_parse_expr(parser_T* parser, scope_T* scope, char* func_name)
{
    switch (parser->current_token->type)
    {
        case TOKEN_STRING: return parser_parse_string(parser, scope, func_name);
        case TOKEN_INT: return parser_parse_int(parser, scope, func_name);
        case TOKEN_ID: return parser_parse_id(parser, scope, func_name);
        default: return 0;
    }
    return init_ast(AST_NOOP);
}

AST_T* parser_parse_function_call(parser_T* parser, scope_T* scope, char* func_name)
{
    AST_T* ast = init_ast(AST_FUNCTION_CALL);
    
    char* func_call_name = parser->prev_token->value;
    ast->function_call_name = calloc(
        strlen(func_call_name) + 1,
        sizeof(char)
    );
    strcpy(ast->function_call_name, func_call_name);

    parser_eat(parser, TOKEN_LPAREN);
    ast->function_call_args = calloc(1, sizeof(struct AST_STRUCT*));
    if (parser->current_token->type != TOKEN_RPAREN) {
        AST_T* ast_expr = parser_parse_expr(parser, scope, func_name);
        ast->function_call_args[0] = ast_expr;
        ast->function_call_args_size++;
    }
    while (parser->current_token->type == TOKEN_COMMA) {
        parser_eat(parser, TOKEN_COMMA);
        AST_T* ast_expr = parser_parse_expr(parser, scope, func_name);
        ast->function_call_args_size++;
        ast->function_call_args = realloc(ast->function_call_args,ast->function_call_args_size * sizeof(struct AST_STRUCT));
        ast->function_call_args[ast->function_call_args_size - 1] = ast_expr;
    }
    parser_eat(parser, TOKEN_RPAREN);

    ast->scope = scope;
    return ast;
}

AST_T* parser_parse_function_definition(parser_T* parser, scope_T* scope)
{
    AST_T* ast = init_ast(AST_FUNCTION_DEFINITION);
    parser_eat(parser, TOKEN_ID);

    char* func_name = parser->current_token->value;
    ast->function_definition_name = calloc(
        strlen(func_name) + 1,
        sizeof(char)
    );
    strcpy(ast->function_definition_name, func_name);
    parser_eat(parser, TOKEN_ID);

    parser_eat(parser, TOKEN_LPAREN);
    parser_eat(parser, TOKEN_RPAREN);
    parser_eat(parser, TOKEN_DO);

    ast->function_definition_body = parser_parse_statements(parser, scope, func_name);

    parser_eat(parser, TOKEN_END);

    ast->scope = scope;
    return ast;
}

AST_T* parser_parse_string(parser_T* parser, scope_T* scope, char* func_name)
{
    AST_T* ast_string = init_ast(AST_STRING);
    ast_string->string_value = parser->current_token->value;
    parser_eat(parser, TOKEN_STRING);
    ast_string->scope = scope;
    return ast_string;
}

AST_T* parser_parse_int(parser_T* parser, scope_T* scope, char* func_name)
{
    char* endPtr;
    long int int_value = strtol(parser->current_token->value, &endPtr, 10);
    parser_eat(parser, TOKEN_INT);
    AST_T* ast_int = init_ast(AST_INT);
    ast_int->int_value = int_value;
    ast_int->scope = scope;
    return ast_int;
}

AST_T* parser_parse_id(parser_T* parser, scope_T* scope, char* func_name)
{
    if (strcmp(parser->current_token->value, "func") == 0)
        return parser_parse_function_definition(parser, scope);

    if (strcmp(func_name, "") == 0) {
        printf("SyntaxError: non-declaration statement outside function body\n");
        exit(1);
    }
    parser_eat(parser, TOKEN_ID);
    return parser_parse_function_call(parser, scope, func_name);
}
