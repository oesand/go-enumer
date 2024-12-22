# go-enumer

[![GoDoc](https://godoc.org/github.com/oesand/go-enumer?status.svg)](https://godoc.org/github.com/oesand/go-enumer)

An enum generator for go

## How it works

go-enumer will take a commented type declaration like this:

```go
// enum(pending, running, completed)
type IntStatus int

// enum(pending, running, completed)
type StrStatus string
```

and generate a file with definitions along various optional niceties that you may need:

```go
// Code generated by go-enumer[https://github.com/oesand/go-enumer]. DO NOT EDIT! 

package main

// IntStatus enum declarations
const (
	IntStatusPending IntStatus = iota
	IntStatusRunning
	IntStatusCompleted
)


func IntStatusNames() []string
func IntStatusFromString(value string) (IntStatus, bool) 

func (en IntStatus) IsValid() bool 
func IntStatusValues() []IntStatus 

func (en IntStatus) String() string 

// StrStatus enum declarations
const (
	StrStatusPending StrStatus = "pending"
	StrStatusRunning StrStatus = "running"
	StrStatusCompleted StrStatus = "completed"
)

func (en StrStatus) IsValid() bool
func StrStatusValues() []StrStatus
func (en StrStatus) String() string
```

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

### Enum options

Generative options are written after `)` either as a word or the parameterized value `key: value`


#### - reversed

Just add `reversed` option ater `)`

```go
// enum(pending, running, completed) reversed
type Status string
```

and generates values ​​with the name inversion

```go
// Status enum declarations
const (
	PendingStatus Status = "pending"
	RunningStatus Status = "running"
	CompletedStatus Status = "completed"
)
```