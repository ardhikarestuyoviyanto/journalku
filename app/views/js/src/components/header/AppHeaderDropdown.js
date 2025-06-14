import React from 'react'
import {
  CAvatar,
  CDropdown,
  CDropdownDivider,
  CDropdownHeader,
  CDropdownItem,
  CDropdownMenu,
  CDropdownToggle,
} from '@coreui/react'
import { cilSettings, cilUser, cilAccountLogout } from '@coreui/icons'
import CIcon from '@coreui/icons-react'
import store from '../../store'
import avatar8 from './../../assets/images/avatars/8.jpg'
import { useTranslation } from 'react-i18next'
import { useDispatch } from 'react-redux'
import { useNavigate } from 'react-router-dom'
import history from '../../history'

const AppHeaderDropdown = () => {
  const state = store.getState()
  const auth = state.auth
  const { t, i18n } = useTranslation()
  const dispatch = useDispatch()

  const handleLogout = () => {
    history.push('/', { logoutSuccess: t('logoutSuccess') })
    dispatch({ type: 'logout' })
  }

  return (
    <CDropdown variant="nav-item">
      <CDropdownToggle placement="bottom-end" className="py-0 pe-0" caret={false}>
        <CAvatar src={avatar8} size="md" />
      </CDropdownToggle>
      <CDropdownMenu className="pt-0" placement="bottom-end">
        <CDropdownHeader className="bg-body-secondary fw-semibold mb-2">
          Hi, {auth.user.name}
        </CDropdownHeader>
        <CDropdownItem href="#">
          <CIcon icon={cilUser} className="me-2" />
          Profile
        </CDropdownItem>
        <CDropdownItem href="#">
          <CIcon icon={cilSettings} className="me-2" />
          Settings
        </CDropdownItem>
        <CDropdownDivider />
        <CDropdownItem
          href="#"
          onClick={() => {
            handleLogout()
          }}
        >
          <CIcon icon={cilAccountLogout} className="me-2" />
          {t('logout')}
        </CDropdownItem>
      </CDropdownMenu>
    </CDropdown>
  )
}

export default AppHeaderDropdown
