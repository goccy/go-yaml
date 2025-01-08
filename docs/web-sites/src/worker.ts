import './wasm_exec';
import { YAMLFuncMap, YAMLProcessResultType, YAMLProcessResult, Token, TokenGroup, GroupedToken, initWASM } from './YAML.ts';

const yaml = initWASM('/yaml.wasm');
const funcMap = yaml.then((v): Promise<YAMLFuncMap> => {
    const decode = (code: string): Promise<YAMLProcessResult> => {
        return new Promise((resolve) => {
            const res = v.decode(code);
            if (res.error !== undefined) {
                resolve({
                    type: YAMLProcessResultType.Decode,
                    result: res.error,
                });
                return
            }
            resolve({
                type: YAMLProcessResultType.Decode,
                result: res.response as string,
            });
        });
    };
    const tokenize = (code: string): Promise<YAMLProcessResult> => {
        return new Promise((resolve) => {
            const res = v.tokenize(code);
            if (res.error !== undefined) {
                resolve({
                    type: YAMLProcessResultType.Lexer,
                    result: res.error,
                });
                return
            }
            resolve({
                type: YAMLProcessResultType.Lexer,
                result: JSON.parse(res.response) as Token[],
            });
        });
    };
    const parseGroup = (code: string): Promise<YAMLProcessResult> => {
        return new Promise((resolve) => {
            const res = v.parseGroup(code);
            if (res.error !== undefined) {
                resolve({
                    type: YAMLProcessResultType.ParserGroup,
                    result: res.error,
                });
                return
            }
            resolve({
                type: YAMLProcessResultType.ParserGroup,
                result: JSON.parse(res.response) as GroupedToken[],
            });           
        });
    };
    const parse = (code: string): Promise<YAMLProcessResult> => {
        return new Promise((resolve) => {
            const res = v.parse(code);
            if (res.error !== undefined) {
                resolve({
                    type: YAMLProcessResultType.Parser,
                    result: res.error,
                });
                return
            }
            resolve({
                type: YAMLProcessResultType.Parser,
                result: res.response as string,
            });
        });
    };
    return new Promise((resolve) => {
        resolve({
            decode: decode,
            tokenize: tokenize,
            parseGroup: parseGroup,
            parse: parse,
        });
    })
});

self.addEventListener('message', (e) => {
    const code = e.data as string;
    funcMap.then((v) => {
        const decode = v.decode(code);
        const tokenize = v.tokenize(code);
        const parseGroup = v.parseGroup(code);
        const parse = v.parse(code);
        Promise.all([decode, tokenize, parseGroup, parse]).then((value) => {
            self.postMessage(value);
        });
    })
});

export default {}