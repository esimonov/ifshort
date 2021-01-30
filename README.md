# ifshort
[![Go Report Card](https://goreportcard.com/badge/github.com/esimonov/ifshort)](https://goreportcard.com/report/github.com/esimonov/ifshort)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-97%25-brightgreen.svg?longCache=true&style=flat)</a>

Go linter that checks if your code uses short syntax for `if`-statements whenever possible.

For example, it will suggest changing this code:

```go
func someFunc(k string, m map[string]interface{}) {
	_, ok := m[k]
	if !ok {
		return
	}

	err := otherFunc1()
	if err != nil {
		otherFunc2(err)
	}
}
```
to this:
```go
func someFunc(k string, m map[string]interface{}) {
	if _, ok := m[k]; !ok {
		return
	}

	if err := otherFunc1(); err != nil {
		otherFunc2(err)
	}
}
```

At the same time, it won't suggest any changes if a variable referenced in the `if`-statement also occurs in other places. E.g, linter won't complain about either of the snippents below:

```go
func someFunc() {
	err := otherFunc1()
	if err != nil {
		otherFunc2(err)
	}
	otherFunc3(err) // Variable referenced in a call outside if-statement.
}
```

```go
func someFunc() interface{} {
	v := getValue()
	if v != nil {
		otherFunc(v)
	}
	return v // Variable referenced in return statement.
}
```

```go
func someFunc() interface{} {
	a, b := getTwoValues()
	if a != nil {
		return a
	}
	return b // Variables a and b are used in different statements.
}
```

```go
func someFunc() {
	v := &someStruct{}

	if v != nil { // cannot be `if v := &someStruct{}; v != nil
		return
	}
}
```
etc.

## Usage

```shell
usage: ifshort [--max-decl-chars {integer}] [--max-decl-lines {integer}] [INPUT]

positional arguments:
  INPUT

options:
  --max-decl-chars
        maximum length of variable declaration measured in number of characters, after which the linter won't suggest using short syntax. (default 30)
  --max-decl-lines
        maximum length of variable declaration measured in number of lines, after which the linter won't suggest using short syntax.  (default 1)
		Has precedence over max-decl-chars.
```

Example usage to check only the variables whose declaration takes no more than 50 characters:

`ifshort --max-decl-chars 50 path/to/myproject`.

With this configuration, `ifshort` won't suggest using short syntax on line 3:

```go
1 func someFunc() {
2	v := getValue("Long long long declaration, linter shouldn't force short syntax for it.") // More than 50 characters long.
3	if v != nil {
4		otherFunc1(v)
5	}
6 }
```

Example usage to check only the variables whose declaration takes no more than 2 lines:

`ifshort --max-decl-lines 2 path/to/myproject`.

```go
func someFunc() {
	v1 := getValue("firstLine",
		    "secondLine") // This declaration will be checked, and short syntax suggested.
	if v1 != nil {
		someFunc(v1)
	}

	v2 := getValue("firstLine",
			"secondLine",
			"thirdLine") // This declaration won't be checked.
	if v2 != nil {
		someFunc(v2)
	}
}
```