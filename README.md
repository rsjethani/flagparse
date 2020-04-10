![GitHub](https://img.shields.io/github/license/rsjethani/flagparse) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/rsjethani/flagparse) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue)](https://pkg.go.dev/github.com/rsjethani/flagparse) [![flagparse](https://circleci.com/gh/rsjethani/flagparse.svg?style=shield)](https://app.circleci.com/pipelines/github/rsjethani/flagparse) ![Codecov](https://img.shields.io/codecov/c/github/rsjethani/flagparse)


# A Powerful Commandline Flags Parser for Go

## Struct Tags Syntax
```
type <struct name> struct {
    Field1    <field type>                                               // 'flagparse' struct tag not given hence ignored
    Field2    <field type>    `flagparse:"<key=value>,<key=value>,..."   // 'flagparse' struct tag given hence parsed
    Field3    <field type>    `flagparse:"<key=value>,<key=value>,..."   // 'flagparse' struct tag given hence parsed
    ...
}
```
**PS:** The fields must be public otherwise the `reflect` package will fail to parse the struct.

## Valid Tag Keys and Values

| Key | Mandatory | Value Type (Go) | Possible Values | Default | Description |
| :---: | :---: | --- | :---: | :---: | :--- |
| `positional` | no | `N/A` | `N/A` | `N/A` | create a positional argument if given otherwise create an optional argument |
| `name` | no | string | a valid string containing alphanumeric charaters and/or '-' | struct field's name in lower case | the name to identify the flag with |
| `nargs` | no | int | a valid int, give `0` if you want a switch flag | `1` | number of values required by the flag |
| `help` | no | string | any valid string, escape `,` as `\\,`  | "" | help message for the user |

## Example

For full examples please refer to `examples/`.
