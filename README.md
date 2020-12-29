# ifshort
[![Go Report Card](https://goreportcard.com/badge/github.com/esimonov/ifshort)](https://goreportcard.com/report/github.com/esimonov/ifshort)
<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-96%25-brightgreen.svg?longCache=true&style=flat)</a>

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
	otherFunc3(err) // Variable referenced in a call outside the if-statement.
}
```

```go
func someFunc() {
	sliceOfStrings := getSliceOfStrings()

	if sliceOfStrings == 1 && sliceOfStrings[0] == "" {
		return
	}

	for _, s := range sliceOfStrings { // Variable referenced in a loop.
		// do stuff
	}
}
```

etc.