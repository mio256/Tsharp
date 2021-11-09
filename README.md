<div align="center">
    <h1> The T# Programming Language</h1>
</div>

[![T# Build](https://github.com/ibukiyoshidaa/Tsharp/actions/workflows/tsharp-ci.yml/badge.svg)](https://github.com/ibukiyoshidaa/Tsharp/actions/workflows/tsharp-ci.yml)

WARNING! THIS LANGUAGE IS A WORK IN PROGRESS! ANYTHING CAN CHANGE AT ANY MOMENT WITHOUT ANY NOTICE!

<!---- Compile to C. ---->

### Install

> Install
```
make
```

### Run

> Run
```
$ ./tsh.out main <filename>.t#
```

> Hello World
```pascal
func main() do
    print("Hello World");
end;
```

> Variable
```pascal
func main() do
    name = "T#";
    print(name);

    num = 1234;
    print(num);
end;
```

> If statement
```pascal
func main() do
    if 1 do
        print("Hello World");
    end;

    if 0 do
        print("Hello World");
    else
        print("else!");
    end;
end;
```

### Vim Syntax Highlighting
- <a href="https://github.com/ibukiyoshidaa/Tsharp/blob/main/editor/tsharp.vim">Vim</a>


### Contributors

<a href="https://github.com/ibukiyoshidaa/Tsharp/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=ibukiyoshidaa/Tsharp" />
</a>
