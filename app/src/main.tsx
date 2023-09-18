import { render } from 'preact'
import { App } from './app.tsx'
import '@fontsource-variable/inter'
import './index.css'

render(<App />, document.getElementById('app') as HTMLElement)
