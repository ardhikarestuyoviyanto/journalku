import { useSelector } from 'react-redux'
import { Navigate } from 'react-router-dom'

const PrivateRoute = ({ children }) => {
  const auth = useSelector((state) => state.auth)
  if (!auth) {
    return <Navigate to="/" replace />
  }

  return children
}

export default PrivateRoute
