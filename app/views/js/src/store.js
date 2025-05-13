import { legacy_createStore as createStore } from 'redux'

const initialState = {
  sidebarShow: true,
  theme: 'light',
  auth: localStorage.getItem('auth') ? JSON.parse(localStorage.getItem('auth')) : null,
}

const changeState = (state = initialState, action) => {
  switch (action.type) {
    case 'set':
      if ('auth' in action) {
        localStorage.setItem('auth', JSON.stringify(action.auth))
      }

      return { ...state, ...action }

    case 'logout':
      localStorage.removeItem('auth')
      return { ...state, auth: null }

    default:
      return state
  }
}

const store = createStore(changeState)
export default store
