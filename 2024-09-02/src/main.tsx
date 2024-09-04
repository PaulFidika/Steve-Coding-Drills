import React, { useState } from 'react';
import ReactDOM from 'react-dom/client';
import PopoverDemo from './components/popover';
import Popover2Demo from './components/popover2';
import Popover3Demo from './components/popover3';
import '@radix-ui/themes/styles.css';
import { Theme, ThemePanel, Button, Dialog, DropdownMenu, Select, Checkbox } from '@radix-ui/themes';
import { Button as MuiButton, ThemeProvider, createTheme } from '@mui/material';

// Create a custom Material-UI theme
const muiTheme = createTheme({
  palette: {
    primary: {
      main: '#2196f3', // Custom blue color
    },
    secondary: {
      main: '#f50057', // Custom pink color
    },
  },
  typography: {
    fontFamily: 'Roboto, Arial, sans-serif',
    h1: {
      fontSize: '2.5rem',
      fontWeight: 500,
    },
    button: {
      textTransform: 'none',
    },
  },
  shape: {
    borderRadius: 8,
  },
});

function App() {
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [message, setMessage] = useState('');
  const [dropdownValue, setDropdownValue] = useState('');
  const [selectedFruit, setSelectedFruit] = useState('apple');
  const [checkboxValues, setCheckboxValues] = useState({
    option1: false,
    option2: false,
    option3: false,
  });

  const handleButtonClick = (buttonColor: string) => {
    setMessage(`You clicked the ${buttonColor} button!`);
    setIsDialogOpen(true);
  };

  const handleCheckboxChange = (option: string) => {
    setCheckboxValues(prev => ({ ...prev, [option]: !prev[option] }));
  };

  return (
  <Theme
    appearance="dark"
    accentColor="green"
    grayColor="mauve"
    panelBackground="solid"
    scaling="100%"
    radius="medium"
  >
    <ThemeProvider theme={muiTheme}>
      <div>
        <h1>Hello, World!</h1>
        <PopoverDemo />
        <Popover2Demo />
        <Popover3Demo />
        <br></br>
        <Button variant="solid" color="blue" onClick={() => handleButtonClick('Blue')}>Blue Button</Button>
        <Button variant="soft" color="red" onClick={() => handleButtonClick('Red')}>Red Button</Button>
        <Button variant="outline" color="green" onClick={() => handleButtonClick('Green')}>Green Button</Button>
        <br></br>
        <MuiButton variant="contained" color="primary" onClick={() => handleButtonClick('MUI Blue')}>MUI Blue Button</MuiButton>
        <MuiButton variant="contained" color="secondary" onClick={() => handleButtonClick('MUI Pink')}>MUI Pink Button</MuiButton>
        <MuiButton variant="outlined" color="success" onClick={() => handleButtonClick('MUI Green')}>MUI Green Button</MuiButton>
        {/* <ThemePanel /> */}
        <Dialog.Root open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <Dialog.Content>
            <Dialog.Title>Button Clicked</Dialog.Title>
            <Dialog.Description>{message}</Dialog.Description>
            <Dialog.Close>
              <Button>Close</Button>
            </Dialog.Close>
          </Dialog.Content>
        </Dialog.Root>
        
        <h2>Radix Dropdowns and Select Components</h2>
        
        <DropdownMenu.Root>
          <DropdownMenu.Trigger>
            <Button variant="soft">Open Dropdown</Button>
          </DropdownMenu.Trigger>
          <DropdownMenu.Content sideOffset={5} className="animate-in slide-in-from-top-2 duration-500">
            <DropdownMenu.Item onSelect={() => setDropdownValue('Edit')}>Edit</DropdownMenu.Item>
            <DropdownMenu.Item onSelect={() => setDropdownValue('Duplicate')}>Duplicate</DropdownMenu.Item>
            <DropdownMenu.Separator />
            <DropdownMenu.Item onSelect={() => setDropdownValue('Archive')}>Archive</DropdownMenu.Item>
            <DropdownMenu.Item onSelect={() => setDropdownValue('Delete')} color="red">Delete</DropdownMenu.Item>
          </DropdownMenu.Content>
        </DropdownMenu.Root>
        <p>Selected dropdown value: {dropdownValue}</p>

        <Select.Root 
          defaultValue="apple" 
          onValueChange={setSelectedFruit}
        >
          <Select.Trigger />
          <Select.Content className="animate-in slide-in-from-top-2 duration-1000">
            <Select.Group>
              <Select.Label>Fruits</Select.Label>
              <Select.Item value="apple">Apple</Select.Item>
              <Select.Item value="banana">Banana</Select.Item>
              <Select.Item value="blueberry">Blueberry</Select.Item>
              <Select.Item value="grapes">Grapes</Select.Item>
              <Select.Item value="pineapple">Pineapple</Select.Item>
            </Select.Group>
          </Select.Content>
        </Select.Root>
        <p>Selected fruit: {selectedFruit}</p>

        <h2>Radix Checkboxes</h2>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
          <Checkbox checked={checkboxValues.option1} onCheckedChange={() => handleCheckboxChange('option1')}>
            Option 1
          </Checkbox>
          <Checkbox checked={checkboxValues.option2} onCheckedChange={() => handleCheckboxChange('option2')}>
            Option 2
          </Checkbox>
          <Checkbox checked={checkboxValues.option3} onCheckedChange={() => handleCheckboxChange('option3')}>
            Option 3
          </Checkbox>
        </div>
        <p>Selected options: {Object.entries(checkboxValues).filter(([_, value]) => value).map(([key]) => key).join(', ')}</p>
      </div>
    </ThemeProvider>
    </Theme>
  );
}

const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
