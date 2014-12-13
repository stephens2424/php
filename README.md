php
===

Parser for PHP written in Go

See [this post](https://stephensearles.com/?p=288) for an introduction.

[![Build Status](https://travis-ci.org/stephens2424/php.svg)](https://travis-ci.org/stephens2424/php) [![GoDoc](https://godoc.org/github.com/stephens2424/php?status.svg)](https://godoc.org/github.com/stephens2424/php)

Test console:

[![console](https://stephensearles.com/wp-content/uploads/2014/07/Screen-Shot-2014-07-27-at-12.02.32-PM.png)](https://phpconsole.stephensearles.com)

## Project Status

This project is under heavy development, though some pieces are more or less stable. Listed here are components that in progress or are ideas for future development

- Lexer and Parser: these pieces are quite useable and haven't seen more than minor changes for a while. They do have notable gaps, though (e.g. namespaces).
- Code search and symbol lookup: basic idea implemented, many many details missing
- Code formatting: basic idea implemented, formatting needs to narrow down to PSR-2
- Transpilation to Go: not begun.
- Scoping: not begun.
- Type inferencing: not begun.
- Dead code analysis: not begun.

## Project Components

- php: the core parser
- php/ast: (Abstract syntax tree) describes the nodes in PHP as parsed by the parser
- php/ast/printer: prints an ast back to source code
- php/cmd: a tool used to debug the parser
- php/lexer: reads a stream of tokens from source code
- php/passes: tools and packages related to modifying or analyzing PHP code (heavily a work in progress)
- php/query: tools and packages related to analyzing and finding things in PHP code (heavily a work in progress)
- php/testfiles: simple examples of PHP that must parse with no errors for tests to pass
- php/token: describes the tokens read by the lexer
