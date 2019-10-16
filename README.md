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

# Install

```
go get -u github.com/goccy/go-yaml
```

# License

MIT
