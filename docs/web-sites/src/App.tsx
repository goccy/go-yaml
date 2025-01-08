/// <reference types="vite-plugin-svgr/client" />
import { AppBar, Toolbar, Typography, Box, Card, CardContent, CardActionArea, Stack, Icon, SvgIcon } from '@mui/material';
import GitHub from '@mui/icons-material/GitHub';
import './App.css'
import CodeEditor from './CodeEditor.tsx'
import YAMLIcon from '/public/yaml.svg?react';

const Header = () => {
  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar position="fixed">
        <Toolbar sx={{ backgroundColor: '#313131' }}>
          <SvgIcon sx={{ md: 2, transform: 'scale(1.8)', marginRight: 3 }}><YAMLIcon /></SvgIcon>
          <Typography variant="h5" sx={{ flexGrow: 1, textAlign: 'left' }}>
            goccy/go-yaml Playground
          </Typography>
          <GitHubLink></GitHubLink>
        </Toolbar>
      </AppBar>
    </Box>
  )
};

const GitHubLink = () => {
  return (
    <Box sx={{}}>
      <Card sx={{
        maxWidth: 200,
      }} variant="outlined">
        <CardActionArea href="https://github.com/goccy/go-yaml">
          <CardContent>
            <Stack direction="row" spacing={1} sx={{ alignItems: 'center' }}>
              <GitHub />
              <Typography variant="caption">
                Visit Our GitHub
              </Typography>
            </Stack>
          </CardContent>
        </CardActionArea>
      </Card>
    </Box>
  )
};

function App() {
  return (
    <>
      <Header />
      <CodeEditor></CodeEditor>
    </>
  )
}

export default App
