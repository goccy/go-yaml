## v1.9.1 - 2021-07-20

### Fix DecodeFromNode ( #237 )

- Fix YAML handling where anchor exists

## v1.9.0 - 2021-07-19

### New features

- Support encoding of comment node ( #233 )
- Support `yaml.NodeToValue(ast.Node, interface{}, ...DecodeOption) error` ( #236 )
  - Can convert a AST node to a value directly

### Fix decoder for comment

- Fix parsing of literal with comment ( #234 )

### Rename API ( #235 )

- Rename `MarshalWithContext` to `MarshalContext`
- Rename `UnmarshalWithContext` to `UnmarshalContext`

## v1.8.10 - 2021-07-02

### Fixed bugs

- Fix searching anchor by alias name ( #212 )
- Fixing Issue 186, scanner should account for newline characters when processing multi-line text. Without this source annotations line/column number (for this and all subsequent tokens) is inconsistent with plain text editors. e.g. https://github.com/goccy/go-yaml/issues/186. This addresses the issue specifically for single and double quote text only. ( #210 )
- Add error for unterminated flow mapping node ( #213 )
- Handle missing required field validation ( #221 )
- Nicely format unexpected node type errors ( #229 )
- Support to encode map which has defined type key ( #231 )

### New features

- Support sequence indentation by EncodeOption ( #232 )

## v1.8.9 - 2021-03-01

### Fixed bugs

- Fix origin buffer for DocumentHeader and DocumentEnd and Directive
- Fix origin buffer for anchor value
- Fix syntax error about map value
- Fix parsing MergeKey ('<<') characters
- Fix encoding of float value
- Fix incorrect column annotation when single or double quotes are used

### New features

- Support to encode/decode of ast.Node directly
