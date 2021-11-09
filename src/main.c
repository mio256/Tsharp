#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "include/lex.h"
#include "include/parse.h"
#include "include/scope.h"
#include "include/io.h"
#include "include/visitor.h"

int main(int argc, char* argv[]) {

    if (argc == 3) {
        lexer_T* lexer = init_lexer(get_file_contents(argv[2]));
        parser_T* parser = init_parser(lexer);
        AST_T* root = parser_parse(parser, parser->scope);
        visitor_T* visitor = init_visitor();
        visitor_visit(visitor, root, "");

        AST_T* fdef = scope_find_func(parser->scope, "main");

        char* str1 = "gcc -o "; char* str2 = argv[1]; char* new_str;
        if ((new_str = malloc(strlen(str1)+strlen(str2)+1)) != NULL) {new_str[0] = '\0';strcat(new_str,str1);strcat(new_str,str2);}
        str1 = new_str; str2 = " tshc/tshc.c";
        if ((new_str = malloc(strlen(str1)+strlen(str2)+1)) != NULL) {new_str[0] = '\0';strcat(new_str,str1);strcat(new_str,str2);}
        system(new_str); free(new_str);
    }
    else
    {
        printf("Usage:\n");
        printf("    tsh main <filename>\n");
        exit(1);
    }

    return 0;
}
