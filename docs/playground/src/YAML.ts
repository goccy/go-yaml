import init from './yaml.wasm?init';

export interface GoValueType {
  response: string
  error: string
}

export interface YAMLGoFuncMap {
  decode: (v: string) => GoValueType;
  tokenize: (v: string) => GoValueType;
  parseGroup: (v: string) => GoValueType;
  parse: (v: string) => GoValueType;
}

export interface YAMLFuncMap {
  decode: (code: string) => Promise<YAMLProcessResult>
  tokenize: (code: string) => Promise<YAMLProcessResult>
  parseGroup: (code: string) => Promise<YAMLProcessResult>
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
  result: Token[] | GroupedToken[] | string
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

export interface GroupedToken {
  token: Token
  group: TokenGroup
  lineComment: Token
}

export interface TokenGroup {
  type: string
  tokens: GroupedToken[]
}

declare function decode(v: string): GoValueType;
declare function tokenize(v: string): GoValueType;
declare function parseGroup(v: string): GoValueType;
declare function parse(v: string): GoValueType;

export const initWASM = async (): Promise<YAMLGoFuncMap> => {
  const go = new Go();
  return new Promise(resolve => {
    init(
      go.importObject,
    ).then((instance) => {
      go.run(instance);
      resolve({
        decode: decode,
        tokenize: tokenize,
        parseGroup: parseGroup,
        parse: parse,
      });
    });
  });
};