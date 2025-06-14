import { useSelector } from 'react-redux'
import { Navigate } from 'react-router-dom'

const PrivateRouteAuth = ({ children }) => {
  const auth = useSelector((state) => state.auth)
  if (!auth) {
    return <Navigate to="/" replace />
  }

  return children
}

export default PrivateRouteAuth
