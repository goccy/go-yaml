# YAML support for the Go language

# Why?

I know [https://github.com/go-yaml/yaml]() . It is the very famous library .
But I want to use the following features no in that library .

- Beautiful syntax or validation error notification
- Manipulate Abstract Syntax Tree for YAML with Go
- Support `Anchor` and `Alias` to Marshaler
- Reference to Anchor defined by the other file

# Status

- This library should be considered alpha grade. API may still change.

# Features

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