import React from 'react';
import logo from './logo.svg';
import './App.css';
import ReactDom from 'react-dom';
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom';

ReactDOM.render(
  <Router>
    <Header />
    <Switch></Switch>
  </Router>
);

function Header() {
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>Home</p>
        <p>Inventory</p>
        <p>Mint</p>
        <p>Connect Wallet</p>
      </header>
    </div>
  );
}

export default App;
