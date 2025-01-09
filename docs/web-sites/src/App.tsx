/// <reference types="vite-plugin-svgr/client" />
import { AppBar, Toolbar, Typography, Box, Card, CardContent, CardActionArea, Stack, SvgIcon, Button, TextField, IconButton } from '@mui/material';
import GitHub from '@mui/icons-material/GitHub';
import Share from '@mui/icons-material/Share';
import ContentCopy from '@mui/icons-material/ContentCopy';
import './App.css'
import CodeEditor from './CodeEditor.tsx'
import YAMLIcon from '/public/yaml.svg?react';
import { useState } from 'react';

const themeBlack = '#313131';
const themeWhite = 'white';

const Header = (content: any) => {
  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar position="fixed">
        <Toolbar sx={{ backgroundColor: themeBlack }}>
          <SvgIcon sx={{ md: 2, transform: 'scale(1.8)', marginRight: 3 }}><YAMLIcon /></SvgIcon>
          <Typography variant="h5" sx={{ flexGrow: 1, textAlign: 'left' }}>
            goccy/go-yaml Playground
          </Typography>
          <Stack sx={{
            backgroundColor: '#A0A0A0',
            paddingLeft: 0.4,
            paddingTop: 0.2,
          }} direction={'row'} alignItems={'center'}>
            <TextField value={'foo'} sx={{
              backgroundColor: themeWhite,
            }} variant="standard"></TextField>
            <Button size='small' sx={{
              backgroundColor: '#A0A0A0',
              color: themeBlack,
            }}>
              <ContentCopy/>
            </Button>
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
          }} startIcon={<Share sx={{ color: themeBlack }} />}>
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

function App() {
  const [content, setContent] = useState<string>('');

  return (
    <>
      <Header content={content}/>
      <CodeEditor></CodeEditor>
    </>
  )
}

export default App
