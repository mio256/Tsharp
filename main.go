package main

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
	"os"
	"strconv"
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
	TOKEN_EQUALS
	TOKEN_IS_EQUALS
	TOKEN_NOT_EQUALS
	TOKEN_LESS_THAN
	TOKEN_GREATER_THAN
	TOKEN_LESS_EQUALS
	TOKEN_GREATER_EQUALS
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
	TOKEN_EQUALS: "TOKEN_EQUALS",
	TOKEN_IS_EQUALS: "TOKEN_IS_EQUALS",
	TOKEN_NOT_EQUALS: "TOKEN_NOT_EQUALS",
	TOKEN_LESS_THAN: "TOKEN_LESS_THAN",
	TOKEN_GREATER_THAN: "TOKEN_GREATER_THAN",
	TOKEN_LESS_EQUALS: "TOKEN_LESS_EQUALS",
	TOKEN_GREATER_EQUALS: "TOKEN_GREATER_EQUALS",
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
					} else {
						return lexer.pos, TOKEN_EQUALS, "="
					}
				} else if r == '<' {
					r, _, err := lexer.reader.ReadRune()
					if err != nil {panic(err)}
					if r == '=' {
						lexer.pos.column++
						return lexer.pos, TOKEN_LESS_EQUALS, "<="
					} else {
						return lexer.pos, TOKEN_LESS_THAN, "<"
					}
				} else if r == '>' {
					r, _, err := lexer.reader.ReadRune()
					if err != nil {panic(err)}
					if r == '=' {
						lexer.pos.column++
						return lexer.pos, TOKEN_GREATER_EQUALS, ">="
					} else {
						return lexer.pos, TOKEN_GREATER_THAN, ">"
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
	ExprId
	ExprPush
	ExprBlockdef
	ExprPrint
	ExprSwap
	ExprImport
	ExprCall
	ExprBool
	ExprIf
	ExprDup
	ExprDrop
	ExprExit
	ExprFor	
	ExprBinop // + - * / %
	ExprCompare // < > == !=
	ExprVardef
)

type Expr struct {
	Type ExprType
	AsInt int
	AsStr string
	AsId string
	AsPush *Push
	AsBlockdef *Blockdef
	AsCall *Call
	AsBool bool
	AsIf *If
	AsFor *For
	AsBiniop int
	AsCompare int
	AsImport string
	AsVardef *Vardef
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

type Vardef struct {
	Name string
	Arg Expr
}


// -----------------------------
// ----------- Parse -----------
// -----------------------------

type Parser struct {
	current_token_type Token
	current_token_value string
	lexer Lexer
	line int
	column int
}

func ParserInit(lexer *Lexer) *Parser {
	pos, tok, val := lexer.Lex()
	return &Parser{
		current_token_type: tok,
		current_token_value: val,
		lexer: *lexer,
		line: pos.line,
		column: pos.column,
	}
}

func (parser *Parser) ParserEat(token Token) {
	if token != parser.current_token_type {
		fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value '%s'", parser.line, parser.column, parser.current_token_value))
		os.Exit(0)
	}
	pos, tok, val := parser.lexer.Lex()
	parser.current_token_type = tok
	parser.current_token_value = val
	parser.line = pos.line
	parser.column = pos.column
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
		case TOKEN_ID:
			expr.Type = ExprId
			expr.AsId = parser.current_token_value
			parser.ParserEat(TOKEN_ID)
		default:
			fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value '%s'", parser.line, parser.column, parser.current_token_value))
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
			} else if parser.current_token_value == "swap" {
				parser.ParserEat(TOKEN_ID)
				expr.Type = ExprSwap
				exprs = append(exprs, expr)
			} else if parser.current_token_value == "import" {
				parser.ParserEat(TOKEN_ID)
				if parser.current_token_type != TOKEN_STRING {
					fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value '%s'", parser.line, parser.column, parser.current_token_value))
					os.Exit(0)
				}
				expr.Type = ExprImport
				expr.AsImport = parser.current_token_value
				parser.ParserEat(TOKEN_STRING)
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
					fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value '%s'", parser.line, parser.column, parser.current_token_value))
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
					fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value '%s'", parser.line, parser.column, parser.current_token_value))
					os.Exit(0)
				}
				expr.Type = ExprCall
				expr.AsCall = &Call{
					Value: parser.current_token_value,
				}
				parser.ParserEat(TOKEN_ID)
				exprs = append(exprs, expr)
			} else {
				vname := parser.current_token_value
				parser.ParserEat(TOKEN_ID)
				if parser.current_token_type != TOKEN_EQUALS {
					fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value '%s'", parser.line, parser.column, parser.current_token_value))
					os.Exit(0)
				}
				parser.ParserEat(TOKEN_EQUALS)
				expr.Type = ExprVardef
				expr.AsVardef = &Vardef{
					Name: vname,
					Arg: ParserParseExpr(parser),
				}
				exprs = append(exprs, expr)
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
		} else if parser.current_token_type == TOKEN_GREATER_EQUALS {
			expr.Type = ExprCompare
			expr.AsCompare = TOKEN_GREATER_EQUALS
			parser.ParserEat(TOKEN_GREATER_EQUALS)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_LESS_EQUALS {
			expr.Type = ExprCompare
			expr.AsCompare = TOKEN_LESS_EQUALS
			parser.ParserEat(TOKEN_LESS_EQUALS)
			exprs = append(exprs, expr)
		} else if parser.current_token_type == TOKEN_END || parser.current_token_type == TOKEN_ELSE || parser.current_token_type == TOKEN_DO || parser.current_token_type == TOKEN_EOF {
			return exprs, *parser
		} else {
			fmt.Println(fmt.Sprintf("SyntaxError:%d:%d: unexpected token value '%s'", parser.line, parser.column, parser.current_token_value))
			os.Exit(0)
		}
	}

	return exprs, *parser
}


// -----------------------------
// ----------- Stack -----------
// -----------------------------

var Stack = []Expr{}

func OpPush(item Expr) {
	if item.Type == ExprId {
		if _, ok := VariableScope[item.AsId]; ok {
			item = VariableScope[item.AsId]
		} else {
			fmt.Println("Error: undefined variable '" + item.AsId + "'")
			os.Exit(0)
		}
	}
	Stack = append(Stack, item)
}

func OpDrop() {
	if len(Stack)-1 < 0 {
		fmt.Println("DropError: the stack is empty.")
		os.Exit(0)
	}

	Stack = Stack[:len(Stack)-1]
}

func OpDup() {
	if len(Stack)-1 < 0 {
		fmt.Println("DupError: the stack is empty.")
		os.Exit(0)
	}

	visitedExpr := Stack[len(Stack)-1]
	Stack = append(Stack, visitedExpr)
}

func OpSwap() {
	if len(Stack) < 2 {
		fmt.Println("SwapError: expected more than two args in stack.")
		os.Exit(0)
	}

	visitedExpr := Stack[len(Stack)-1]
	visitedExprSecond := Stack[len(Stack)-2]
	OpDrop()
	OpDrop()
	OpPush(visitedExpr)
	OpPush(visitedExprSecond)
}

func OpPrint() {
	if len(Stack) == 0 {
		fmt.Println("PrintError: the stack is empty")
		os.Exit(0)
	}

	visitedExpr := Stack[len(Stack)-1]
	switch (visitedExpr.Type) {
		case ExprInt: fmt.Println(visitedExpr.AsInt)
		case ExprStr: fmt.Println(visitedExpr.AsStr)
		case ExprBool: fmt.Println(visitedExpr.AsBool)
	}

	OpDrop()
}

func OpCompare(value int) (bool) {
	if len(Stack) < 2 {
		fmt.Println("Error: expected more than two args in stack.")
		os.Exit(0)
	}

	visitedExpr := Stack[len(Stack)-1]
	visitedExprSecond := Stack[len(Stack)-2]

	OpDrop()
	OpDrop()

	if value == TOKEN_IS_EQUALS {
		if visitedExpr.Type != visitedExprSecond.Type {
			return false
		}

		if visitedExpr.Type == ExprInt {
			return visitedExpr.AsInt == visitedExprSecond.AsInt
		}

		if visitedExpr.Type == ExprStr {
			return visitedExpr.AsStr == visitedExprSecond.AsStr
		}

		if visitedExpr.Type == ExprBool {
			return visitedExpr.AsBool == visitedExprSecond.AsBool
		}
	}

	if value == TOKEN_NOT_EQUALS {
		if visitedExpr.Type != visitedExprSecond.Type {
			return true
		}

		if visitedExpr.Type == ExprInt {
			return visitedExpr.AsInt != visitedExprSecond.AsInt
		}

		if visitedExpr.Type == ExprStr {
			return visitedExpr.AsStr != visitedExprSecond.AsStr
		}

		if visitedExpr.Type == ExprBool {
			return visitedExpr.AsBool != visitedExprSecond.AsBool
		}
	}
    
	if visitedExpr.Type != ExprInt || visitedExprSecond.Type != ExprInt {
		fmt.Println("TypeError: '<' expected type int")
		os.Exit(0)
	}

	if value == TOKEN_LESS_THAN {
		return visitedExprSecond.AsInt < visitedExpr.AsInt
	}

	if value == TOKEN_GREATER_THAN {
		return visitedExprSecond.AsInt > visitedExpr.AsInt
	}

	if value == TOKEN_GREATER_EQUALS {
		return visitedExprSecond.AsInt >= visitedExpr.AsInt
	}

	if value == TOKEN_LESS_EQUALS {
		return visitedExprSecond.AsInt <= visitedExpr.AsInt
	}

	return false
}

func RetBool() (bool) {
	if len(Stack)-1 < 0 {
		fmt.Println("Error: the stack is empty. couldn't find bool.")
		os.Exit(0)
	}

	visitedExpr := Stack[len(Stack)-1]
	if visitedExpr.Type != ExprBool {
		fmt.Println("Error: if op should be bool")
		os.Exit(0)
	}
	bool_value := visitedExpr.AsBool
	OpDrop()
	return bool_value
}

func OpIf(expr Expr) {
	VisitExpr(expr.AsIf.Op)
	bool_value := RetBool()
	if bool_value {
		VisitExpr(expr.AsIf.Body)
	} else {
		if expr.AsIf.ElseBody != nil {
			VisitExpr(expr.AsIf.ElseBody)
		}
	}
}

func OpCondition(expr Expr) {
	bool_value := OpCompare(expr.AsCompare)
	BoolExpr := Expr{}
	BoolExpr.Type = ExprBool
	BoolExpr.AsBool = bool_value
	PushExpr := Expr{}
	PushExpr.Type = ExprPush
	PushExpr.AsPush = &Push{
		Arg: BoolExpr,
	}
	OpPush(PushExpr.AsPush.Arg)
}

func OpBinop(value int) {
	if len(Stack) < 2 {
		fmt.Println("Error: expected more than two args in stack.")
		os.Exit(0)
	}

	var finalInt int
	visitedExpr := Stack[len(Stack)-1]
	visitedExprSecond := Stack[len(Stack)-2]
	OpDrop()
	OpDrop()

	if visitedExpr.Type != ExprInt || visitedExprSecond.Type != ExprInt {
		fmt.Println("TypeError: binary operation expected type int")
		os.Exit(0)
	}

	switch (value) {
		case TOKEN_PLUS:
			finalInt = visitedExpr.AsInt + visitedExprSecond.AsInt
		case TOKEN_MINUS:
			finalInt = visitedExprSecond.AsInt - visitedExpr.AsInt
		case TOKEN_MUL:
			finalInt = visitedExpr.AsInt * visitedExprSecond.AsInt
		case TOKEN_DIV:
			finalInt = visitedExprSecond.AsInt / visitedExpr.AsInt
		case TOKEN_REM:
			finalInt = visitedExprSecond.AsInt % visitedExpr.AsInt
	}

	IntExpr := Expr{}
	IntExpr.Type = ExprInt
	IntExpr.AsInt = finalInt
	PushExpr := Expr{}
	PushExpr.Type = ExprPush
	PushExpr.AsPush = &Push{
		Arg: IntExpr,
	}
	OpPush(PushExpr.AsPush.Arg)
}

func OpImport(expr Expr) {
	file, err := os.Open(expr.AsImport)
	if err != nil {
		panic(err)
	}
	lexer := LexerInit(file)
	parser := ParserInit(lexer)
	exprs, _ := ParserParse(parser)
	VisitExpr(exprs)
}

func OpFor(expr Expr) {
	VisitExpr(expr.AsFor.Op)
	for RetBool() {
		VisitExpr(expr.AsFor.Body)
		VisitExpr(expr.AsFor.Op)
	}
}


// -----------------------------
// ---------- Variable ---------
// -----------------------------

var VariableScope = map[string]Expr{}

func OpVardef(expr Expr) {
	if expr.AsVardef.Arg.Type == ExprId {
		if _, ok := VariableScope[expr.AsVardef.Arg.AsId]; ok {
			value := VariableScope[expr.AsVardef.Arg.AsId]
			VariableScope[expr.AsVardef.Name] = value
		} else {
			fmt.Println("Error: undefined variable '" + expr.AsVardef.Arg.AsId + "'")
			os.Exit(0)
		}
	} else {
		VariableScope[expr.AsVardef.Name] = expr.AsVardef.Arg
	}
}


// -----------------------------
// ----------- Block -----------
// -----------------------------

var BlockScope = map[string][]Expr{}

func OpBlockdef(expr Expr) {
	if _, ok := BlockScope[expr.AsBlockdef.Name]; ok {
		fmt.Println("Error: block '%s' is already defined", expr.AsBlockdef.Name)
		os.Exit(0)
	}
	BlockScope[expr.AsBlockdef.Name] = expr.AsBlockdef.Body
}

func OpCallBlock(expr Expr) {
	if _, ok := BlockScope[expr.AsCall.Value]; ok {
		BlockBody := BlockScope[expr.AsCall.Value]
		VisitExpr(BlockBody)
	} else {
		fmt.Println("Error: undefined block '" + expr.AsCall.Value + "'")
		os.Exit(0)
	}
}


// -----------------------------
// ----------- Visit -----------
// -----------------------------

func VisitExpr(exprs []Expr) {
	for _, expr := range exprs {
		switch expr.Type {
			case ExprPush:
				OpPush(expr.AsPush.Arg)
			case ExprPrint:
				OpPrint()
			case ExprSwap:
				OpSwap()
			case ExprImport:
				OpImport(expr)
			case ExprDup:
				OpDup()
			case ExprDrop:
				OpDrop()
			case ExprExit:
				os.Exit(0)
			case ExprBinop:
				OpBinop(expr.AsBiniop)
			case ExprCompare:
				OpCondition(expr)
			case ExprBlockdef:
				OpBlockdef(expr)
			case ExprCall:
				OpCallBlock(expr)
			case ExprIf:
				OpIf(expr)
			case ExprFor:
				OpFor(expr)
			case ExprVardef:
				OpVardef(expr)
		}
	}
}


// -----------------------------
// ----------- Main -----------
// -----------------------------

func Usage() {
	fmt.Println("Usage:")
	fmt.Println("  tsh <filename>.t#")
	os.Exit(0)
}


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
}
