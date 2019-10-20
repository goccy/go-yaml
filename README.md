# YAML support for the Go language

[![GoDoc](https://godoc.org/github.com/goccy/go-yaml?status.svg)](https://godoc.org/github.com/goccy/go-yaml)
[![CircleCI](https://circleci.com/gh/goccy/go-yaml.svg?style=shield)](https://circleci.com/gh/goccy/go-yaml)
[![codecov](https://codecov.io/gh/goccy/go-yaml/branch/master/graph/badge.svg)](https://codecov.io/gh/goccy/go-yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/goccy/go-yaml)](https://goreportcard.com/report/github.com/goccy/go-yaml)

<img width="300px" src="https://user-images.githubusercontent.com/209884/67159116-64d94b80-f37b-11e9-9b28-f8379636a43c.png"></img>

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

## 1. Simple Encode/Decode

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

## 2. Reference to `Anchor` defined by the other file

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

## 3. Encode with `Anchor` and `Alias`

### 3.1. Explicitly declaration `Anchor` name and `Alias` name

If you want to use `anchor` or `alias`,
it can define as tag in struct.

```go
type T struct {
	A int
	B string
}
var v struct {
	A *T `yaml:"a,anchor=c"`
	B *T `yaml:"b,alias=c"`
}
v.A = &T{A: 1, B: "hello"}
v.B = v.A
bytes, err := yaml.Marshal(v)
if err != nil {
	...
}
fmt.Printf("%s\n", string(bytes)) // "a: &c\n  a: 1\n  b: hello\nb: *c\n"
```

### 3.2. Implicitly declaration `Anchor` name and `Alias` name

If omitted anchor name, assigned default rendering name ( `strings.ToLower(FieldName)` ) as anchor name.
If omitted alias name and it's field type is **pointer** type, assigned anchor name automatically from same pointer address.

```go
type T struct {
	I int
	S string
}
var v struct {
	A *T `yaml:"a,anchor"`
	B *T `yaml:"b,anchor"`
	C *T `yaml:"c,alias"`
	D *T `yaml:"d,alias"`
}
v.A = &T{I: 1, S: "hello"}
v.B = &T{I: 2, S: "world"}
v.C = v.A // C has same pointer address to A
v.D = v.B // D has same pointer address to B
bytes, err := yaml.Marshal(v)
if err != nil {
	...
}
fmt.Println(string(bytes)) 
/*
a: &a
  i: 1
  s: hello
b: &b
  i: 2
  s: world
c: *a
d: *b
*/
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
