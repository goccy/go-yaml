export interface GoValueType {
  response: string
  error: string
}

export interface YAMLGoFuncMap {
  decode: (v: string) => GoValueType;
  tokenize: (v: string) => GoValueType;
  parse: (v: string) => GoValueType;
}

export interface YAMLFuncMap {
  decode: (code: string) => Promise<YAMLProcessResult>
  tokenize: (code: string) => Promise<YAMLProcessResult>
  parse: (code: string) => Promise<YAMLProcessResult>
}

export enum YAMLProcessResultType {
  Decode,
  Lexer,
  ParserGroup,
  Parser,
}

export interface YAMLProcessResult {
  type: YAMLProcessResultType
  result: Token[] | string
}

export interface Token {
  type: string
  value: string
  origin: string
  error: string
  line: number
  column: number
  offset: number
}

declare function decode(v: string): GoValueType;
declare function tokenize(v: string): GoValueType;
declare function parse(v: string): GoValueType;

export const initWASM = async (path: string): Promise<YAMLGoFuncMap> => {
  const go = new Go();
  const result = await WebAssembly.instantiateStreaming(
    fetch(path),
    go.importObject
  );
  const instance = result.instance;
  go.run(instance);
  return {
    decode: decode,
    tokenize: tokenize,
    parse: parse,
  };
};