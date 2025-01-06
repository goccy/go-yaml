import { forwardRef, useState, useRef, useEffect } from 'react';
import { editor } from 'monaco-editor'
import MonacoEditor from '@monaco-editor/react';
import AST from './AST.tsx'
import { Box, Tabs, Tab } from '@mui/material';
import Grid from '@mui/material/Grid2';
import { useXTerm } from 'react-xtermjs';
import { FitAddon } from '@xterm/addon-fit';
import { YAMLProcessResult, YAMLProcessResultType } from './YAML.ts';
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

const TerminalComponent = () => {
    const { instance, ref } = useXTerm()
    const fitAddon = new FitAddon()

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

        fitAddon.fit();

        const handleResize = () => fitAddon.fit()

        // Write custom message on your terminal
        instance?.writeln("\x1b[31mRED\x1b[0m\n");
        //instance?.onData((data) => instance?.write(data))
        window.addEventListener('resize', handleResize)
        return () => {
            window.removeEventListener('resize', handleResize)
        }
    }, [ref, instance])

    return <div ref={ref} style={{ height: 400, width: '100%', textAlign: 'left' }} />
}

interface Token {
    value: string
    origin: string
    error: string
    line: number
    column: number
    offset: number
}

const Lexer = (v: any) => {
    console.log(v.tokens);
    return (
        <>
            {
                (v.tokens as Token[]).map((token) => {
                    return <span style={{
                        backgroundColor: 'gray',
                        margin: 4,
                        padding: 5,
                    }}>{token.origin}</span>;
                })
            }
        </>
    )
};

function CodeEditor() {
    const editorRef = useRef<editor.IStandaloneCodeEditor | null>(null);
    const workerRef = useRef<Worker | null>(null);
    const [tokens, setTokens] = useState<Token[]>([]);
    useEffect(() => {
        workerRef.current = new yamlWorker();
        workerRef.current.onmessage = (event) => {
            const data = event.data as YAMLProcessResult[];
            data.forEach((v) => {
                console.log(v);
                switch (v.type) {
                    case YAMLProcessResultType.Lexer:
                        if (typeof v.result === 'string') {
                            console.error(v.result);
                        } else {
                            setTokens(v.result);
                        }
                        break;
                    case YAMLProcessResultType.Parser:
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
                            <TerminalComponent />
                        </TabPanel>
                        <TabPanel value={value} index={1}>
                            <Lexer tokens={tokens}></Lexer>
                        </TabPanel>
                        <TabPanel value={value} index={2}>
                            Parser(Grouping)
                        </TabPanel>
                        <TabPanel value={value} index={3}>
                            <AST></AST>
                        </TabPanel>
                    </Box>
                </Grid>
            </Grid>
        </>
    )
}

export default CodeEditor;
