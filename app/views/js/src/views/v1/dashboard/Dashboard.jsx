import React from 'react'
import { AppSidebar, AppFooter, AppHeader } from '../../../components/index'
import { CContainer } from '@coreui/react'
import { Helmet } from 'react-helmet'
import { useTranslation } from 'react-i18next'

const Dashboard = () => {
  const { t, i18n } = useTranslation()

  const breadCrumbs = [
    {
      name: 'Dashboard',
      url: '/dashboard',
      isActive: true,
    },
  ]

  return (
    <div>
      <Helmet>
        <title>{`${t('dashboard')} | ${import.meta.env.VITE_APP_NAME}`}</title>
      </Helmet>

      <AppSidebar />
      <div className="wrapper d-flex flex-column min-vh-100">
        <AppHeader breadCrumbs={breadCrumbs} />
        <div className="body flex-grow-1">
          <CContainer className="px-4" lg>
            Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum
            has been the industry's standard dummy text ever since the 1500s, when an unknown
            printer took a galley of type and scrambled it to make a type specimen book. It has
            survived not only five centuries, but also the leap into electronic typesetting,
            remaining essentially unchanged. It was popularised in the 1960s with the release of
            Letraset sheets containing Lorem Ipsum passages, and more recently with desktop
            publishing software like Aldus PageMaker including versions of Lorem Ipsum.
          </CContainer>
        </div>
        <AppFooter />
      </div>
    </div>
  )
}

export default Dashboard
