import { useSelector } from 'react-redux'
import { Navigate } from 'react-router-dom'

const PrivateRouteCurrentCompany = ({ children }) => {
  const auth = useSelector((state) => state.auth)
  if (!auth) {
    return <Navigate to="/" replace />
  }

  // Cek If CurrentCompanyId Exits
  if (auth.user.currentCompany == null) {
    return <Navigate to="/choose-company" replace />
  }

  return children
}

export default PrivateRouteCurrentCompany
