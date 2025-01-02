# go-enumer

[![GoDoc](https://godoc.org/github.com/oesand/go-enumer?status.svg)](https://godoc.org/github.com/oesand/go-enumer)

Another one code generator for go

## How it works

go-enumer will take a commented type declaration and generate all the necessary constants or functions [(see wiki)](https://github.com/oesand/go-enumer/wiki#supported-generations).

## Installation

Just run:
``` bash
go install github.com/oesand/go-enumer
```

Or you can download a release directly from [github](https://github.com/oesand/go-enumer/releases).

## Adding it to your project

### Using go generate

1. Add a go:generate line to your file like so... `//go:generate go-enumer gen`
1. Run go generate like so `go generate ./...`

### Syntax

The parser looks for comments on your type definitions and parses enum declarations from them. 
The parser will look for `enum(` (case-insensitive) and continue looking for comma-separated values until it finds `)`. 
You can put values on one line or on multiple lines. Generative options are also supported. 
They are written after `)` either as a word or the parameterized value `key: value`

[Examples can be found in the example folder](./example/)
