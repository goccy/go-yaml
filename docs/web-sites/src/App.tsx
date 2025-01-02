import { Typography } from '@mui/material'
import './App.css'
import CodeEditor from './CodeEditor.tsx'

function App() {
  return (
    <>
      <Typography variant="h4">
        goccy/go-yaml Playground
      </Typography>
      <CodeEditor></CodeEditor>
    </>
  )
}

export default App
