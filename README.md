# meego
A miniature bootstrapped Golang to 'Vanilla' C++ transpiler.

Lexer/Parser generated using [gocc](https://github.com/goccmack/gocc), following
the stand-alone DFA which recognizes grammer for a regular language in BNF (Backus Naur Form).

* Install meego and gocc
```
$ go get github.com/aniketp/meego
$ go get github.com/goccmack/gocc	(make deps)
```

* Generate Lexer and Parser
```
$ cd src; gocc lang.bnf			(make run)
```

* Run tests
```
$ cd test; go test -v			(make test)
```

* Compile a simple program
```
$ go run main.go input/example.meego
5
Requiescat in pace, Ezio!
```

This project is my attempt to learn about Compiler Design, and was done
in a short duration following this [medium article](https://medium.freecodecamp.org/write-a-compiler-in-go-quick-guide-30d2f33ac6e0),
including my own variations on the top. As a result, the grammer is a tiny subset
of Golang (with a mix of Typescript syntax).

To avoid the complexity of Intermediate Language generation and optimization, I
switched the Target Language to a simple subset of C++11.

### References
* [shivansh/gogo](https://github.com/shivansh/gogo) ([lang.bnf](./src/lang.bnf))
* [Lebonesco/go-compiler](https://github.com/Lebeneco/go-compiler)