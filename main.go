// I will rewrite this today (Tuesday, November 16~17, 2021)

package main

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
	"os"
	_"strconv"
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
}

func (t Token) String() string {
	return tokens[t]
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
					lit := lexer.lexInt()
					return startPos, TOKEN_INT, lit
				} else if unicode.IsLetter(r) {
					startPos := lexer.pos
					lexer.backup()
					lit := lexer.lexId()
					if lit == "end" {
						return startPos, TOKEN_END, lit
					}
					return startPos, TOKEN_ID, lit
				} else if r == '"' {
					startPos := lexer.pos
					lexer.backup()
					lit := lexer.lexString()
					r, _, err = lexer.reader.ReadRune()
					return startPos, TOKEN_STRING, lit
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
	var lit string
	for {
		r, _, err := lexer.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit
			}
		}
        lexer.pos.column++
		if unicode.IsLetter(r) {
			lit = lit + string(r)
		} else {
			lexer.backup()
			return lit
		}
	}
}

func (lexer *Lexer) lexInt() string {
	var lit string
	for {
		r, _, err := lexer.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit
			}
		}
		lexer.pos.column++
		if unicode.IsDigit(r) {
			lit = lit + string(r)
		} else {
			lexer.backup()
			return lit
		}
	}
}

func (lexer *Lexer) lexString() string {
	var lit string
	r, _, err := lexer.reader.ReadRune()
	for {
		r, _, err = lexer.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return lit
			}
		}
		lexer.pos.column++
		if r != '"' {
			lit = lit + string(r)
		} else {
			lexer.backup()
			return lit
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
)

type Expr struct {
	Type ExprType
	AsInt int
	AsStr string
	AsPush *Push
	AsBlockdef *Blockdef
}

type Push struct {
	Arg Expr
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

func ParseInit(lexer *Lexer) (Parser) {
	_, tok, lit := lexer.Lex()
	parser := Parser{current_token_type: tok, current_token_value: lit, lexer: *lexer}
	return parser
}

func ParserEat(parser Parser, token Token) (Parser) {
	if token != parser.current_token_type {
		fmt.Println("Error: unexpected token value '" + parser.current_token_value + "'")
		os.Exit(0)
	}

	//fmt.Println("ParserEat func eated value: '" + parser.current_token_value + "'")
	_, tok, lit := parser.lexer.Lex()

	parser = Parser{current_token_type: tok, current_token_value: lit, lexer: parser.lexer}

	return parser
}

func ParserParseExpr(parser Parser) (Expr, Parser) {
	expr := Expr{}
	switch parser.current_token_type {
	case TOKEN_INT:
		expr.Type = ExprStr
		expr.AsStr = parser.current_token_value
		parser = ParserEat(parser, TOKEN_INT)
	case TOKEN_STRING:
		expr.Type = ExprStr
		expr.AsStr = parser.current_token_value
		parser = ParserEat(parser, TOKEN_STRING)
	default:
		fmt.Println("Error: unexpected token value '" + parser.current_token_value + "'")
		os.Exit(0)
	}
	return expr, parser
}

func ParserParseId(parser Parser) (Expr, Parser) {
	expr := Expr{}
	if parser.current_token_value == "push" {
		parser = ParserEat(parser, TOKEN_ID)
		expr.Type = ExprPush
		argExpr, parser := ParserParseExpr(parser)
		expr.AsPush = &Push{
			Arg: argExpr,
		}
		//fmt.Println("current token value: '" + parser.current_token_value + "'")

		return expr, parser
	} else if parser.current_token_value == "block" {
		parser = ParserEat(parser, TOKEN_ID)
		expr.Type = ExprBlockdef
		name := parser.current_token_value
		parser = ParserEat(parser, TOKEN_ID)
		parser = ParserEat(parser, TOKEN_ID)
		exprs := []Expr{}
		exprs, parser = ParserParse(parser)
		expr.AsBlockdef = &Blockdef{
			Name: name,
			Body: exprs,
		}
		//fmt.Println("current token value after parsing block: '" + parser.current_token_value + "'")
		parser = ParserEat(parser, TOKEN_END)
		return expr, parser
	} else if parser.current_token_value == "print" {
		parser = ParserEat(parser, TOKEN_ID)
		expr.Type = ExprPrint
		//fmt.Println("current token value: '" + parser.current_token_value + "'")

		return expr, parser
	} else {
		fmt.Println("Error: unexpected token value '" + parser.current_token_value + "'")
		os.Exit(0)
	}

	return expr, parser
}

func ParserParse(parser Parser) ([]Expr, Parser) {
	exprs := []Expr{}
    expr := Expr{}

	for parser.current_token_type != TOKEN_EOF || parser.current_token_type != TOKEN_END {
		if parser.current_token_type == TOKEN_ID {
			expr, parser = ParserParseId(parser)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_EOF {
			return exprs, parser
		} else if parser.current_token_type == TOKEN_END {
			return exprs, parser
		} else {
			fmt.Println("SyntaxError: unexpected token value '" + parser.current_token_value + "'")
			os.Exit(0)
		}
	}

	return exprs, parser
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
		fmt.Println(stack.Values[len(stack.Values)-1].string_value)
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
// ---------- Visitor ----------
// -----------------------------

func VisitExpr(exprs[] Expr) {
	for _, expr := range exprs {
		switch expr.Type {
		case ExprPush:
			if expr.AsPush.Arg.Type == ExprInt {
				theStack.OpPush(StackItem{int_value: expr.AsPush.Arg.AsInt})
				// fmt.Println("push value: ", expr.AsPush.Arg.AsInt)
			} else if expr.AsPush.Arg.Type == ExprStr {
				theStack.OpPush(StackItem{string_value: expr.AsPush.Arg.AsStr})
				//fmt.Println("push value: ", expr.AsPush.Arg.AsStr)
			}
		case ExprPrint:
			theStack.OpPrint()
		case ExprBlockdef:
			//fmt.Println("----------- Block Body ------------")
			VisitExpr(expr.AsBlockdef.Body)
			//fmt.Println("------------ Block End ------------")
		}
	}
}

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	lexer := LexerInit(file)

	parser := ParseInit(lexer)

	exprs, _ := ParserParse(parser)

	VisitExpr(exprs)

	return
}

