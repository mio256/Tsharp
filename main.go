package main

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
	"os"
	"strconv"
)


// -----------------------------
// ----------- Lexer -----------
// -----------------------------

type Token int
const (
	TOKEN_EOF = iota
	TOKEN_ILLEGAL
	TOKEN_ID
	TOKEN_STRING
	TOKEN_INT
	TOKEN_PLUS
	TOKEN_MINUS
	TOKEN_END
	TOKEN_DO
)

var tokens = []string{
	TOKEN_EOF:  "TOKEN_EOF",
	TOKEN_ILLEGAL: "TOKEN_ILLEGAL",
	TOKEN_ID:   "TOKEN_ID",
	TOKEN_STRING: "TOKEN_STRING",
	TOKEN_INT:  "TOKEN_INT",
	TOKEN_PLUS: "TOKEN_PLUS",
	TOKEN_MINUS: "TOKEN_MINUS",
	TOKEN_END: "TOKEN_END",
	TOKEN_DO: "TOKEN_DO",
}

func (token Token) String() string {
	return tokens[token]
}

type Position struct {
	line int
	column int
}

type Lexer struct {
	pos Position
	reader *bufio.Reader
}

func LexerInit(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position {line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}

func (lexer *Lexer) Lex() (Position, Token, string) {
	for {
		r, _, err := lexer.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				err = nil
				return lexer.pos, TOKEN_EOF, "EOF"
			}
			panic(err)
		}
		lexer.pos.column++
		switch r {
			case '\n': lexer.resetPosition()
			case '+': return lexer.pos, TOKEN_PLUS, "+"
			case '-': return lexer.pos, TOKEN_MINUS, "-"
			default:
				if unicode.IsSpace(r) {
					continue
				} else if unicode.IsDigit(r) {
					startPos := lexer.pos
					lexer.backup()
					val := lexer.lexInt()
					return startPos, TOKEN_INT, val
				} else if unicode.IsLetter(r) {
					startPos := lexer.pos
					lexer.backup()
					val := lexer.lexId()
					if val == "end" {
						return startPos, TOKEN_END, val
					} else if val == "do" {
						return startPos, TOKEN_DO, val
					}
					return startPos, TOKEN_ID, val
				} else if r == '"' {
					startPos := lexer.pos
					lexer.backup()
					val := lexer.lexString()
					r, _, err = lexer.reader.ReadRune()
					return startPos, TOKEN_STRING, val
				}
        }
	}
}

func (lexer *Lexer) backup() {
	if err := lexer.reader.UnreadRune(); err != nil {
		panic(err)
	}
	lexer.pos.column--
}

func (lexer *Lexer) lexId() string {
	var val string
	for {
		r, _, err := lexer.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return val
			}
		}
        lexer.pos.column++
		if unicode.IsLetter(r) {
			val = val + string(r)
		} else {
			lexer.backup()
			return val
		}
	}
}

func (lexer *Lexer) lexInt() string {
	var val string
	for {
		r, _, err := lexer.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return val
			}
		}
		lexer.pos.column++
		if unicode.IsDigit(r) {
			val = val + string(r)
		} else {
			lexer.backup()
			return val
		}
	}
}

func (lexer *Lexer) lexString() string {
	var val string
	r, _, err := lexer.reader.ReadRune()
	for {
		r, _, err = lexer.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return val
			}
		}
		lexer.pos.column++
		if r != '"' {
			val = val + string(r)
		} else {
			lexer.backup()
			return val
		}
	}
}

func (lexer *Lexer) resetPosition() {
	lexer.pos.line++
	lexer.pos.column = 0
}


// -----------------------------
// ------------ AST ------------
// -----------------------------

type ExprType int
const (
	ExprVoid ExprType = iota
	ExprInt
	ExprStr
	ExprPush
	ExprBlockdef
	ExprPrint
	ExprCall
	ExprPlus
	ExprMinus
)

type Expr struct {
	Type ExprType
	AsInt int
	AsStr string
	AsPush *Push
	AsBlockdef *Blockdef
	AsCall *Call
}

type Push struct {
	Arg Expr
}

type Call struct {
	Value string
}

type Blockdef struct {
	Name string
	Body []Expr
}


// -----------------------------
// ----------- Parse -----------
// -----------------------------

type Parser struct {
	current_token_type Token
	current_token_value string
	lexer Lexer
}

func ParserInit(lexer *Lexer) *Parser {
	_, tok, val := lexer.Lex()
	return &Parser{
		current_token_type: tok,
		current_token_value: val,
		lexer: *lexer,
	}
}

func (parser *Parser) ParserEat(token Token) {
	if token != parser.current_token_type {
		fmt.Println("Error: unexpected token value '" + parser.current_token_value + "'")
		os.Exit(0)
	}
	_, tok, val := parser.lexer.Lex()
	parser.current_token_type = tok
	parser.current_token_value = val
}

func StrToInt(num string) int {
	i, err := strconv.Atoi(num)
	if err != nil {
		panic(err)
	}
	return i
}

func ParserParseExpr(parser *Parser) (Expr) {
	expr := Expr{}
	switch parser.current_token_type {
	case TOKEN_INT:
		expr.Type = ExprInt
		expr.AsInt = StrToInt(parser.current_token_value)
		parser.ParserEat(TOKEN_INT)
	case TOKEN_STRING:
		expr.Type = ExprStr
		expr.AsStr = parser.current_token_value
		parser.ParserEat(TOKEN_STRING)
	default:
		fmt.Println("Error: unexpected token value '" + parser.current_token_value + "'")
		os.Exit(0)
	}
	return expr
}

func ParserParse(parser *Parser)  ([]Expr, Parser) {
	exprs := []Expr{}

	for parser.current_token_type != TOKEN_EOF || parser.current_token_type != TOKEN_END {
		expr := Expr{}
		if parser.current_token_type == TOKEN_ID {
			if parser.current_token_value == "push" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprPush
				expr.AsPush = &Push{
					Arg: ParserParseExpr(parser),
				}
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "print" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprPrint
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "block" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprBlockdef
				if parser.current_token_type != TOKEN_ID {
					fmt.Println("SyntaxError: unexpected token value '" + parser.current_token_value + "'")
					os.Exit(0)
				}
				name := parser.current_token_value
				parser.ParserEat(TOKEN_ID)
				parser.ParserEat(TOKEN_DO)
				body, _ := ParserParse(parser)
				expr.AsBlockdef = &Blockdef{
					Name: name,
					Body: body,
				}
				parser.ParserEat(TOKEN_END)
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "call" {
				parser.ParserEat(TOKEN_ID)
				if parser.current_token_type != TOKEN_ID {
					fmt.Println("Error: unexpected token value '" + parser.current_token_value + "'")
					os.Exit(0)
				}
				expr.Type = ExprCall
				expr.AsCall = &Call{
					Value: parser.current_token_value,
				}
				parser.ParserEat(TOKEN_ID)
				exprs = append(exprs, expr)
			} else {
				fmt.Println("Error: unexpected token value '" + parser.current_token_value + "'")
				os.Exit(0)
			}
		} else if parser.current_token_type == TOKEN_PLUS {
			expr.Type = ExprPlus
			parser.ParserEat(TOKEN_PLUS)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_MINUS {
			expr.Type = ExprMinus
			parser.ParserEat(TOKEN_MINUS)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_END {
			return exprs, *parser
		} else if parser.current_token_type == TOKEN_EOF {
			return exprs, *parser
		} else {
			fmt.Println("SyntaxError: unexpected token value '" + parser.current_token_value + "'")
			os.Exit(0)
		}
	}
	
	return exprs, *parser
}


// -----------------------------
// ----------- Stack -----------
// -----------------------------

type StackItem struct {
	string_value string
	int_value int
}

type Stack struct {
	Values []StackItem
}

var theStack Stack = Stack{}

func (stack *Stack) OpPush(item StackItem) {
	stack.Values = append(stack.Values, item)
}

func (stack *Stack) OpPrint() {
	if len(stack.Values)-1 < 0 {
		fmt.Println("PrintError: the stack is empty")
		os.Exit(3)
	}
	if (len(stack.Values[len(stack.Values)-1].string_value) == 0) {
		fmt.Println(stack.Values[len(stack.Values)-1].int_value)
	} else {
		fmt.Print(stack.Values[len(stack.Values)-1].string_value)
	}
	stack.Values = stack.Values[:len(stack.Values)-1]
}

func (stack *Stack) OpPlus() {
	a := stack.Values[len(stack.Values)-1].int_value
	b := stack.Values[len(stack.Values)-2].int_value
	x := a + b
	stack.Values = stack.Values[:len(stack.Values)-1]
	stack.Values = stack.Values[:len(stack.Values)-1]
	theStack.OpPush(StackItem{int_value: x})
}

func (stack *Stack) OpMinus() {
	a := stack.Values[len(stack.Values)-1].int_value
	b := stack.Values[len(stack.Values)-2].int_value
	x := b - a
	stack.Values = stack.Values[:len(stack.Values)-1]
	stack.Values = stack.Values[:len(stack.Values)-1]
	theStack.OpPush(StackItem{int_value: x})
}


// -----------------------------
// ------ Visitor / Scope ------
// -----------------------------

//BlockScope := make(map[string][]Expr)

var BlockScope = map[string][]Expr{} // Block Scope

func VisitExpr(exprs []Expr) {
	for _, expr := range exprs {
		switch expr.Type {
		case ExprPush:
			if expr.AsPush.Arg.Type == ExprInt {
				theStack.OpPush(StackItem{int_value: expr.AsPush.Arg.AsInt})
			} else if expr.AsPush.Arg.Type == ExprStr {
				theStack.OpPush(StackItem{string_value: expr.AsPush.Arg.AsStr})
			}
		case ExprPrint:
			theStack.OpPrint()
		case ExprBlockdef:
			if _, ok := BlockScope[expr.AsBlockdef.Name]; ok {
				fmt.Println("Error: we can't define blocks that are the same name")
				os.Exit(0)
			}
			BlockScope[expr.AsBlockdef.Name] = expr.AsBlockdef.Body
		case ExprCall:
			if _, ok := BlockScope[expr.AsCall.Value]; ok {
				BlockBody := BlockScope[expr.AsCall.Value]
				VisitExpr(BlockBody)
			} else {
				fmt.Println("Error: undefined block '" + expr.AsCall.Value + "'")	
			}
		case ExprPlus:
			theStack.OpPlus()
		case ExprMinus:
			theStack.OpMinus()
		}
	}
	return
}


// -----------------------------
// ------------ Main -----------
// -----------------------------

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./main <filename>.t#")
		os.Exit(0)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error: file '" + os.Args[1] + "' does not exist")
		os.Exit(0)
	}

	lexer := LexerInit(file)
	parser := ParserInit(lexer)
	exprs, _ := ParserParse(parser)
	VisitExpr(exprs)

	return
}
