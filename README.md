# YAML support for the Go language

[![GoDoc](https://godoc.org/github.com/goccy/go-yaml?status.svg)](https://godoc.org/github.com/goccy/go-yaml)
[![CircleCI](https://circleci.com/gh/goccy/go-yaml.svg?style=shield)](https://circleci.com/gh/goccy/go-yaml)
[![codecov](https://codecov.io/gh/goccy/go-yaml/branch/master/graph/badge.svg)](https://codecov.io/gh/goccy/go-yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/goccy/go-yaml)](https://goreportcard.com/report/github.com/goccy/go-yaml)

# Why?

I know [https://github.com/go-yaml/yaml](https://github.com/go-yaml/yaml) . It is the very famous library .
But I want to use the following features no in that library .

- Beautiful syntax or validation error notification
- Manipulate Abstract Syntax Tree for YAML with Go
- Support `Anchor` and `Alias` to Marshaler
- Reference to Anchor defined by the other file

# Status

- This library should be considered alpha grade. API may still change.

# Features ( Goals )

- Beautiful syntax or validation error notification
- Support `Scanner` or `Lexer` or `Parser` as public API
- Support `Anchor` and `Alias` to Marshaler
- Reference to `Anchor` defined by the other file

# Synopsis

## Simple Encode/Decode

Support compatible interface to `go-yaml/yaml` by using `reflect`

```go
var v struct {
	A int
	B string
}
v.A = 1
v.B = "hello"
bytes, err := yaml.Marshal(v)
if err != nil {
	...
}
fmt.Println(string(bytes)) // "a: 1\nb: hello\n"
```

```go
	yml := `
%YAML 1.2
---
a: 1
b: c
`
var v struct {
	A int
	B string
}
if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
	...
}
```

## Reference to `Anchor` defined by the other file

`testdata` directory includes `anchor.yml` file

```shell
├── testdata
   └── anchor.yml
```

And `anchor.yml` is defined the following.

```yaml
a: &a
  b: 1
  c: hello
```

Then, if `yaml.ReferenceDirs("testdata")` option passed to `yaml.Decoder`, 
 `Decoder` try to find anchor definition from YAML files the under `testdata` directory.
 
```go
buf := bytes.NewBufferString("a: *a\n")
dec := yaml.NewDecoder(buf, yaml.ReferenceDirs("testdata"))
var v struct {
	A struct {
		B int
		C string
	}
}
if err := dec.Decode(&v); err != nil {
	...
}
fmt.Printf("%+v\n", v) // {A:{B:1 C:hello}}
```

# Install

```
$ go get -u github.com/goccy/go-yaml
```

# Tools

## ycat 

print yaml file with color

<img width="713" alt="ycat" src="https://user-images.githubusercontent.com/209884/66986084-19b00600-f0f9-11e9-9f0e-1f91eb072fe0.png">

### Install

```
$ go get -u github.com/goccy/go-yaml/cmd/ycat
```

# License

MIT
