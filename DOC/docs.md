# T# Documentation

## Introduction

T# is a Stack-based programming language designed for building software.
It's similar to Porth and Forth.

## Install
```bash
$ git clone https://github.com/Tsharp-lang/Tsharp
$ cd tsharp
$ go build main.go
```

## Hello World
```pascal
push "Hello World" print
```

'push' will push the value to the stack.
'print' will print the first arg inside the stack, then remove it.

## Block
```pascal
block main do
    push "Hello World" print
end
```

'block' is like a Function in other languages.

