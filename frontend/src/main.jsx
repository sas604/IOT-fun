import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App.jsx';
import 'normalize.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/700.css';
import { GlobaStyle } from './styles/GlobalStyles.js';

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <GlobaStyle />
    <App />
  </React.StrictMode>
);
