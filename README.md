<div align="center">    
    <img width="100px" src="https://user-images.githubusercontent.com/81926489/143374038-059715ef-a83d-479d-a8c3-56ea57b8cc8e.PNG">
    <h1> The T# Programming Language</h1>
    <a href="https://github.com/Tsharp-lang/Tsharp/blob/main/DOC/docs.md">Docs</a>
    |
    <a href="https://github.com/Tsharp-lang/Tsharp/blob/main/DOC/ドキュメント.md">Docs(日本語)</a>
    |
    <a href="https://github.com/Tsharp-lang/Tsharp/blob/main/editor/tsharp.vim">Vim</a>
    |
    <a href="https://marketplace.visualstudio.com/items?itemName=akamurasaki.tsharplanguage-color">VSCode</a>
</div>

[![Ubuntu](https://github.com/Tsharp-lang/Tsharp/actions/workflows/tsharp-ci-ubuntu.yml/badge.svg)](https://github.com/Tsharp-lang/Tsharp/actions/workflows/tsharp-ci-ubuntu.yml)
[![Windows](https://github.com/Tsharp-lang/Tsharp/actions/workflows/tsharp-ci-windows.yml/badge.svg)](https://github.com/Tsharp-lang/Tsharp/actions/workflows/tsharp-ci-windows.yml)
[![CodeQL](https://github.com/Tsharp-lang/Tsharp/actions/workflows/codeql-analysis.yml/badge.svg?branch=main)](https://github.com/Tsharp-lang/Tsharp/actions/workflows/codeql-analysis.yml)

WARNING! THIS LANGUAGE IS A WORK IN PROGRESS! ANYTHING CAN CHANGE AT ANY MOMENT WITHOUT ANY NOTICE!

It's like Forth and Porth, but written in Go.
<a href="https://en.wikipedia.org/wiki/Stack-oriented_programming">Stack-oriented programming</a>

### TODO
- [ ] Self-hosted

### Install

> Install
```
go build main.go
```

### Run

> Run
```
$ ./main <filename>.t#

or

$ ./main.exe <filename>.t#
```

> Hello World
```pascal
"Hello World" print
```

> Block
```pascal
block Main do
    "Hello World" print
end

call Main
```

> If Statement
```pascal
if true do
    "Hello World!" print
end
```

```pascal
if false do
    "Hello World" print
else
    "Hello World else body" print
end
```

> For loop
```pascal
for true do
    "Hello World!" print
end
```

> Dup
```pascal
"Hello World!" dup print print
```

> Drop
```pascal
"Hello World" "T# Programming Language" drop print
```

> Arithmetic
```pascal
34 35 + print

100 40 - print

200 5 / print

10 2 * print
```

> Variable
```pascal
10 -> x drop

x -> y drop

y print
```

> FizzBuzz
```pascal
1
for dup 101 < do
    if dup 3 % 0 == do
        if dup 15 % 0 == do
            "FizzBuzz" print
        else
            "Fizz" print
        end
    else
        if dup 5 % 0 == do
            "Buzz" print
        else
            dup print
        end
    end
    inc
end drop
```

