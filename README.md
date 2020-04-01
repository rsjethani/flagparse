# A Powerful Argument Parser for Go

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
| `type` | no | string | `pos`/`opt`/`switch` | `opt` | create a positional argument if given otherwise create an optional argument |
| `name` | no | string | a valid string containing alphanumeric charaters and/or '-' | struct field's name in lower case | the name to identify the argument with |
| `nargs` | no | int | a valid int | `1` if `type=pos\|opt`, `0` if `type=switch` | number of values required by the argument |
| `help` | no | string | any valid string, escape `,` as `\\,`  | "" | help message for the user |

## Example

For full examples please refer to `examples/`.
