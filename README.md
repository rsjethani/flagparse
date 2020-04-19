![GitHub](https://img.shields.io/github/license/rsjethani/flagparse?color=blue) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/rsjethani/flagparse) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue)](https://pkg.go.dev/github.com/rsjethani/flagparse) [![Go Report Card](https://goreportcard.com/badge/github.com/rsjethani/flagparse)](https://goreportcard.com/report/github.com/rsjethani/flagparse) [![flagparse](https://circleci.com/gh/rsjethani/flagparse.svg?style=shield)](https://app.circleci.com/pipelines/github/rsjethani/flagparse) ![Codecov](https://img.shields.io/codecov/c/github/rsjethani/flagparse)


# A Powerful Commandline Flags Parser for Go

Inspired by the Go's standard [flag package](https://golang.org/pkg/flag) and [Python's argparse module](https://docs.python.org/3/library/argparse.html).

## Features
- Support for both positional and optional flags.
- The flags can take multiple arguments from 0 to unlimited.
- Built-in support for common types like bool, int, uint, etc. and their slice counterparts.
- A simple interface similar to the standard flag package for using your own types with the package.
- Concisely describe your flags using the simple struct tags based syntax. The API based approach is also supported.

### Future Enhancements
- Support for mutually exclusive set of flags.
- Constraint flag arguments to a pre-defined set of choices.
- Support sub-commands like `git commit`, `docker container ls` etc.


## Installation
`$ go get github.com/rsjethani/flagparse`

## Usage

Please visit [Go Docs Page](https://pkg.go.dev/github.com/rsjethani/flagparse) for detailed documentation and examples.

## PS: The package in still under development so use with caution.
