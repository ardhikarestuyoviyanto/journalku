import { CNavGroup, CNavItem } from '@coreui/react'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { library } from '@fortawesome/fontawesome-svg-core'
import * as solidIcons from '@fortawesome/free-solid-svg-icons'

// Tambahkan semua ikon ke library
const validIcons = Object.values(solidIcons).filter(
  (icon) => icon && typeof icon === 'object' && 'iconName' in icon && 'prefix' in icon,
)
library.add(...validIcons)

const convertMenuToNav = (menuList, language = 'id') => {
  return menuList
    .sort((a, b) => a.order - b.order)
    .map((item) => {
      const name = language === 'en' ? item.nameEn : item.nameId

      // Dapatkan nama ikon PascalCase, contoh: faGaugeHigh
      const iconDef = solidIcons[item.icon]

      const icon = iconDef ? <FontAwesomeIcon icon={iconDef} className="nav-icon" /> : null

      if (item.child && item.child.length > 0) {
        return {
          component: CNavGroup,
          name,
          to: item.url,
          icon,
          items: convertMenuToNav(item.child, language),
        }
      } else {
        return {
          component: CNavItem,
          name,
          to: item.url,
          icon,
        }
      }
    })
}

const menuRaw = JSON.parse(localStorage.getItem('auth'))?.user?.currentCompany?.menu || []
const language = localStorage.getItem('lang') || 'id'
const _nav = convertMenuToNav(menuRaw, language)

export default _nav
