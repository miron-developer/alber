import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router } from 'react-router-dom';

import { IsLogged, UserOnline } from 'utils/user';
import App from './App'

import reportWebVitals from './reportWebVitals';

const check = async () => {
  const id = await IsLogged();
  if (id) return await UserOnline(id);
}

const init = async () => {
  await check();
  ReactDOM.render(
    <React.StrictMode >
      <Router >
        <App />
      </Router>
    </React.StrictMode>,
    document.getElementById('root')
  );
}

init();

reportWebVitals();
