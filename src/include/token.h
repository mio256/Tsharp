#ifndef TOKEN_H
#define TOKEN_H

typedef struct TOKEN_STRUCT
{
    enum
    {
        TOKEN_STRING,
        TOKEN_ID,
        TOKEN_SEMI,
        TOKEN_LPAREN,
        TOKEN_RPAREN,
        TOKEN_EQUALS,
        TOKEN_COMMA,
        TOKEN_DOT,
        TOKEN_END,
        TOKEN_DO,
        TOKEN_ELSE,
        TOKEN_INT,
        TOKEN_EOF
    } type;

    char* value;
} token_T;

token_T* init_token(int type, char* value);

#endif