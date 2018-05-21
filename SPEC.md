# json canonical form

All JSON is stored in UTF-8.

## Value

A value is the top-level element type.

```
element = string | number | object | array | "true" | "false" | "null"
```

## Strings

```
string-element = unicode value except \ and "
               | "\\", "\\"
               | "\\", "\""
string = "\"" { string-element } "\""
```

## Numbers

Numbers correspond to the JSON definition of numbers, but are implicitly defined
into two separate categories, described below.

```
number = integer | floating-point
```

### Integers

Integers are numbers inside the range `[-(2**53)+1, (2**53)-1]` that have a zero
fractional part when written without an exponent.

An integer is emitted as a number with no fractional part and no exponent,
without leading zeroes (with the exception of zero):

```
integer = 0
        | [ "-" ], nonzero-digit, { digit }
```

### Floating point numbers

A floating point value is any number outside of the range `[-(2**53)+1, (2**53)-1]`,
or any value with a nonzero fractional part when written without an exponent.

A floating point number is emitted as a single nonzero digit, followed by an
optional fractional part, followed by an exponent. The fractional part must end
with a nonzero digit:

```
trailing-nonzero = nonzero-digit
                 | digit, trailing-nonzero

floating-point = nonzero-digit, [ ".", trailing-nonzero ], "E", integer
```

Note that, although the EBNF definition of `floating-point` allows numbers such
as `1E0`, they are not valid due to the constraint described at the top of this
section.

## Objects

Objects are printed as JSON objects, where the key-value pairs are ordered
lexicographically on keys.

```
pair = string, ":", element
pairs = pair
      | pair, ",", pairs
object = "{", [ pairs ], "}"
```

## Arrays

Arrays consist of zero or more elements, separated by a comma.

```
elements = element
         | element, ",", elements
array = "[", [ elements ], "]"
```