# Struct Field Init Order

[![Go](https://github.com/manuelarte/structfieldinitorder/actions/workflows/go.yml/badge.svg)](https://github.com/manuelarte/structfieldinitorder/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/manuelarte/structfieldinitorder)](https://goreportcard.com/report/github.com/manuelarte/structfieldinitorder)
[![version](https://img.shields.io/github/v/release/manuelarte/structfieldinitorder)](https://img.shields.io/github/v/release/manuelarte/structfieldinitorder)

This linter checks whether when a struct instantiates, the fields order follows the same order as in the struct declaration.

## ⬇️  Getting Started

To install it run:

```bash
go install github.com/manuelarte/structfieldinitorder@latest
```

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
