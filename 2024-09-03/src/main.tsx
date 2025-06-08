import React, { useState } from 'react'
import ReactDOM from 'react-dom/client'
import { NextUIProvider, Checkbox, Button, CircularProgress, Dropdown, DropdownTrigger, DropdownMenu, DropdownItem, Select, SelectItem } from '@nextui-org/react'
import './index.css'

const App = () => {
  const [selectedFruit, setSelectedFruit] = useState(new Set<string>(["apple"]))
  const [isDrawerOpen, setIsDrawerOpen] = useState(false)

  const fruits = [
    { key: "apple", label: "Apple" },
    { key: "banana", label: "Banana" },
    { key: "orange", label: "Orange" }
  ]

  const toggleDrawer = () => {
    setIsDrawerOpen(!isDrawerOpen)
  }

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold mb-4">Hello, World!</h1>
      <p className="mb-4">This is a React app with NextUI components.</p>

      <div className="space-y-4">
        <Checkbox defaultSelected>Option 1</Checkbox>
        <Checkbox>Option 2</Checkbox>

        <Button color="primary">Click me</Button>

        <CircularProgress value={70} color="secondary" />

        <Dropdown>
          <DropdownTrigger>
            <Button 
              variant="bordered" 
            >
              Open Menu
            </Button>
          </DropdownTrigger>
          <DropdownMenu aria-label="Static Actions">
            <DropdownItem key="new">New file</DropdownItem>
            <DropdownItem key="copy">Copy link</DropdownItem>
            <DropdownItem key="edit">Edit file</DropdownItem>
            <DropdownItem key="delete" className="text-danger" color="danger">
              Delete file
            </DropdownItem>
          </DropdownMenu>
        </Dropdown>

        <div className="flex w-full flex-wrap md:flex-nowrap gap-4">
          <Select 
            label="Select a fruit" 
            className="max-w-xs"
            selectedKeys={selectedFruit}
            onSelectionChange={(key) => setSelectedFruit(key as Set<string>)}
          >
            {fruits.map((fruit) => (
              <SelectItem key={fruit.key} value={fruit.key}>
                {fruit.label}
              </SelectItem>
            ))}
          </Select>
          <Select
            label="Favorite Fruit"
            placeholder="Select a fruit"
            className="max-w-xs"
          >
            {fruits.map((fruit) => (
              <SelectItem key={fruit.key} value={fruit.key}>
                {fruit.label}
              </SelectItem>
            ))}
          </Select>
        </div>

        <Button onClick={toggleDrawer}>Toggle Drawer</Button>

        {isDrawerOpen && (
          <div className="fixed inset-y-0 right-0 w-64 bg-white shadow-lg p-4 transition-transform transform translate-x-0">
            <h2 className="text-xl font-bold mb-4">Drawer Content</h2>
            <p>This is the drawer content.</p>
            <Button onClick={toggleDrawer} className="mt-4">Close Drawer</Button>
          </div>
        )}
      </div>
    </div>
  )
}

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
      <NextUIProvider>
        <main className="dark text-foreground bg-background">
          <div className="bg-purple-900 min-h-screen">
            <App />
          </div>
        </main>
      </NextUIProvider>
  </React.StrictMode>
)