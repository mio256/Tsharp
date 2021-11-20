package main

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
	"os"
	"strconv"
	"reflect"
	"github.com/fatih/color"
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
	TOKEN_BOOL
	TOKEN_ELSE
	TOKEN_DIV
	TOKEN_MUL
	TOKEN_IS_EQUALS
	TOKEN_NOT_EQUALS
	TOKEN_LESS_THAN
	TOKEN_GREATER_THAN
	TOKEN_REM
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
	TOKEN_BOOL: "TOKEN_BOOL",
	TOKEN_ELSE: "TOKEN_ELSE",
	TOKEN_DIV: "TOKEN_DIV",
	TOKEN_MUL: "TOKEN_MUL",
	TOKEN_IS_EQUALS: "TOKEN_IS_EQUALS",
	TOKEN_NOT_EQUALS: "TOKEN_NOT_EQUALS",
	TOKEN_LESS_THAN: "TOKEN_LESS_THAN",
	TOKEN_GREATER_THAN: "TOKEN_GREATER_THAN",
	TOKEN_REM: "TOKEN_REM",
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
			case '/': return lexer.pos, TOKEN_DIV, "/"
			case '*': return lexer.pos, TOKEN_MUL, "*"
		    case '<': return lexer.pos, TOKEN_LESS_THAN, "<"
			case '>': return lexer.pos, TOKEN_GREATER_THAN, ">"
			case '%': return lexer.pos, TOKEN_REM, "%"
			default:
				if unicode.IsSpace(r) {
					continue
				} else if r == '=' {
					r, _, err := lexer.reader.ReadRune()
					if r == '\n' {break}
					if err != nil {panic(err)}
					lexer.pos.column++
					if r == '=' {
						return lexer.pos, TOKEN_IS_EQUALS, "=="
					}
				} else if r == '!' {
					r, _, err := lexer.reader.ReadRune()
					if r == '\n' {break}
					if err != nil {panic(err)}
					lexer.pos.column++
					if r == '=' {
						return lexer.pos, TOKEN_NOT_EQUALS, "!="
					}
				} else if r == '#' {
					for {
						r, _, err := lexer.reader.ReadRune()
						if r == '\n' {break}
						if err != nil {panic(err)}
						lexer.pos.column++
					}
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
					} else if val == "true" || val == "false" {
						return startPos, TOKEN_BOOL, val
					} else if val == "else" {
						return startPos, TOKEN_ELSE, val
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
	ExprBool
	ExprIf
	ExprDup
	ExprDrop
	ExprExit
	ExprFor	
	ExprBinop // + - * / %
	ExprCompare // < > == !=
)

type Expr struct {
	Type ExprType
	AsInt int
	AsStr string
	AsPush *Push
	AsBlockdef *Blockdef
	AsCall *Call
	AsBool bool
	AsIf *If
	AsFor *For
	AsBiniop int
	AsCompare int
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

type If struct {
	Op []Expr
	Body []Expr
	ElseBody []Expr
}

type For struct {
	Op []Expr
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
	case TOKEN_BOOL:
		expr.Type = ExprBool
		if parser.current_token_value == "true" {
			expr.AsBool = true
		} else {
			expr.AsBool = false
		}
		parser.ParserEat(TOKEN_BOOL)
	default:
		fmt.Println("Error: unexpected token value '" + parser.current_token_value + "'")
		os.Exit(0)
	}
	return expr
}

func ParserParse(parser *Parser)  ([]Expr, Parser) {
	exprs := []Expr{}

	for {
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
			} else if parser.current_token_value == "dup" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprDup
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "drop" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprDrop
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "exit" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprExit
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
			} else if parser.current_token_value == "for" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprFor
				op, _ := ParserParse(parser)
				parser.ParserEat(TOKEN_DO)
				body, _ := ParserParse(parser)
				parser.ParserEat(TOKEN_END)
				expr.AsFor = &For{
					Op: op,
					Body: body,
				}
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "if" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprIf
				op, _ := ParserParse(parser)
				parser.ParserEat(TOKEN_DO)
				body, _ := ParserParse(parser)
				if parser.current_token_type == TOKEN_ELSE {
					parser.ParserEat(TOKEN_ELSE)
					ElseBody, _ := ParserParse(parser)
					parser.ParserEat(TOKEN_END)
					expr.AsIf = &If{
						Op: op,
						Body: body,
						ElseBody: ElseBody,
					}
					exprs = append(exprs, expr)
				} else {
					parser.ParserEat(TOKEN_END)
					expr.AsIf = &If{
						Op: op,
						Body: body,
					}
					exprs = append(exprs, expr)
				}
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
			expr.Type = ExprBinop
			expr.AsBiniop = TOKEN_PLUS
			parser.ParserEat(TOKEN_PLUS)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_MINUS {
			expr.Type = ExprBinop
			expr.AsBiniop = TOKEN_MINUS
			parser.ParserEat(TOKEN_MINUS)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_DIV {
			expr.Type = ExprBinop
			expr.AsBiniop = TOKEN_DIV
			parser.ParserEat(TOKEN_DIV)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_MUL {
			expr.Type = ExprBinop
			expr.AsBiniop = TOKEN_MUL
			parser.ParserEat(TOKEN_MUL)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_REM {
			expr.Type = ExprBinop
			expr.AsBiniop = TOKEN_REM
			parser.ParserEat(TOKEN_REM)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_IS_EQUALS {
			expr.Type = ExprCompare
			expr.AsCompare = TOKEN_IS_EQUALS
			parser.ParserEat(TOKEN_IS_EQUALS)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_NOT_EQUALS {
			expr.Type = ExprCompare
			expr.AsCompare = TOKEN_NOT_EQUALS
			parser.ParserEat(TOKEN_NOT_EQUALS)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_LESS_THAN {
			expr.Type = ExprCompare
			expr.AsCompare = TOKEN_LESS_THAN
			parser.ParserEat(TOKEN_LESS_THAN)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_GREATER_THAN {
			expr.Type = ExprCompare
			expr.AsCompare = TOKEN_GREATER_THAN
			parser.ParserEat(TOKEN_GREATER_THAN)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_END {
			return exprs, *parser
		} else if parser.current_token_type == TOKEN_ELSE {
			return exprs, *parser
		} else if parser.current_token_type == TOKEN_DO {
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
	string_value *string
	int_value *int
	bool_value *bool
}

type Stack struct {
	Values []StackItem
}

var theStack Stack = Stack{}

func (stack *Stack) OpPush(item StackItem) {
	stack.Values = append(stack.Values, item)
}

func (stack *Stack) OpPrint() {
	if len(stack.Values) == 0 {
		fmt.Println("PrintError: the stack is empty")
		os.Exit(0)
	}

	if stack.Values[len(stack.Values)-1].bool_value != nil {
		fmt.Println(*stack.Values[len(stack.Values)-1].bool_value)
	} else if stack.Values[len(stack.Values)-1].string_value != nil {
		fmt.Println(*stack.Values[len(stack.Values)-1].string_value)
	} else if stack.Values[len(stack.Values)-1].int_value != nil {
		fmt.Println(*stack.Values[len(stack.Values)-1].int_value)
	}

	stack.Values = stack.Values[:len(stack.Values)-1]
	return
}

func (stack *Stack) OpDrop() {
	if len(stack.Values) == 0 {
		fmt.Println("DropError: the stack is empty")
		os.Exit(0)
	}
	stack.Values = stack.Values[:len(stack.Values)-1]
}

func (stack *Stack) OpBinop(value int) {
	var x int
	a := stack.Values[len(stack.Values)-1].int_value
	b := stack.Values[len(stack.Values)-2].int_value
	switch (value) {
		case TOKEN_PLUS:
			x = *a + *b
		case TOKEN_MINUS:
			x = *b - *a
		case TOKEN_MUL:
			x = *a * *b
		case TOKEN_DIV:
			x = *b / *a
		case TOKEN_REM:
			x = *b % *a
	}
	stack.Values = stack.Values[:len(stack.Values)-2]
	theStack.OpPush(StackItem{int_value: &x})
}

func (stack *Stack) OpDup() {
	if stack.Values[len(stack.Values)-1].int_value != nil {
		a := stack.Values[len(stack.Values)-1].int_value
		theStack.OpPush(StackItem{int_value: a})
	} else if stack.Values[len(stack.Values)-1].string_value != nil {
		a := stack.Values[len(stack.Values)-1].string_value
		theStack.OpPush(StackItem{string_value: a})
	} else if stack.Values[len(stack.Values)-1].bool_value != nil {
		a := stack.Values[len(stack.Values)-1].bool_value
		theStack.OpPush(StackItem{bool_value: a})
	}
}


// I will rewrite this function later
func (stack *Stack) OpCompare(value int) (bool) {
	if len(stack.Values) < 2 {
		fmt.Println("Error: expected more than two args in stack.")
		os.Exit(0)
	}
    var bool_value bool
	switch (value) {
		case TOKEN_IS_EQUALS:
			if reflect.TypeOf(stack.Values[len(stack.Values)-1]) != reflect.TypeOf(stack.Values[len(stack.Values)-2]) {
				bool_value = false
			} else if stack.Values[len(stack.Values)-1].int_value != nil {
				a := stack.Values[len(stack.Values)-1].int_value
				b := stack.Values[len(stack.Values)-2].int_value
				if *a != *b {bool_value = false} else {bool_value = true}
			} else if stack.Values[len(stack.Values)-1].string_value != nil {
				a := stack.Values[len(stack.Values)-1].string_value
				b := stack.Values[len(stack.Values)-2].string_value
				if *a != *b {bool_value = false} else {bool_value = true}
			} else if stack.Values[len(stack.Values)-1].bool_value != nil {
				a := stack.Values[len(stack.Values)-1].bool_value
				b := stack.Values[len(stack.Values)-2].bool_value
				if *a == *b {bool_value = true} else {bool_value = false}
			}
			stack.Values = stack.Values[:len(stack.Values)-2]
			return bool_value
		case TOKEN_NOT_EQUALS:
			if reflect.TypeOf(stack.Values[len(stack.Values)-1]) != reflect.TypeOf(stack.Values[len(stack.Values)-2]) {
				bool_value = true
			} else if stack.Values[len(stack.Values)-1].int_value != nil {
				a := stack.Values[len(stack.Values)-1].int_value
				b := stack.Values[len(stack.Values)-2].int_value
				if *a != *b {bool_value = true} else {bool_value = false}
			} else if stack.Values[len(stack.Values)-1].string_value != nil {
				a := stack.Values[len(stack.Values)-1].string_value
				b := stack.Values[len(stack.Values)-2].string_value
				if *a != *b {bool_value =  true} else {bool_value =  false}
			} else if stack.Values[len(stack.Values)-1].bool_value != nil {
				a := stack.Values[len(stack.Values)-1].bool_value
				b := stack.Values[len(stack.Values)-2].bool_value
				if *a == *b {bool_value = false} else {bool_value = true}
			}
			stack.Values = stack.Values[:len(stack.Values)-2]
			return bool_value
		case TOKEN_LESS_THAN:
			if stack.Values[len(stack.Values)-1].int_value == nil || stack.Values[len(stack.Values)-2].int_value == nil {
				fmt.Println("Error: type must be int")
				os.Exit(0)
			}
			a := stack.Values[len(stack.Values)-1].int_value
			b := stack.Values[len(stack.Values)-2].int_value
			stack.Values = stack.Values[:len(stack.Values)-2]
			return *b < *a
		case TOKEN_GREATER_THAN:
			if stack.Values[len(stack.Values)-1].int_value == nil || stack.Values[len(stack.Values)-2].int_value == nil {
				fmt.Println("Error: type must be int")
				os.Exit(0)
			}
			a := stack.Values[len(stack.Values)-1].int_value
			b := stack.Values[len(stack.Values)-2].int_value
			stack.Values = stack.Values[:len(stack.Values)-2]
			return *b > *a
		default:
			fmt.Println("Error: undifined type")
			os.Exit(0)
	}

	return true
}

func (stack *Stack) RetBool() (bool) {
	if len(stack.Values)-1 < 0 {
		fmt.Println("IfStatementError: the stack is empty. couldn't find bool.")
		os.Exit(0)
	}

	if stack.Values[len(stack.Values)-1].bool_value == nil {
		fmt.Println("Error: if op should be bool")
		os.Exit(0)
	}

	bool_value := *stack.Values[len(stack.Values)-1].bool_value
	stack.Values = stack.Values[:len(stack.Values)-1]
	return bool_value
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
				theStack.OpPush(StackItem{int_value: &expr.AsPush.Arg.AsInt})
			} else if expr.AsPush.Arg.Type == ExprStr {
				theStack.OpPush(StackItem{string_value: &expr.AsPush.Arg.AsStr})
			} else if expr.AsPush.Arg.Type == ExprBool {
				theStack.OpPush(StackItem{bool_value: &expr.AsPush.Arg.AsBool})
			}
		case ExprPrint:
			theStack.OpPrint()
		case ExprDup:
			theStack.OpDup()
		case ExprDrop:
			theStack.OpDrop()
		case ExprExit:
			os.Exit(0)
		case ExprBinop:
			theStack.OpBinop(expr.AsBiniop)
		case ExprCompare:
			bool_value := theStack.OpCompare(expr.AsCompare)
			theStack.OpPush(StackItem{bool_value: &bool_value})
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
		case ExprIf:
			VisitExpr(expr.AsIf.Op)
			bool_value := theStack.RetBool()
			if bool_value {
				VisitExpr(expr.AsIf.Body)
			} else {
				if expr.AsIf.ElseBody != nil {
					VisitExpr(expr.AsIf.ElseBody)
				}
			}
		case ExprFor:
			VisitExpr(expr.AsFor.Op)
			for theStack.RetBool() {
				VisitExpr(expr.AsFor.Body)
				VisitExpr(expr.AsFor.Op)
			}
		}
	}
	return
}

func Usage() {
	fmt.Println("Usage:")
	fmt.Println("  tsh <filename>.t#")
	os.Exit(0)
}

// -----------------------------
// ------------ Main -----------
// -----------------------------

func main() {
	if len(os.Args) != 2 || os.Args[1] == "help" {
		Usage()
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error: file '" + os.Args[1] + "' does not exist")

		whilte := color.New(color.FgWhite)

		fmt.Print("Run ")
		boldWhite := whilte.Add(color.BgCyan)
		boldWhite.Print(" tsh help ")
		fmt.Println(" for usage")

		os.Exit(0)
	}

	lexer := LexerInit(file)
	parser := ParserInit(lexer)
	exprs, _ := ParserParse(parser)
	VisitExpr(exprs)

	return
}
