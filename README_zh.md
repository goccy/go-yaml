# Go语言的YAML支持

[![PkgGoDev](https://pkg.go.dev/badge/github.com/goccy/go-yaml)](https://pkg.go.dev/github.com/goccy/go-yaml)
![Go](https://github.com/goccy/go-yaml/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/goccy/go-yaml/branch/master/graph/badge.svg)](https://codecov.io/gh/goccy/go-yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/goccy/go-yaml)](https://goreportcard.com/report/github.com/goccy/go-yaml)

<img width="300px" src="https://user-images.githubusercontent.com/209884/67159116-64d94b80-f37b-11e9-9b28-f8379636a43c.png"></img>

**语言版本**: [English](README.md) | 中文

## 本库与go-yaml/yaml库**没有**任何关系

> [!IMPORTANT]
> 本库是从零开始开发的，用于替代[`go-yaml/yaml`](https://github.com/go-yaml/yaml)。
> 如果您正在寻找更好的YAML库，这个库应该会对您有所帮助。

# 为什么需要一个新的库？

在撰写本文时，Go语言已经存在一个事实上的标准YAML处理库：[https://github.com/go-yaml/yaml](https://github.com/go-yaml/yaml)。然而，我们认为需要一个新的YAML库，原因如下：

- 缺乏积极维护
- `go-yaml/yaml`将用C语言编写的libyaml移植到Go，因此源代码不是Go风格的
- 有很多内容无法解析
- YAML通常用于配置，通常需要包含验证。然而，`go-yaml/yaml`中的错误不够直观，很难提供有意义的验证错误
- 在创建使用YAML的工具时，有时需要对YAML进行可逆转换。然而，要对包含注释或锚点/别名的内容执行可逆转换，操作AST是唯一的选择
- 不直观的[Marshaler](https://pkg.go.dev/gopkg.in/yaml.v3#Marshaler) / [Unmarshaler](https://pkg.go.dev/gopkg.in/yaml.v3#Unmarshaler)

顺便说一下，像[ghodss/yaml](https://github.com/ghodss/yaml)和[sigs.k8s.io/yaml](https://github.com/kubernetes-sigs/yaml)这样的库也依赖于go-yaml/yaml，所以如果您使用这些库，同样的问题也适用：它们无法解析go-yaml/yaml无法解析的内容，并且继承了go-yaml/yaml的许多问题。

# 特性

- 无依赖
- 比`go-yaml/yaml`更好的解析器
  - [支持递归处理](https://github.com/apple/device-management/blob/release/docs/schema.yaml)
  - 在[YAML测试套件](https://github.com/yaml/yaml-test-suite?tab=readme-ov-file)中有更高的覆盖率
    - YAML测试套件总共包含402个测试用例，其中`gopkg.in/yaml.v3`通过了`295`个。除了通过所有这些测试用例外，`goccy/go-yaml`还成功通过了近60个额外的测试用例（2024/12/15）
    - 测试代码在[这里](https://github.com/goccy/go-yaml/blob/master/yaml_test_suite_test.go#L77)
- 易于维护和可持续性
  - 主要维护者是[@goccy](https://github.com/goccy)，但我们也在建立一个系统，与可信赖的开发者团队一起开发
  - 由于是从零开始编写的，代码对Gopher来说易于阅读
- 不仅允许使用`Encoder`/`Decoder`，还允许使用`Tokenizer`和`Parser`功能的API结构
  - [lexer.Tokenize](https://pkg.go.dev/github.com/goccy/go-yaml@v1.15.4/lexer#Tokenize)
  - [parser.Parse](https://pkg.go.dev/github.com/goccy/go-yaml@v1.15.4/parser#Parse)
- 使用YAML路径过滤、替换和合并YAML内容
- 对包含锚点、别名和注释的YAML进行可逆转换，无需使用AST
- 为原始类型和第三方库类型自定义Marshal/Unmarshal行为（[RegisterCustomMarshaler](https://pkg.go.dev/github.com/goccy/go-yaml#RegisterCustomMarshaler)，[RegisterCustomUnmarshaler](https://pkg.go.dev/github.com/goccy/go-yaml#RegisterCustomUnmarshaler)）
- 遵循`encoding/json`行为
  - 接受`json`标签。请注意，在解析YAML文档时，并非`json`标签的所有选项都有意义。如果两个标签都存在，`yaml`标签将优先
  - [json.Marshaler](https://pkg.go.dev/encoding/json#Marshaler)风格的[marshaler](https://pkg.go.dev/github.com/goccy/go-yaml#BytesMarshaler)
  - [json.Unmarshaler](https://pkg.go.dev/encoding/json#Unmarshaler)风格的[unmarshaler](https://pkg.go.dev/github.com/goccy/go-yaml#BytesUnmarshaler)
  - 使用`MarshalJSON`和`UnmarshalJSON`的选项（[UseJSONMarshaler](https://pkg.go.dev/github.com/goccy/go-yaml#UseJSONMarshaler)，[UseJSONUnmarshaler](https://pkg.go.dev/github.com/goccy/go-yaml#UseJSONUnmarshaler)）
- 错误通知的美观格式
- 与[go-playground/validator](https://github.com/go-playground/validator)结合的智能验证处理
  - [示例测试代码在这里](https://github.com/goccy/go-yaml/blob/45889c98b0a0967240eb595a1bd6896e2f575106/testdata/validate_test.go#L12)
- 允许通过锚点引用在另一个文件中声明的元素

# 用户

使用goccy/go-yaml的仓库列在这里。

- <https://github.com/goccy/go-yaml/wiki/Users>

源数据在[这里](https://github.com/goccy/go-yaml/network/dependents)。
它已经在许多仓库中使用。现在轮到您了😄

# 演练场

演练场可视化了go-yaml如何处理YAML文本。使用它来协助您的调试或问题报告。

<https://goccy.github.io/go-yaml>

# 安装

```sh
go get github.com/goccy/go-yaml
```

# 概要

## 1. 简单编码/解码

具有类似`go-yaml/yaml`的接口，使用`reflect`

```go
var v struct {
 A int
 B string
}
v.A = 1
v.B = "hello"
bytes, err := yaml.Marshal(v)
if err != nil {
 //...
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
 //...
}
```

要控制marshal/unmarshal行为，您可以使用`yaml`标签。

```go
 yml := `---
foo: 1
bar: c
`
var v struct {
 A int    `yaml:"foo"`
 B string `yaml:"bar"`
}
if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
 //...
}
```

为了方便起见，我们也接受`json`标签。请注意，在解析YAML文档时，并非`json`标签的所有选项都有意义。如果两个标签都存在，`yaml`标签将优先。

```go
 yml := `---
foo: 1
bar: c
`
var v struct {
 A int    `json:"foo"`
 B string `json:"bar"`
}
if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
 //...
}
```

对于自定义marshal/unmarshaling，实现marshaler/unmarshaler的`Bytes`或`Interface`变体。区别在于`BytesMarshaler`/`BytesUnmarshaler`的行为类似于[`encoding/json`](https://pkg.go.dev/encoding/json)，而`InterfaceMarshaler`/`InterfaceUnmarshaler`的行为类似于[`gopkg.in/yaml.v2`](https://pkg.go.dev/gopkg.in/yaml.v2)。

语义上两者是相同的，但在性能上有所不同。因为缩进在YAML中很重要，您不能简单地从Marshaler接受有效的YAML片段，并期望它在附加到父容器的序列化形式时能够工作。因此，当我们接收使用`BytesMarshaler`（返回`[]byte`）时，我们必须解码一次以确定如何在给定上下文中使其工作。如果您使用`InterfaceMarshaler`，我们可以跳过解码。

如果您重复编组复杂对象，后者在性能方面总是更好。但是，如果您只是提供一个只读一次的配置文件格式之间的选择，前者可能更容易编码。

## 2. 引用在另一个文件中声明的元素

`testdata`目录包含`anchor.yml`文件：

```shell
├── testdata
   └── anchor.yml
```

`anchor.yml`定义如下：

```yaml
a: &a
  b: 1
  c: hello
```

然后，如果将`yaml.ReferenceDirs("testdata")`选项传递给`yaml.Decoder`，
`Decoder`会尝试从`testdata`目录下的YAML文件中查找锚点定义。

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
 //...
}
fmt.Printf("%+v\n", v) // {A:{B:1 C:hello}}
```

## 3. 使用`Anchor`和`Alias`编码

### 3.1. 显式声明的`Anchor`名称和`Alias`名称

如果您想使用`anchor`，可以将其定义为结构标签。
如果为锚点指定的值是指针类型，并且找到与指针相同地址的值，则该值会自动设置为别名。
如果指定了显式别名名称，如果其值与锚点中指定的值不同，则会引发错误。

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

### 3.2. 隐式声明的`Anchor`和`Alias`名称

如果您没有显式声明锚点名称，默认行为是使用相当于`strings.ToLower($FieldName)`作为锚点的名称。
如果为锚点指定的值是指针类型，并且找到与指针相同地址的值，则该值会自动设置为别名。

```go
type T struct {
 I int
 S string
}
var v struct {
 A *T `yaml:"a,anchor"`
 B *T `yaml:"b,anchor"`
 C *T `yaml:"c"`
 D *T `yaml:"d"`
}
v.A = &T{I: 1, S: "hello"}
v.B = &T{I: 2, S: "world"}
v.C = v.A // C与A具有相同的指针地址
v.D = v.B // D与B具有相同的指针地址
bytes, err := yaml.Marshal(v)
if err != nil {
 //...
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

### 3.3 合并键和别名

合并键和别名（`<<: *alias`）可以通过嵌入带有`inline,alias`标签的结构来使用。

```go
type Person struct {
 *Person `yaml:",omitempty,inline,alias"` // 嵌入Person类型作为默认值
 Name    string `yaml:",omitempty"`
 Age     int    `yaml:",omitempty"`
}
defaultPerson := &Person{
 Name: "John Smith",
 Age:  20,
}
people := []*Person{
 {
  Person: defaultPerson, // 分配默认值
  Name:   "Ken",         // 覆盖Name属性
  Age:    10,            // 覆盖Age属性
 },
 {
  Person: defaultPerson, // 仅分配默认值
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
 //...
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

## 4. 美观格式化的错误

解析过程中产生的错误值比常规错误值有两个额外的特性。

首先，默认情况下，它们包含来自源YAML文档的错误位置的额外信息，以便更容易找到错误位置。

其次，错误消息可以选择性地着色。

如果您想精确控制输出的外观，请考虑使用`yaml.FormatError`，它接受两个布尔值来控制开启或关闭这些特性。

<img src="https://user-images.githubusercontent.com/209884/67358124-587f0980-f59a-11e9-96fc-7205aab77695.png"></img>

## 5. 使用YAMLPath

```go
yml := `
store:
  book:
    - author: john
      price: 10
    - author: ken
      price: 12
  bicycle:
    color: red
    price: 19.95
`
path, err := yaml.PathString("$.store.book[*].author")
if err != nil {
  //...
}
var authors []string
if err := path.Read(strings.NewReader(yml), &authors); err != nil {
  //...
}
fmt.Println(authors)
// [john ken]
```

### 5.1 使用YAML源代码打印自定义错误

```go
package main

import (
  "fmt"

  "github.com/goccy/go-yaml"
)

func main() {
  yml := `
a: 1
b: "hello"
`
  var v struct {
    A int
    B string
  }
  if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
    panic(err)
  }
  if v.A != 2 {
    // 使用YAML源输出错误
    path, err := yaml.PathString("$.a")
    if err != nil {
      panic(err)
    }
    source, err := path.AnnotateSource([]byte(yml), true)
    if err != nil {
      panic(err)
    }
    fmt.Printf("a值期望为2但实际为%d:\n%s\n", v.A, string(source))
  }
}
```

输出结果如下：

<img src="https://user-images.githubusercontent.com/209884/84148813-7aca8680-aa9a-11ea-8fc9-37dece2ebdac.png"></img>

# 工具

## ycat

彩色打印yaml文件

<img width="713" alt="ycat" src="https://user-images.githubusercontent.com/209884/66986084-19b00600-f0f9-11e9-9f0e-1f91eb072fe0.png">

### 安装

```sh
git clone https://github.com/goccy/go-yaml.git
cd go-yaml/cmd/ycat && go install .
```

# 开发者须知

> [!NOTE]
> 在这个项目中，我们在`testdata`目录下管理这样的测试代码，以避免在顶级`go.mod`文件中添加仅用于测试的库依赖。因此，如果您想添加使用第三方库的测试用例，请将测试代码添加到`testdata`目录中。

# 寻找赞助商

我正在为这个库寻找赞助商。这个库是作为个人项目在我的业余时间开发的。如果您希望在项目中使用这个库时获得快速响应或问题解决，请注册为[赞助商](https://github.com/sponsors/goccy)。我会尽力配合。当然，这个库是以MIT许可证开发的，所以您可以免费自由使用。

# 许可证

MIT
