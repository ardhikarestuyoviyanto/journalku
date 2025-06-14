import React, { useEffect, useRef, useState } from 'react'
import { useSelector, useDispatch } from 'react-redux'
import {
  CContainer,
  CDropdown,
  CDropdownItem,
  CDropdownMenu,
  CDropdownToggle,
  CHeader,
  CHeaderNav,
  CHeaderToggler,
  CNavLink,
  CNavItem,
  useColorModes,
  CBreadcrumb,
  CBreadcrumbItem,
  CInputGroup,
  CButton,
} from '@coreui/react'
import CIcon from '@coreui/icons-react'
import { cilContrast, cilMenu, cilMoon, cilSun, cilGlobeAlt, cilPlaylistAdd } from '@coreui/icons'
import store from '../store'
import { AppHeaderDropdown } from './header/index'
import { useTranslation } from 'react-i18next'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faShoppingBasket, faShoppingBag, faCreditCard } from '@fortawesome/free-solid-svg-icons'
import fetchWithAuth from '../helpers/fetch'
import { toast } from 'react-toastify'
import { Link } from 'react-router-dom'

const navButtonStyle = {
  display: 'flex',
  alignItems: 'center',
  gap: '6px',
  padding: '6px 12px',
  border: '1px solid #ddd',
  borderRadius: '8px',
  color: '#4a4a4a',
  fontWeight: '500',
  backgroundColor: 'white',
  transition: 'all 0.2s ease-in-out',
  textDecoration: 'none',
}

const iconStyle = {
  color: '#6c63ff',
}

const AppHeader = ({ breadCrumbs }) => {
  const headerRef = useRef()
  const { colorMode, setColorMode } = useColorModes('coreui-free-react-admin-template-theme')
  const { t, i18n } = useTranslation()
  const dispatch = useDispatch()
  const sidebarShow = useSelector((state) => state.sidebarShow)
  const state = store.getState()
  const auth = state.auth
  const [loading, setLoading] = useState(false)
  const [dataInit, setDataInit] = useState(null)

  const handleGetInit = async () => {
    const data = await fetchWithAuth(`${import.meta.env.VITE_API_BASE_URL}/init`, {
      method: 'GET',
    })
    if (data && data.success) {
      setDataInit(data.data)
    } else {
      toast.error(data.error)
    }
  }

  const handleChooseCompany = async (companyId) => {
    const formData = new FormData()
    formData.append('companyId', companyId)
    setLoading(true)
    const response = await fetchWithAuth(`${import.meta.env.VITE_API_BASE_URL}/choose-company`, {
      method: 'POST',
      body: formData,
    })
    setLoading(false)
    if (response.success) {
      dispatch({
        type: 'set',
        auth: response.data,
      })
      toast.success(t('switchingCompanySuccess'))
    } else {
      toast.error(data.error)
    }
  }

  useEffect(() => {
    document.addEventListener('scroll', () => {
      headerRef.current &&
        headerRef.current.classList.toggle('shadow-sm', document.documentElement.scrollTop > 0)
    })
    // Get Init Data
    handleGetInit()
  }, [])

  return (
    <CHeader position="sticky" className="mb-4 p-0" ref={headerRef}>
      <CContainer className="border-bottom px-4" fluid>
        <CHeaderToggler
          onClick={() => dispatch({ type: 'set', sidebarShow: !sidebarShow })}
          style={{ marginInlineStart: '-14px' }}
        >
          <CIcon icon={cilMenu} size="lg" />
        </CHeaderToggler>
        <CHeaderNav className="d-none d-md-flex" style={{ gap: '8px' }}>
          <CNavItem>
            <CNavLink href="#" style={navButtonStyle}>
              <FontAwesomeIcon icon={faShoppingBasket} style={iconStyle} />
              {t('sell')}
            </CNavLink>
          </CNavItem>
          <CNavItem>
            <CNavLink href="#" style={navButtonStyle}>
              <FontAwesomeIcon icon={faShoppingBag} style={iconStyle} />
              {t('buy')}
            </CNavLink>
          </CNavItem>
          <CNavItem>
            <CNavLink href="#" style={navButtonStyle}>
              <FontAwesomeIcon icon={faCreditCard} style={iconStyle} />
              {t('cost')}
            </CNavLink>
          </CNavItem>
        </CHeaderNav>
        <CHeaderNav className="ms-auto align-items-center">
          {/* Nama Perusahaan */}
          <div className="px-3">
            {dataInit != null && (
              <CInputGroup>
                <select
                  disabled={loading}
                  onChange={(e) => {
                    handleChooseCompany(e.target.value)
                  }}
                  className="form-control"
                  defaultValue={auth?.user?.currentCompany?.id}
                >
                  {dataInit?.companyAccess?.map((comp, index) => (
                    <option key={index} value={comp.id}>
                      {comp.name}
                    </option>
                  ))}
                </select>
                <Link to={'/choose-company/create'}>
                  <CButton type="button" color="secondary" variant="outline" id="button-addon2">
                    <CIcon icon={cilPlaylistAdd} size="lg" />
                  </CButton>
                </Link>
              </CInputGroup>
            )}
          </div>

          {/* Garis Vertikal */}
          <li className="nav-item py-1">
            <div className="vr h-100 mx-2 text-body text-opacity-75"></div>
          </li>

          {/* Dropdown Bahasa */}
          <CDropdown variant="nav-item" className="px-2">
            <CDropdownToggle caret={false} className="d-flex align-items-center gap-2">
              <span>ID</span>
              <CIcon icon={cilGlobeAlt} size="lg" />
            </CDropdownToggle>
            <CDropdownMenu>
              <CDropdownItem>English</CDropdownItem>
              <CDropdownItem>Bahasa Indonesia</CDropdownItem>
            </CDropdownMenu>
          </CDropdown>

          {/* Garis Vertikal */}
          <li className="nav-item py-1">
            <div className="vr h-100 mx-2 text-body text-opacity-75"></div>
          </li>

          {/* Dropdown Mode Warna */}
          <CDropdown variant="nav-item" placement="bottom-end" className="px-1">
            <CDropdownToggle caret={false}>
              {colorMode === 'dark' ? (
                <CIcon icon={cilMoon} size="lg" />
              ) : colorMode === 'auto' ? (
                <CIcon icon={cilContrast} size="lg" />
              ) : (
                <CIcon icon={cilSun} size="lg" />
              )}
            </CDropdownToggle>
            <CDropdownMenu>
              <CDropdownItem
                active={colorMode === 'light'}
                className="d-flex align-items-center"
                as="button"
                type="button"
                onClick={() => setColorMode('light')}
              >
                <CIcon className="me-2" icon={cilSun} size="lg" /> Light
              </CDropdownItem>
              <CDropdownItem
                active={colorMode === 'dark'}
                className="d-flex align-items-center"
                as="button"
                type="button"
                onClick={() => setColorMode('dark')}
              >
                <CIcon className="me-2" icon={cilMoon} size="lg" /> Dark
              </CDropdownItem>
              <CDropdownItem
                active={colorMode === 'auto'}
                className="d-flex align-items-center"
                as="button"
                type="button"
                onClick={() => setColorMode('auto')}
              >
                <CIcon className="me-2" icon={cilContrast} size="lg" /> Auto
              </CDropdownItem>
            </CDropdownMenu>
          </CDropdown>

          {/* Garis Vertikal */}
          <li className="nav-item py-1">
            <div className="vr h-100 mx-2 text-body text-opacity-75"></div>
          </li>

          {/* User Dropdown */}
          <AppHeaderDropdown />
        </CHeaderNav>
      </CContainer>
      <CContainer className="px-4" fluid>
        <CBreadcrumb className="my-0">
          {breadCrumbs.map((bread, index) => (
            <CBreadcrumbItem
              key={index}
              active={bread.isActive}
              {...(!bread.isActive ? { href: bread.url } : {})}
            >
              {bread.name}
            </CBreadcrumbItem>
          ))}
        </CBreadcrumb>
      </CContainer>
    </CHeader>
  )
}

export default AppHeader
