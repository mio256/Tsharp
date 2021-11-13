#include "include/lex.h"
#include <stdlib.h>
#include <string.h>
#include <ctype.h>
#include <stdio.h>

lexer_T* init_lexer(char* contents)
{
    lexer_T* lexer = calloc(1, sizeof(struct LEXER_STRUCT));
    lexer->contents = contents;
    lexer->i = 0;
    lexer->c = contents[lexer->i];
    lexer->line_n = 1;
    return lexer;
}

void lexer_advance(lexer_T* lexer)
{
    if (lexer->c != '\0' && lexer->i < strlen(lexer->contents))
    {
        lexer->i += 1;
        lexer->c = lexer->contents[lexer->i];
    }
    
    if (lexer->c == '\n' || lexer->c == 10)
        lexer->line_n += 1;
}

void lexer_skip_whitespace(lexer_T* lexer)
{
    while (lexer->c == ' ' || lexer->c == 10 || lexer->c == 13 || lexer->c == '	')
        lexer_advance(lexer);
}

void lexer_skip_inline_comment(lexer_T* lexer)
{
    while (lexer->c != '\n' && lexer->c != 10 && lexer->c != '\0')
        lexer_advance(lexer);
}

void lexer_skip_block_comment(lexer_T* lexer)
{
    while (1)
    {
        lexer_advance(lexer);
        if (lexer->c == '*')
        {
            lexer_advance(lexer);
            if (lexer->c == '/')
            {
                lexer_advance(lexer);
                return;
            }
        }
    }
}

token_T* lexer_get_next_token(lexer_T* lexer)
{
    while (lexer->c != '\0' && lexer->i < strlen(lexer->contents))
    {
        if (lexer->c == ' ' || lexer->c == 10 || lexer->c == 13 || lexer->c == '	')
            lexer_skip_whitespace(lexer);

        if (isalpha(lexer->c))
            return lexer_collect_id(lexer);

        if (lexer->c == '"')
            return lexer_collect_string(lexer);

        if (lexer->c == '\'')
            return lexer_collect_string_single_quot(lexer);

        if (isdigit(lexer->c))
            return lexer_collect_int(lexer);

        if (lexer->c == '/')
        {
            char* value = lexer_get_current_char_as_string(lexer);
            lexer_advance(lexer);
            if (lexer->c == '/')
            {
                lexer_advance(lexer);
                lexer_skip_inline_comment(lexer);
                continue;
            }
            if (lexer->c == '*')
            {
                lexer_advance(lexer);
                lexer_skip_block_comment(lexer);
                continue;
            }
        }

        switch (lexer->c)
        {
            case ':': return lexer_advance_token(lexer, init_token(TOKEN_COLON, lexer_get_current_char_as_string(lexer))); break;
            case ';': return lexer_advance_token(lexer, init_token(TOKEN_SEMI, lexer_get_current_char_as_string(lexer))); break;
            case '\0': return init_token(TOKEN_EOF, "\0"); break;
            default: printf("SyntaxError: unexpected '%c' (line %d)\n", lexer->c, lexer->line_n); exit(1); break;
        }
    }

    return init_token(TOKEN_EOF, "\0");
}

token_T* lexer_collect_string(lexer_T* lexer)
{
    lexer_advance(lexer);
    char* value = calloc(1, sizeof(char));
    value[0] = '\0';
    while (lexer->c != '"')
    {
        char* s = lexer_get_current_char_as_string(lexer);
        value = realloc(value, (strlen(value) + strlen(s) + 1) * sizeof(char));
        strcat(value, s);
        lexer_advance(lexer);
    }
    lexer_advance(lexer);
    return init_token(TOKEN_STRING, value);
}

token_T* lexer_collect_string_single_quot(lexer_T* lexer)
{
    lexer_advance(lexer);
    char* value = calloc(1, sizeof(char));
    value[0] = '\0';
    while (lexer->c != '\'')
    {
        char* s = lexer_get_current_char_as_string(lexer);
        value = realloc(value, (strlen(value) + strlen(s) + 1) * sizeof(char));
        strcat(value, s);
        lexer_advance(lexer);
    }
    lexer_advance(lexer);
    return init_token(TOKEN_STRING, value);
}

token_T* lexer_collect_int(lexer_T* lexer)
{
    char* value = calloc(1, sizeof(char));
    value[0] = '\0';
    while(isdigit(lexer->c))
    {
        if (isdigit(lexer->c))
        {
            char* s = lexer_get_current_char_as_string(lexer);
            value = realloc(value, (strlen(value) + strlen(s) + 1) * sizeof(char));
            strcat(value, s);

            lexer_advance(lexer);
        }
        else
        {
            printf("Error: expecting integer '%c' (line %d)\n", lexer->c, lexer->line_n);
            exit(1);
        }
    }
    return init_token(TOKEN_INT, value);
}

token_T* lexer_collect_id(lexer_T* lexer)
{
    char* value = calloc(1, sizeof(char));
    value[0] = '\0';
    while (isalnum(lexer->c) || lexer->c == '_')
    {
        char* s = lexer_get_current_char_as_string(lexer);
        value = realloc(value, (strlen(value) + strlen(s) + 1) * sizeof(char));
        strcat(value, s);
        lexer_advance(lexer);
    }

    if (strcmp(value, "end") == 0)
        return init_token(TOKEN_END, value);

    if (strcmp(value, "do") == 0)
        return init_token(TOKEN_DO, value);

    return init_token(TOKEN_ID, value);
}

token_T* lexer_advance_token(lexer_T* lexer, token_T* token)
{
    lexer_advance(lexer);
    return token;
}

char* lexer_get_current_char_as_string(lexer_T* lexer)
{
    char* str = calloc(2, sizeof(char));
    str[0] = lexer->c;
    str[1] = '\0';
    return str;
}
