# YAML support for the Go language

[![GoDoc](https://godoc.org/github.com/goccy/go-yaml?status.svg)](https://godoc.org/github.com/goccy/go-yaml)
[![CircleCI](https://circleci.com/gh/goccy/go-yaml.svg?style=shield)](https://circleci.com/gh/goccy/go-yaml)
[![codecov](https://codecov.io/gh/goccy/go-yaml/branch/master/graph/badge.svg)](https://codecov.io/gh/goccy/go-yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/goccy/go-yaml)](https://goreportcard.com/report/github.com/goccy/go-yaml)

<img width="300px" src="https://user-images.githubusercontent.com/209884/67159116-64d94b80-f37b-11e9-9b28-f8379636a43c.png"></img>

# Why a new library?

As of this writing, there already exists a defacto standard library for YAML processing Go: [https://github.com/go-yaml/yaml](https://github.com/go-yaml/yaml). However we feel that some features are lacking, namely:

- Pretty format for error notifacations
- Directly manipulate the YAML abstract syntax tree
- Support `Anchor` and `Alias` when marshaling
- Allow referencing elements declared in another file via anchors

# Status

- This library should be considered alpha grade. API may still change.

# Features ( Goals )

- Pretty format for error notifacations
- Support `Scanner` or `Lexer` or `Parser` as public API
- Support `Anchor` and `Alias` to Marshaler
- Allow referencing elements declared in another file via anchors

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

## 2. Reference elements in declared in another file

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

If you want to use `anchor` or `alias`, you can define it as a struct tag.

```go
type T struct {
  A int
  B string
}
var v struct {
  C *T `yaml:"c,anchor=x"`
  D *T `yaml:"d,alias=x"`
}
v.C = &T{A: 1, B: "hello"}
v.D = v.C
bytes, err := yaml.Marshal(v)
if err != nil {
  panic(err)
}
fmt.Println(string(bytes))
/*
c: &x
  a: 1
  b: hello
d: *x
*/
```

### 3.2. Implicitly declared `Anchor` and `Alias` names

If you do not explicitly declare the anchor name, the default behavior is to
use the equivalent of `strings.ToLower($FieldName)` as the name of the anchor.

If you do not explicitly declare the alias name AND the value is a pointer
to another element, we look up the anchor name by finding out which anchor
field the value is assigned to by looking up its pointer address.

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

### 3.3 MergeKey and Alias

Merge key and alias ( `<<: *alias` ) can be used by embedding a structure with the `inline,alias` tag .

```go
type Person struct {
	*Person `yaml:",omitempty,inline,alias"` // embed Person type for default value
	Name    string `yaml:",omitempty"`
	Age     int    `yaml:",omitempty"`
}
defaultPerson := &Person{
	Name: "John Smith",
	Age:  20,
}
people := []*Person{
	{
		Person: defaultPerson, // assign default value
		Name:   "Ken",         // override Name property
		Age:    10,            // override Age property
	},
	{
		Person: defaultPerson, // assign default value only
	},
}
var doc struct {
	Default *Person   `yaml:"default,anchor"`
	People  []*Person `yaml:"people"`
}
doc.Default = defaultPerson
doc.People = people
bytes, err := yaml.Marshal(doc)
if err != nil {
	...
}
fmt.Println(string(bytes))
/*
default: &default
  name: John Smith
  age: 20
people:
- <<: *default
  name: Ken
  age: 10
- <<: *default
*/
```

# 4. Pretty Formatted Errors

<img src="https://user-images.githubusercontent.com/209884/67358124-587f0980-f59a-11e9-96fc-7205aab77695.png"></img>

If parser receives invalid YAML, the error is printed along with highlighted source code.

※ If you do not want to source code in your error message (colored or otherwise), set `errors.ColoredErr = false` or `errors.WithSourceCode = false` .

# Installation

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
