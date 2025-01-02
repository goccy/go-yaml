import { useState, useRef } from 'react';
import { editor } from 'monaco-editor'
import MonacoEditor from '@monaco-editor/react';
import AST from './AST.tsx'
import { Box, Tabs, Tab } from '@mui/material';
import Grid from '@mui/material/Grid2';

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

function CodeEditor() {
    const editorRef = useRef<editor.IStandaloneCodeEditor | null>(null);
    const onChange = () => {
        console.log('code = ', editorRef?.current?.getValue());
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
                <Grid size={{ xs: 6, md: 6 }}>
                    <MonacoEditor
                        height={400}
                        language="yaml"
                        theme="vs-dark"
                        value={'foo: bar'}
                        options={{
                            selectOnLineNumbers: true,
                            renderWhitespace: 'all',
                        }}
                        onChange={onChange}
                        onMount={onMount}
                    />
                </Grid>
                <Grid size={{ xs: 6, md: 6 }}>
                    <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                        <Tabs value={value} onChange={handleChange} aria-label="basic tabs example">
                            <Tab label="Console" {...a11yProps(0)} />
                            <Tab label="Lexer" {...a11yProps(1)} />
                            <Tab label="Parser" {...a11yProps(2)} />
                            <Tab label="Three" {...a11yProps(3)} />
                        </Tabs>
                        <TabPanel value={value} index={0}>
                            OUT
                        </TabPanel>
                        <TabPanel value={value} index={1}>
                            Lexer
                        </TabPanel>
                        <TabPanel value={value} index={2}>
                            <AST></AST>
                        </TabPanel>
                        <TabPanel value={value} index={3}>
                            Item Three
                        </TabPanel>
                    </Box>
                </Grid>
            </Grid>
        </>
    )
}

export default CodeEditor;
