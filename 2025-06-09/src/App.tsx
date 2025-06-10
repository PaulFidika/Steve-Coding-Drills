import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import { Button } from '@/components/ui/button'

function App() {
  const [count, setCount] = useState(0)

  return (
    <div className="font-sans antialiased relative min-h-screen">
      <div className="texture" />
      <div className="relative z-10">
        <div>
          <a href="https://vite.dev" target="_blank">
            <img src={viteLogo} className="logo" alt="Vite logo" />
          </a>
          <a href="https://react.dev" target="_blank">
            <img src={reactLogo} className="logo react" alt="React logo" />
          </a>
        </div>
        <h1 className="font-serif text-4xl font-bold mb-8">Vite + React</h1>
        <div className="card">
          <Button onClick={() => setCount((count) => count + 1)} className="mb-4">
            count is {count}
          </Button>
          <p className="font-sans">
            Edit <code>src/App.tsx</code> and save to test HMR
          </p>
        </div>
        <p className="read-the-docs font-sans text-sm opacity-75">
          Click on the Vite and React logos to learn more
        </p>
      </div>
    </div>
  )
}

export default App
