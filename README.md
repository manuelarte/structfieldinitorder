# Struct Field Init Order

[![Go](https://github.com/manuelarte/structfieldinitorder/actions/workflows/go.yml/badge.svg)](https://github.com/manuelarte/structfieldinitorder/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/manuelarte/structfieldinitorder)](https://goreportcard.com/report/github.com/manuelarte/structfieldinitorder)
[![version](https://img.shields.io/github/v/release/manuelarte/structfieldinitorder)](https://img.shields.io/github/v/release/manuelarte/structfieldinitorder)

This linter checks whether when a struct instantiates, the fields order follows the same order as in the struct declaration.

## ⬇️  Getting Started

### As a Binary

To install it run:

```bash
go install github.com/manuelarte/structfieldinitorder/cmd/structfieldinitorder@latest
```

And then use it as:

```bash
structfieldinitorder ./...
```

### As a golangci-lint module plugin

Take a look at the example [golangciplugin](./examples/golangciplugin) on how to run this linter as a plugin.
As a summary ([more info](https://golangci-lint.run/plugins/module-plugins)):

+ Create `.custom-gcl.yml`.
+ Run `golangci-lint custom`.
+ Define the plugin inside the `linters.settings.custom` section with the type `module`.
+ Run the resulting custom binary of golangci-lint

## 🚀 Features

Check fields order:

```go
type Person struct {
  Name      string
  Surname   string
  Birthdate time.Time
}
```

<table>
<thead><tr><th>❌ Bad</th><th>✅ Good</th></tr></thead>
<tbody>
<tr><td>

```go
// ❌ Not following Name, Surname, Birthdate 
var Me = Person {
  Surname: "Doe",
  Name: "John",
  Birthdate: time.Now(),
}
```

</td><td>

```go
// ✅ Name, Surname, Birthdate
var Me = Person {
  Name: "John",
  Surname: "Doe",
  Birthdate: time.Now(),
}
```

</td></tr>

</tbody>
</table>
