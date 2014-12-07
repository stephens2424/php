php
===

Parser for PHP written in Go

See [this post](https://stephensearles.com/?p=288) for an introduction.

[![Build Status](https://travis-ci.org/stephens2424/php.svg)](https://travis-ci.org/stephens2424/php) [![GoDoc](https://godoc.org/github.com/stephens2424/php?status.svg)](https://godoc.org/github.com/stephens2424/php)

Test console:

[![console](https://stephensearles.com/wp-content/uploads/2014/07/Screen-Shot-2014-07-27-at-12.02.32-PM.png)](https://phpconsole.stephensearles.com)

## Project Components

- php: the core parser
- php/ast: (Abstract syntax tree) describes the nodes in PHP as parsed by the parser
- php/ast/printer: prints an ast back to source code (heavily a work in progress)
- php/cmd: a tool used to debug the parser
- php/lexer: reads a stream of tokens from source code
- php/passes: tools and packages related to modifying or analyzing PHP code (heavily a work in progress)
- php/query: tools and packages related to analyzing and finding things in PHP code (heavily a work in progress)
- php/testfiles: simple examples of PHP that must parse with no errors for tests to pass
- php/token: describes the tokens read by the lexer
