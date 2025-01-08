import { useState, useRef, useEffect } from 'react';
import { editor } from 'monaco-editor'
import MonacoEditor from '@monaco-editor/react';
import { Box, Tabs, Tab, Tooltip, TooltipProps, tooltipClasses, styled } from '@mui/material';
import Grid from '@mui/material/Grid2';
import { useXTerm } from 'react-xtermjs';
import { FitAddon } from '@xterm/addon-fit';
import { Token, TokenGroup, GroupedToken, YAMLProcessResult, YAMLProcessResultType } from './YAML.ts';
import yamlWorker from "./worker?worker";
import '@xterm/xterm/css/xterm.css';
import { ArrowBackIosNew } from '@mui/icons-material'

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

const CustomTooltip = styled(({ className, ...props }: TooltipProps) => (
    <Tooltip {...props} classes={{ popper: className }} />
))({
    [`& .${tooltipClasses.tooltip}`]: {
        fontSize: 16,
        whiteSpace: 'pre-wrap',
    },
});

const EmptyTokenContainer = () => {
    return (
        <Box sx={{
            textAlign: 'left',
            backgroundColor: '#001435',
            fontSize: 16,
            fontWeight: 'bold',
            fontFamily: 'monospace',
            paddingLeft: '1em',
            height: 400,
        }}></Box>
    )
};

const TokenContainer = ({ children }) => {
    return (
        <Box sx={{
            textAlign: 'left',
            backgroundColor: '#001435',
            fontSize: 16,
            fontWeight: 'bold',
            fontFamily: 'monospace',
            paddingLeft: '1em',
            height: 400,
        }}>{children}</Box>
    )
};

const Lexer = (v: any) => {
    if (!v.tokens) {
        return <EmptyTokenContainer />
    }
    const tokens = v.tokens as Token[];
    return (
        <>
            <TokenContainer>
                {
                    tokens.map((tk, idx) => {
                        if (tokens.length > idx + 1) {
                            return tokenToComponent(tk, tokens[idx + 1]);
                        }
                        return tokenToComponent(tk, null)
                    })
                }
            </TokenContainer>
        </>
    )
};

const groupedTokenToComponent = (g: GroupedToken, groups: string[]) => {
    if (g.token) {
        return tokenToComponentWithoutTip(g.token);
    }
    return groupTokenToComponent(g.group, groups);
};

const groupTokenToComponent = (g: TokenGroup, groups: string[]) => {
    if (!g) {
        return <></>
    }
    const groupColor = () => {
        switch (g.type) {
            case "directive":
                return "dimgray";
            case "directive_name":
                return "orange"
            case "document":
                return "hotpink";
            case "document_body":
                return "olive"
            case "anchor":
                return "deepskyblue";
            case "anchor_name":
                return "brown"
            case "alias":
                return "limegreen";
            case "literal":
                return "gold"
            case "folded":
                return "gold"
            case "scalar_tag":
                return "blueviolet";
            case "map_key":
                return "coral";
            case "map_key_value":
                return "teal"
        }
        return 'white';
    }
    const newGroups = [...groups, g.type];
    return (
        <CustomTooltip
            sx={{ zIndex: 1500 + newGroups.length }}
            title={newGroups.join(' > ')}
            followCursor
            arrow
        >
            <Box sx={{
                backgroundColor: groupColor(),
                borderRadius: 1,
                paddingLeft: 1,
                paddingRight: 1,
                marginRight: '1em'
            }} component="span">
                {
                    g.tokens.map(tk => {
                        return groupedTokenToComponent(tk, newGroups);
                    })
                }
            </Box>
        </CustomTooltip>
    )
};

const tokenToComponent = (tk: Token, nextTk: Token | null) => {
    if (!tk) {
        return <></>
    }
    const orgs = tk.origin.split('\n').map((v, idx) => {
        if (idx > 0) {
            return ["\n", v];
        }
        return v;
    });
    const color = () => {
        if (nextTk && nextTk.type === 'MappingValue') {
            return '#008b8b';
        }
        switch (tk.type) {
            case "String":
                return '#ff7f50';
        }
        return 'black';
    };
    const prop = `type:   ${tk.type}
origin: ${JSON.stringify(tk.origin)}
value:  ${JSON.stringify(tk.value)}
line:   ${tk.line}
column: ${tk.column}
`
    return (
        <CustomTooltip title={prop}>
            <Box component="span" sx={{
                backgroundColor: 'white',
                paddingRight: 1,
                marginRight: 1,
                borderRadius: 1,
                color: color(),
            }}>
                {
                    orgs.map((v) => {
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
};

const tokenToComponentWithoutTip = (tk: Token) => {
    if (!tk) {
        return <></>
    }
    const orgs = tk.origin.split('\n').map((v, idx) => {
        if (idx > 0) {
            return ["\n", v];
        }
        return v;
    });
    return (
        <Box component="span" sx={{
            backgroundColor: 'white',
            paddingRight: 1,
            borderRadius: 1,
        }}>
            {
                orgs.map((v) => {
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
    )
};


const ParserGroup = (v: any) => {
    if (!v.tokens) {
        return <EmptyTokenContainer />
    }
    const tokens = v.tokens as GroupedToken[];
    return (
        <TokenContainer>
            {
                tokens.map(tk => { return groupedTokenToComponent(tk, []); })
            }
        </TokenContainer>
    )
};

const AST = (v: any) => {
    if (!v?.svg) {
        return <Box sx={{height: 400, backgroundColor: '#001435'}}></Box>
    }
    const parser = new DOMParser()
    const dom = parser.parseFromString(v.svg, 'text/xml');
    const g = dom.getElementById('graph0');
    if (!g) {
        return <Box sx={{height: 400, backgroundColor: '#001435'}}></Box>
    }
    const viewBox = g.parentElement!.getAttribute('viewBox')!;
    return (
        <Box sx={{height: 400}}>
            <svg width={'100%'} height={'100%'} viewBox={viewBox} dangerouslySetInnerHTML={{ __html: g.outerHTML }}></svg>
        </Box>
    )
}

function CodeEditor() {
    const editorRef = useRef<editor.IStandaloneCodeEditor | null>(null);
    const workerRef = useRef<Worker | null>(null);
    const [tokens, setTokens] = useState<Token[]>([]);
    const [groupedTokens, setGroupedTokens] = useState<GroupedToken[]>([]);
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
                            setTokens(v.result as Token[]);
                        }
                        break;
                    case YAMLProcessResultType.ParserGroup:
                        if (typeof v.result === 'string') {
                            console.error(v.result);
                        } else {
                            setGroupedTokens(v.result as GroupedToken[]);
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
                <Grid marginTop={12} size={{ xs: 12, md: 6 }}>
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
                <Grid size={{ xs: 12, md: 6 }}>
                    <Tabs
                        textColor='secondary'
                        indicatorColor='secondary'
                        value={value}
                        onChange={handleChange}
                        variant="scrollable"
                        scrollButtons="auto"
                        aria-label="tabs">
                        <Tab icon={<ArrowBackIosNew />} iconPosition="end" style={{ marginLeft: 20 }} label="OUTPUT" {...a11yProps(0)} />
                        <Tab icon={<ArrowBackIosNew />} iconPosition="end" style={{ marginLeft: 0 }} label="AST" {...a11yProps(1)} />
                        <Tab icon={<ArrowBackIosNew />} iconPosition="end" style={{ marginLeft: 0 }} label="GROUPED TOKENS" {...a11yProps(2)} />
                        <Tab style={{ marginLeft: 0 }} label="TOKENS" {...a11yProps(3)} />
                    </Tabs>
                    <TabPanel value={value} index={0}>
                        <TerminalComponent out={out} />
                    </TabPanel>
                    <TabPanel value={value} index={1}>
                        <AST svg={svg}></AST>
                    </TabPanel>
                    <TabPanel value={value} index={2}>
                        <ParserGroup tokens={groupedTokens}></ParserGroup>
                    </TabPanel>
                    <TabPanel value={value} index={3}>
                        <Lexer tokens={tokens}></Lexer>
                    </TabPanel>
                </Grid>
            </Grid>
        </>
    )
}

export default CodeEditor;
