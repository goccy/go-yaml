import './wasm_exec';
import { YAMLFuncMap, YAMLProcessResultType, YAMLProcessResult, Token, GroupedToken, initWASM } from './YAML.ts';

const yaml = initWASM();
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
    const data = e.data as {
        code: string
        tabIndex: number
    }
    funcMap.then((v) => {
        switch (data.tabIndex) {
            case 0:
                v.decode(data.code).then((value) => {
                    self.postMessage(value);
                })
                break;
            case 1:
                v.parse(data.code).then((value) => {
                    self.postMessage(value);
                })
                break;
            case 2:
                v.parseGroup(data.code).then((value) => {
                    self.postMessage(value);
                })
                break;
            case 3:
                v.tokenize(data.code).then((value) => {
                    self.postMessage(value);
                })
                break;
            default:
                break;
        }
    })
});

export default {}