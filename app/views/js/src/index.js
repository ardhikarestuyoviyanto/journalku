import React from 'react'
import { createRoot } from 'react-dom/client'
import { Provider } from 'react-redux'
import { unstable_HistoryRouter as HistoryRouter } from 'react-router-dom'

import App from './App'
import store from './store'
import history from './history'
import './helpers/i18n'

createRoot(document.getElementById('root')).render(
  <HistoryRouter history={history}>
    <Provider store={store}>
      <App />
    </Provider>
  </HistoryRouter>,
)
