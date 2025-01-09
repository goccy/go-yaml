/// <reference types="vite-plugin-svgr/client" />
import { AppBar, Toolbar, Typography, Box, Stack, SvgIcon, Button, TextField, Tabs, Tab, Tooltip, TooltipProps, tooltipClasses, styled, Snackbar, SnackbarCloseReason } from '@mui/material';
import GitHub from '@mui/icons-material/GitHub';
import Share from '@mui/icons-material/Share';
import ContentCopy from '@mui/icons-material/ContentCopy';
import './App.css'
import YAMLIcon from '/public/yaml.svg?react';
import { useState, useRef, useEffect } from 'react';
import { editor } from 'monaco-editor'
import MonacoEditor, { loader } from '@monaco-editor/react';
import Grid from '@mui/material/Grid2';
import { FitAddon } from '@xterm/addon-fit';
import { Token, TokenGroup, GroupedToken, YAMLProcessResult, YAMLProcessResultType } from './YAML.ts';
import yamlWorker from "./worker?worker";
import '@xterm/xterm/css/xterm.css';
import { ArrowBackIosNew } from '@mui/icons-material'
import { Terminal } from '@xterm/xterm';

const themeBlack = '#313131';
const themeWhite = 'white';

const Header = (content: any) => {
  const v = btoa(content.content as string);
  const origin = window.location.origin;
  const shareURL = `${origin}/?content=${v}`;
  const [open, setOpen] = useState(false);
  const [shareURLFieldVisibility, setShareURLFieldVisibility] = useState('hidden');
  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar position="fixed">
        <Toolbar sx={{ backgroundColor: themeBlack }}>
          <SvgIcon sx={{ md: 2, transform: 'scale(1.8)', marginRight: 3 }}><YAMLIcon /></SvgIcon>
          <Typography variant="h5" sx={{ flexGrow: 1, textAlign: 'left' }}>
            goccy/go-yaml Playground
          </Typography>
          <Stack sx={{
            visibility: shareURLFieldVisibility,
            backgroundColor: '#A0A0A0',
            paddingLeft: 0.4,
            paddingTop: 0.2,
          }} direction={'row'} alignItems={'center'}>
            <TextField value={shareURL} sx={{
              backgroundColor: themeWhite,
            }} variant="standard"></TextField>
            <Button size='small' sx={{
              backgroundColor: '#A0A0A0',
              color: themeBlack,
            }}
              onClick={() => {
                navigator.clipboard.writeText(shareURL).then(() => {
                  setOpen(true);
                });
              }}>
              <ContentCopy />
            </Button>
            <Snackbar
              open={open}
              autoHideDuration={1500}
              onClose={
                (event: React.SyntheticEvent | Event,
                  reason?: SnackbarCloseReason) => {
                  if (reason === 'clickaway') {
                    return;
                  }
                  setOpen(false);
                }}
              message="Copied"
            />
          </Stack>
          <Button variant="contained" sx={{
            backgroundColor: themeWhite,
            color: themeBlack,
            marginLeft: 4,
            marginRight: 4,
            textTransform: 'none',
            fontWeight: 'bold',
            paddingRight: 2,
            paddingLeft: 2,
            borderRadius: 0,
          }}
           startIcon={<Share sx={{ color: themeBlack }} />}
           onClick={() => {
            setShareURLFieldVisibility('visible');
           }}
           >
            Share
          </Button>
          <Button variant="contained" sx={{
            color: themeBlack,
            backgroundColor: 'white',
            fontWeight: 'bold',
            paddingRight: 2,
            paddingLeft: 2,
            textTransform: 'none',
          }}
            startIcon={<GitHub sx={{ color: themeBlack }}></GitHub>}
            href="https://github.com/goccy/go-yaml"
          >
            Visit Our GitHub
          </Button>
        </Toolbar>
      </AppBar>
    </Box>
  )
};

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
  const terminalRef = useRef<HTMLDivElement>(null)
  const [terminalInstance, setTerminalInstance] = useState<Terminal | null>(null)
  const fitAddon = new FitAddon();
  const out = v.out;

  useEffect(() => {
    const instance = new Terminal({
      cursorInactiveStyle: 'none',
      cursorStyle: 'bar',
      letterSpacing: 4,
      fontFamily: 'monospace',
      fontSize: 16,
      fontWeightBold: 'bold',
      convertEol: true,
      theme: {
        background: '#313131',
        cursor: '#313131',
        brightRed: '#da433a',
      },
    });

    instance.loadAddon(fitAddon);

    if (terminalRef.current) {
      instance.open(terminalRef.current)
    }

    setTerminalInstance(instance)

    return () => {
      instance.dispose()
      setTerminalInstance(null)
    }
  }, [
    terminalRef,
  ]);

  useEffect(() => {
    if (!terminalInstance) {
      return;
    }

    terminalInstance.loadAddon(fitAddon);
    fitAddon.fit();
    terminalInstance.clear();
    terminalInstance.writeln(out);
    const handleResize = () => fitAddon.fit()

    window.addEventListener('resize', handleResize)
    return () => {
      window.removeEventListener('resize', handleResize);
    }
  }, [terminalRef, terminalInstance, out]);

  return (
    <Box component="div" sx={{
      height: 400,
      width: '100%',
      backgroundColor: '#313131',
    }}>
      <Box
        component="div"
        ref={terminalRef}
        sx={{
          height: '100%',
          width: '90%',
          textAlign: 'left',
          paddingLeft: 1,
          paddingTop: 1,
        }} />
    </Box>
  )
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
      backgroundColor: '#313131',
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
      backgroundColor: '#313131',
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
    return <Box sx={{ height: 400, backgroundColor: '#313131' }}></Box>
  }
  const parser = new DOMParser()
  const dom = parser.parseFromString(v.svg, 'text/xml');
  const g = dom.getElementById('graph0');
  if (!g) {
    return <Box sx={{ height: 400, backgroundColor: '#313131' }}></Box>
  }
  const viewBox = g.parentElement!.getAttribute('viewBox')!;
  return (
    <Box sx={{ height: 400 }}>
      <svg width={'100%'} height={'100%'} viewBox={viewBox} dangerouslySetInnerHTML={{ __html: g.outerHTML }}></svg>
    </Box>
  )
}

function App() {
  const [content, setContent] = useState<string>('');
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
    setContent(code);
    if (workerRef.current) {
      workerRef.current.postMessage(code);
    }
  };
  loader.init().then((monaco) => {
    monaco.editor.defineTheme('go-yaml-theme', {
      base: 'vs-dark',
      inherit: true,
      rules: [],
      colors: {
        'editor.background': '#313131',
        'editor.selectionHighlightBorder': '#313131',
        'editor.lineHighlightBackground': '#313131',
        'editor.selectionBackground': '#313131',
      },
    });
  });

  const search = window.location.search;
  const yamlDataBinary = new URLSearchParams(search).get('content');

  const onMount = (editor: editor.IStandaloneCodeEditor, monaco) => {
    editorRef.current = editor;
    if (yamlDataBinary) {
      const code = atob(yamlDataBinary)
      editor.setValue(code);
      setContent(code);
      if (workerRef.current) {
        workerRef.current.postMessage(code);
      }
    }
  };
  const [value, setValue] = useState(0);
  const handleChange = (event: React.SyntheticEvent, newValue: number) => {
    setValue(newValue);
  };
  return (
    <>
      <Header content={content} />
      <Grid container>
        <Grid marginTop={12} size={{ xs: 12, md: 6 }}>
          <MonacoEditor
            height={400}
            language="yaml"
            theme="go-yaml-theme"
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

export default App
