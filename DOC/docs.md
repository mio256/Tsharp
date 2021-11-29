# T# Documentation

## Introduction

T# is a Stack-based programming language designed for building software.
It's similar to Porth and Forth.

## Install & Run
```bash
$ git clone https://github.com/Tsharp-lang/Tsharp
$ cd tsharp
$ go build main.go
$ ./main examples/main.t#
or
$ ./main.exe examples/main.t#
```

## Hello World
```pascal
push "Hello World!" print
```

'push' will push the value to the stack.
'print' will print the top element of the stack, then remove it.

## Comments
```python
# Sample comment
```

## Import
```python
import "main.t#"
```

## Block
```pascal
block main do
    push "Hello World!" print
end

call main
```

'block' is like Function in other languages.

## If Statement
```pascal
if push false do
    push "Hello World" print
else
    push "Hello World else body!" print
end

push 10 push 10 == print
push 20 push 10 != print
push 2 push 10 < print
push 10 push 2 > print
```

## Dup
```pascal
push "Hello World" dup print print
```
'dup' duplicate element on top of the stack.

## Drop
```pascal
push "Hello World" push "T# Programming Language" drop print
```
'drop' drops the top element of the stack.

## For loop
```pascal
for push true do
    push "Hello World!" print
    break
end
```

## Arithmetic
```pascal
push 34 push 35 + print

push 100 push 40 - print

push 200 push 5 / print

push 10 push 2 * print
```

## Variable
```python
x = "Hello World 'x'"

push x print

y = "New variable 'y'"

x = y

push x print
```

## Type
```python
push int # 12345
push string # "Hello World!"
push bool # true false
push type # int string bool type
```

## Typeof
```python
push "Hello World" dup typeof print
```

## FizzBuzz
```pascal
push 1
for dup push 101 < do
    if dup push 3 % push 0 == do
        if dup push 15 % push 0 == do
            push "FizzBuzz" print
        else
            push "Fizz" print
        end
    else
        if dup push 5 % push 0 == do
            push "Buzz" print
        else
            dup print
        end
    end
    push 1 +
end drop
```
