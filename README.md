<div align="center">
    <h1> The T# Programming Language</h1>
</div>

[![Ubuntu](https://github.com/Tsharp-lang/Tsharp/actions/workflows/tsharp-ci-ubuntu.yml/badge.svg)](https://github.com/Tsharp-lang/Tsharp/actions/workflows/tsharp-ci-ubuntu.yml)
[![Windows](https://github.com/Tsharp-lang/Tsharp/actions/workflows/tsharp-ci-windows.yml/badge.svg)](https://github.com/Tsharp-lang/Tsharp/actions/workflows/tsharp-ci-windows.yml)

WARNING! THIS LANGUAGE IS A WORK IN PROGRESS! ANYTHING CAN CHANGE AT ANY MOMENT WITHOUT ANY NOTICE!

Something like Forth and Porth, but written in Go.
<a href="https://en.wikipedia.org/wiki/Stack-oriented_programming">Stack-oriented programming</a>

### Install

> Install
```
go build main.go
```

### Run

> Run
```
$ ./main <filename>.t#
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


### Vim Syntax Highlighting
- <a href="https://github.com/ibukiyoshidaa/Tsharp/blob/main/editor/tsharp.vim">Vim</a>


### Contributors

<a href="https://github.com/ibukiyoshidaa/Tsharp/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=ibukiyoshidaa/Tsharp" />
</a>
