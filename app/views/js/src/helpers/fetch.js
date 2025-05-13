import store from '../store'
import i18next from 'i18next'
import history from '../history'

const fetchWithAuth = async (url, options = {}) => {
  const state = store.getState()
  const token = state.auth?.token

  const headers = {
    ...options.headers,
    Authorization: `Bearer ${token}`,
    'Accept-Lang': i18next.language,
  }

  try {
    const res = await fetch(url, {
      ...options,
      headers,
    })

    if (res.status === 401) {
      store.dispatch({ type: 'logout' })
      history.push('/', { unauthorized: true })
      return null
    }

    const data = await res.json()
    return data
  } catch (err) {
    console.error('Fetch error:', err)
    throw err
  }
}

export default fetchWithAuth
