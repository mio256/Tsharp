#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "include/lex.h"
#include "include/parse.h"
#include "include/scope.h"
#include "include/io.h"
#include "include/visitor.h"

int main(int argc, char* argv[]) {
    if (argc == 2) {
        lexer_T* lexer = init_lexer(get_file_contents(argv[1]));
        parser_T* parser = init_parser(lexer);
        AST_T* node = parser_parse(parser, parser->scope);
        visitor_T* visitor = init_visitor();
        visitor_visit(visitor, node);

        AST_T* fdef = scope_find_func(parser->scope, "main");
    }
    else
    {
        printf("Usage:\n");
        printf("    tsh <filename>.t#\n");
        exit(1);
    }

    return 0;
}
