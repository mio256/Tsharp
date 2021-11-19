<div align="center">
    <h1> The T# Programming Language</h1>
    <a href="https://github.com/Tsharp-lang/Tsharp/blob/main/DOC/docs.md">Docs</a>
    |
    <a href="https://github.com/Tsharp-lang/Tsharp/blob/main/DOC/ドキュメント.md">Docs(日本語)</a>
    |
    <a href="https://github.com/ibukiyoshidaa/Tsharp/blob/main/editor/tsharp.vim">Vim</a>
</div>

[![Ubuntu](https://github.com/Tsharp-lang/Tsharp/actions/workflows/tsharp-ci-ubuntu.yml/badge.svg)](https://github.com/Tsharp-lang/Tsharp/actions/workflows/tsharp-ci-ubuntu.yml)
[![Windows](https://github.com/Tsharp-lang/Tsharp/actions/workflows/tsharp-ci-windows.yml/badge.svg)](https://github.com/Tsharp-lang/Tsharp/actions/workflows/tsharp-ci-windows.yml)

WARNING! THIS LANGUAGE IS A WORK IN PROGRESS! ANYTHING CAN CHANGE AT ANY MOMENT WITHOUT ANY NOTICE!

Something like Forth and Porth, but written in Go.
<a href="https://en.wikipedia.org/wiki/Stack-oriented_programming">Stack-oriented programming</a>

### Install

> Install
```
go build -o tsh main.go
```

### Run

> Run
```
$ ./tsh <filename>.t#
```

> Hello World
```pascal
push "Hello World" print
```

> Block
```pascal
block Main do
    push "Hello World" print
end

call Main
```

> If Statement
```pascal
if push true do
    push "Hello World!" print
end
```

```pascal
if push false do
    push "Hello World" print
else
    push "Hello World else body" print
end
```

> For loop
```pascal
for push true do
    push "Hello World!" print
end
```

> Dup
```pascal
push "Hello World!" dup print print
```

> Drop
```pascal
push "Hello World" push "T# Programming Language" drop print
```

### Contributors

<a href="https://github.com/ibukiyoshidaa/Tsharp/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=ibukiyoshidaa/Tsharp" />
</a>
