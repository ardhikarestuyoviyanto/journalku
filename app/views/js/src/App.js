import React, { Suspense, useEffect } from 'react'
import { HashRouter, Route, Routes } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { CSpinner, useColorModes } from '@coreui/react'
import './scss/style.scss'
import './scss/examples.scss'
import { ToastContainer } from 'react-toastify'
import PrivateRouteAuth from './PrivateRouteAuth'
import PrivateRouteCurrentCompany from './PrivateRouteCurrentCompany'

// Pages
const SignIn = React.lazy(() => import('./views/v1/auth/SignIn'))
const SignUp = React.lazy(() => import('./views/v1/auth/SignUp'))
const Page404 = React.lazy(() => import('./views/pages/page404/Page404'))
const Page500 = React.lazy(() => import('./views/pages/page500/Page500'))
const ChooseCompany = React.lazy(() => import('./views/v1/company/ChooseCompany'))
const CreateCompany = React.lazy(() => import('./views/v1/company/CreateCompany'))
const Dashboard = React.lazy(() => import('./views/v1/dashboard/Dashboard'))

const App = () => {
  const { isColorModeSet, setColorMode } = useColorModes('coreui-free-react-admin-template-theme')
  const storedTheme = useSelector((state) => state.theme)

  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.href.split('?')[1])
    const theme = urlParams.get('theme') && urlParams.get('theme').match(/^[A-Za-z0-9\s]+/)[0]
    if (theme) {
      setColorMode(theme)
    }

    if (isColorModeSet()) {
      return
    }

    setColorMode(storedTheme)
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <Suspense
      fallback={
        <div className="pt-3 text-center">
          <CSpinner color="primary" variant="grow" />
        </div>
      }
    >
      <Routes>
        <Route exact path="/" name="SignIn Page" element={<SignIn />} />
        <Route exact path="/signup" name="SignUp Page" element={<SignUp />} />
        <Route exact path="/404" name="Page 404" element={<Page404 />} />
        <Route exact path="/500" name="Page 500" element={<Page500 />} />

        {/* ProtectedRoute */}
        <Route
          path="/choose-company"
          element={
            <PrivateRouteAuth>
              <ChooseCompany />
            </PrivateRouteAuth>
          }
        />
        <Route
          path="/choose-company/create"
          element={
            <PrivateRouteAuth>
              <CreateCompany />
            </PrivateRouteAuth>
          }
        />
        <Route
          path="/dashboard"
          element={
            <PrivateRouteCurrentCompany>
              <Dashboard />
            </PrivateRouteCurrentCompany>
          }
        />
      </Routes>
      <ToastContainer position="top-right" autoClose={3000} />
    </Suspense>
  )
}

export default App
