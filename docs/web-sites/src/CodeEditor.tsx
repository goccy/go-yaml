import { useState, useRef, useEffect } from 'react';
import { editor } from 'monaco-editor'
import MonacoEditor from '@monaco-editor/react';
import { Box, Tabs, Tab, Tooltip, TooltipProps, tooltipClasses, styled } from '@mui/material';
import Grid from '@mui/material/Grid2';
import { useXTerm } from 'react-xtermjs';
import { FitAddon } from '@xterm/addon-fit';
import { Token, YAMLProcessResult, YAMLProcessResultType } from './YAML.ts';
import yamlWorker from "./worker?worker";
import '@xterm/xterm/css/xterm.css';

function TabPanel(props: { children: any, value: number, index: number }) {
    const { children, value, index, ...other } = props;

    return (
        <div
            role="tabpanel"
            hidden={value !== index}
            id={`simple-tabpanel-${index}`}
            aria-labelledby={`simple-tab-${index}`}
            {...other}
        >
            {value === index && (
                <Box p={3}>
                    {children}
                </Box>
            )}
        </div>
    );
}

function a11yProps(index: number) {
    return {
        id: `simple-tab-${index}`,
        'aria-controls': `simple-tabpanel-${index}`,
    };
}

const TerminalComponent = (v: any) => {
    const { instance, ref } = useXTerm()
    const fitAddon = new FitAddon()
    const out = v.out;

    useEffect(() => {
        if (instance === null) {
            return
        }
        // Load the fit addon
        instance.loadAddon(fitAddon)
        instance.options.cursorInactiveStyle = 'none';
        instance.options.cursorStyle = 'bar';
        instance.options.letterSpacing = 4;
        instance.options.fontFamily = 'monospace';
        instance.options.fontSize = 16;
        instance.options.convertEol = true;

        fitAddon.fit();
        instance.clear();
        instance.writeln(out);
        const handleResize = () => fitAddon.fit()

        window.addEventListener('resize', handleResize)
        return () => {
            window.removeEventListener('resize', handleResize)
        }
    }, [ref, instance, out]);

    return <div ref={ref} style={{ height: 400, width: '100%', textAlign: 'left' }} />
}


interface RenderedToken {
    origins: string[]
    color: string
    prop: string
}

const CustomTooltip = styled(({ className, ...props }: TooltipProps) => (
    <Tooltip {...props} classes={{ popper: className }} />
))({
    [`& .${tooltipClasses.tooltip}`]: {
        fontSize: 16,
        whiteSpace: 'pre-wrap',
    },
});

const Lexer = (v: any) => {
    if (!v.tokens) {
        return <Box sx={{
            backgroundColor: '#001435',
            height: 400,
        }}></Box>
    }
    const tokens = v.tokens as Token[];
    const tks: RenderedToken[] = (tokens as Token[]).map((token, tokenIndex) => {
        const orgs = token.origin.split('\n').map((v, idx) => {
            if (idx > 0) {
                return ["\n", v];
            }
            return v;
        });
        const color = () => {
            if (tokens.length > tokenIndex + 1) {
                const nextToken = v.tokens[tokenIndex + 1]
                if (nextToken.type === 'MappingValue') {
                    return '#008b8b';
                }
            }
            switch (token.type) {
                case "String":
                    return '#ff7f50';
            }
            return 'black';
        };
        const prop = `type:   ${token.type}
origin: ${JSON.stringify(token.origin)}
value:  ${JSON.stringify(token.value)}
line:   ${token.line}
column: ${token.column}
`
        return {
            origins: orgs.flat().filter((v) => { return v != "" }),
            color: color(),
            prop: prop,
        }
    })
    return (
        <>
            <Box sx={{
                textAlign: 'left',
                backgroundColor: '#001435',
                fontSize: 16,
                fontWeight: 'bold',
                fontFamily: 'monospace',
                paddingLeft: '1em',
                height: 400,
            }}>
                {
                    tks.map((tk) => {
                        return (
                            <CustomTooltip title={tk.prop}>
                                <Box component="span" sx={{
                                    backgroundColor: 'white',
                                    paddingRight: 1,
                                    marginRight: 1,
                                    borderRadius: 1,
                                    color: tk.color,
                                }}>
                                    {
                                        tk.origins.map((v) => {
                                            return (
                                                <Box component="span"
                                                    sx={{
                                                        whiteSpace: "pre-wrap",
                                                        paddingLeft: 1,
                                                    }}
                                                >{v}</Box>
                                            )
                                        })
                                    }
                                </Box>
                            </CustomTooltip>
                        )
                    })
                }
            </Box>
        </>
    )
};

const AST = (v: any) => {
    if (!v?.svg) {
      return <></>
    }
    const parser = new DOMParser()
    const dom = parser.parseFromString(v.svg, 'text/xml');
    const g = dom.getElementById('graph0');
    if (!g) {
        return <></>
    }
    const viewBox = g.parentElement!.getAttribute('viewBox')!;
    return (
        <>
            <svg width={'100%'} height={'100%'} viewBox={viewBox} dangerouslySetInnerHTML={{__html: g.outerHTML}}></svg>
        </>
    )
}

function CodeEditor() {
    const editorRef = useRef<editor.IStandaloneCodeEditor | null>(null);
    const workerRef = useRef<Worker | null>(null);
    const [tokens, setTokens] = useState<Token[]>([]);
    const [out, setOut] = useState<string>('');
    const [svg, setSvg] = useState<string>('');
    useEffect(() => {
        workerRef.current = new yamlWorker();
        workerRef.current.onmessage = (event) => {
            const data = event.data as YAMLProcessResult[];
            data.forEach((v) => {
                console.log(v);
                switch (v.type) {
                    case YAMLProcessResultType.Decode:
                        setOut(v.result as string);
                        break;
                    case YAMLProcessResultType.Lexer:
                        if (typeof v.result === 'string') {
                            console.error(v.result);
                        } else {
                            setTokens(v.result);
                        }
                        break;
                    case YAMLProcessResultType.Parser:
                        setSvg(v.result as string);
                        break;
                    default:
                        break;
                }
            });
        };
        return () => {
            workerRef.current?.terminate();
        };
    }, []);
    const onChange = () => {
        const code = editorRef?.current?.getValue()!;
        console.log('code = ', code);
        if (workerRef.current) {
            workerRef.current.postMessage(code);
        }
    };
    const onMount = (editor: editor.IStandaloneCodeEditor, monaco) => {
        editorRef.current = editor;
    };
    const [value, setValue] = useState(0);
    const handleChange = (event: React.SyntheticEvent, newValue: number) => {
        setValue(newValue);
    };
    return (
        <>
            <Grid container>
                <Grid marginTop={10} size={{ xs: 6, md: 6 }}>
                    <MonacoEditor
                        height={400}
                        language="yaml"
                        theme="vs-dark"
                        value={'foo: bar'}
                        options={{
                            fontSize: 16,
                            selectOnLineNumbers: true,
                            renderWhitespace: 'all',
                            autoIndent: 'none',
                        }}
                        onChange={onChange}
                        onMount={onMount}
                    />
                </Grid>
                <Grid size={{ xs: 6, md: 6 }}>
                    <Box marginTop={1}>
                        <Tabs
                            textColor='secondary'
                            indicatorColor='secondary'
                            value={value}
                            onChange={handleChange}
                            variant="scrollable"
                            scrollButtons="auto"
                            aria-label="tabs">
                            <Tab style={{ marginLeft: 20 }} label="Console" {...a11yProps(0)} />
                            <Tab style={{ marginLeft: 20 }} label="Lexer" {...a11yProps(1)} />
                            <Tab style={{ marginLeft: 20 }} label="Parser(Grouping)" {...a11yProps(2)} />
                            <Tab style={{ marginLeft: 20 }} label="Parser(AST)" {...a11yProps(3)} />
                        </Tabs>
                        <TabPanel value={value} index={0}>
                            <TerminalComponent out={out} />
                        </TabPanel>
                        <TabPanel value={value} index={1}>
                            <Lexer tokens={tokens}></Lexer>
                        </TabPanel>
                        <TabPanel value={value} index={2}>
                            Parser(Grouping)
                        </TabPanel>
                        <TabPanel value={value} index={3}>
                            <AST svg={svg}></AST>
                        </TabPanel>
                    </Box>
                </Grid>
            </Grid>
        </>
    )
}

export default CodeEditor;
