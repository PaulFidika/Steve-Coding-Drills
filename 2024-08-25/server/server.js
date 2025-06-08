const express = require('express');
const path = require('path');
const cors = require('cors');
const app = express();
const port = 3001;

// Enable CORS for all routes
app.use(cors());

// Serve static files from the current directory
app.use(express.static(__dirname));

// Define routes for Lazy1.js, Lazy2.js, and Lazy3.js
app.get('/Lazy1.js', (req, res) => {
  res.type('application/javascript');
  res.sendFile(path.join(__dirname, 'Lazy1.js'));
});

app.get('/Lazy2.js', (req, res) => {
  res.type('application/javascript');
  res.sendFile(path.join(__dirname, 'Lazy2.js'));
});

app.get('/Lazy3.js', (req, res) => {
  res.type('application/javascript');
  res.sendFile(path.join(__dirname, 'Lazy3.js'));
});

// Start the server
app.listen(port, () => {
  console.log(`Server running at http://localhost:${port}`);
});
