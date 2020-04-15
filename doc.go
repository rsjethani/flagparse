/*
Package flagparse provides a powerful and feature-rich command line flags parser.

Usage

There are two ways in which you can you use this package to parse command line arguments.

Using the package API to create the various objects yourself:

This approach is similar to ``flag'' package's approach. We create flags from existing varialbes,
add these flags to a flagset then call Parse() on the flagset. This approach is much more verbose
and should be used when you have only a handful of flags to deal with. See the ``APIApproach''
example for basic steps.

Using struct tags similar to standard library's json package:

In this approach all the variables from which we want to create flags must be struct fields with
appropriate tag. The package then uses reflect package's ability to parse the struct tags and
creates the flagset object for you. Then you simply need to call Parse() on the flagset. This
approach is more concise and preferred if your application requires a lot of flags to be processed.
See the ``StructTagApproach'' example for basic steps or read on for more details on struct tag
syntax.

Struct Tag Syntax

A struct field must be exported and tagged with ``flagparse'' in order to be parsed succcessfully.
General syntax:

	type <struct name> struct {
		Field1    <field type>                                               	 // field exported but not tagged hence ignored
		field2    <field type>    `flagparse:"[<key=value>],[<key=value>],..."   // field tagged but is un-exported hence ignored
		Field3    <field type>    `flagparse:"[<key=value>],[<key=value>],..."   // field exported and tagged hence parsed
		Field4    <field type>    `flagparse:"[<key=value>],[<key=value>],..."   // field exported and tagged hence parsed
		...
	}

The tag value is a sequence of key-value pairs where each pair should separated by a ``,''.
The key and its value are themselves separated by a ``=''. What goes in the value part depends
on the key it is being assigned to. To use ``,'' in the value part escape it with a pair of
back-slashes, see examples below. Following are the valid keys:

``name''

Specifies the name to use for the flag. The value can be a sequence of upper/lower case
alpha-numeric characters and a``-''. If omitted, the field's name in lower-case is used as value.
You don't have to prefix the value with ``--'' when creating an optional flag, the will be done
automatically. More precisely the value of FlagSet.OptPrefix is used as the prefix, ``--'' being
the default.

``positional''

By default an optional flag is created from the struct field but if this key is given then a
positional flag is created. This key does not take any value.

``usage''

Specifies the usage string for the flag. It can be any valid string. If omitted, then an empty
string is used as value.

``nargs''

Specifies the number of arguments the flag requires. The value can a vaild integer. If omitted, then
``1'' is used as value.

A negative integer means unlimited number of arguments. Both optional and positional flags can have
unlimited arguments with following rules: 1) Only one flag positional or optional can have unlimited
arguments in a flagset. 2) For positional flags only the last one specified can have unlimited
arguments.

The value ``0'' is also a bit special. When specified for a positional flag it results in error
since positional flag must have at least one argument. For an optional flag specifying ``0'' means
the flag doesn't require any arguments i.e. it is essentially a switch.


Some examples:

	type <struct name> struct {
		// an optional flag with name="field1",nargs="1",usage=""
		Field1  int  `flagparse:""`

		// an optional flag with name="--loc-data",nargs="3",usage="hello, world!"
		Field2  []int  `flagparse:"name=loc-data,nargs=3,usage=hello\\, world!"`

		// a positional flag with name="field2",nargs="1",usage=""
		Field3  int  `flagparse:"positional"`

		// a positional flag with name="f4" and unlimited arguments
		Field4  []int  `flagparse:"positional,name=f4,nargs=-1"`

		// Error: nargs cannot be 0 for a positional flag
		Field5  int  `flagparse:"positional,nargs=0"`

		// a switch flag with name="--f6"
		Field6  int  `flagparse:"name=f6,nargs=0"`
	}


User Defined Types

The package provides support for common built-in types but it is easy to extend this support
to other types including your own by simply implementing the Value interface. Please see the
interface's documentation for more details.
*/
package flagparse
