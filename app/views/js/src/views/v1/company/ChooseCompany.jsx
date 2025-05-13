import { useSelector } from 'react-redux'
import React, { useEffect, useState } from 'react'
import {
  CAvatar,
  CBadge,
  CButton,
  CCard,
  CCardBody,
  CCol,
  CContainer,
  CListGroup,
  CListGroupItem,
  CPlaceholder,
  CRow,
  CSpinner,
} from '@coreui/react'
import { useTranslation } from 'react-i18next'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { toast } from 'react-toastify'
import fetchWithAuth from '../../../helpers/fetch'
import { CAlert } from '@coreui/react'
import history from '../../../history'

const defaultPhoto = `${import.meta.env.VITE_API_DOMAIN}/storage/image/empty_image.jpg`

const ChooseCompany = () => {
  const auth = useSelector((state) => state.auth)
  const { t, i18n } = useTranslation()
  const navigate = useNavigate()
  const location = useLocation()
  const message = location.state?.message

  const [company, setCompany] = useState([])
  const [loading, setLoading] = useState(false)

  const handleGetCompany = async () => {
    const data = await fetchWithAuth(`${import.meta.env.VITE_API_BASE_URL}/company`, {
      method: 'GET',
    })

    if (data.success) {
      setCompany(data.data)
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
      alert('Sukses')
      console.log(response)
    } else {
      toast.error(response.error)
    }
  }

  useEffect(() => {
    if (message != undefined) {
      toast.success(message)
      history.replace({ ...location, state: { ...location.state, message: undefined } })
    }

    handleGetCompany()
  }, [])

  return (
    <>
      <div className="bg-body-tertiary min-vh-100 d-flex flex-row align-items-center">
        <CContainer>
          <CRow className="justify-content-center">
            <CCol md={9} lg={7} xl={6}>
              <CCard className="mx-4">
                <CCardBody className="p-4">
                  <h1>{t('listCompany')}</h1>
                  <div className="card-text text-muted">{t('createCompanyCaption1')}</div>

                  {company.length == 0 ? (
                    <CAlert color="danger" className="d-flex align-items-center mt-3">
                      <svg
                        className="flex-shrink-0 me-2"
                        width="24"
                        height="24"
                        viewBox="0 0 512 512"
                      >
                        <rect
                          width="32"
                          height="176"
                          x="240"
                          y="176"
                          fill="var(--ci-primary-color, currentColor)"
                          className="ci-primary"
                        ></rect>
                        <rect
                          width="32"
                          height="32"
                          x="240"
                          y="384"
                          fill="var(--ci-primary-color, currentColor)"
                          className="ci-primary"
                        ></rect>
                        <path
                          fill="var(--ci-primary-color, currentColor)"
                          d="M274.014,16H237.986L16,445.174V496H496V445.174ZM464,464H48V452.959L256,50.826,464,452.959Z"
                          className="ci-primary"
                        ></path>
                      </svg>
                      <div>{t('emptyCompany')}</div>
                    </CAlert>
                  ) : (
                    <>
                      {loading ? (
                        <div className="d-flex align-items-center mt-5">
                          <strong role="status">{t('Loading')}</strong>
                          <CSpinner className="ms-auto" color="primary" />
                        </div>
                      ) : (
                        <CListGroup className="mt-3 mb-3">
                          {company.map((comp, index) => (
                            <CListGroupItem
                              as="a"
                              href="#"
                              onClick={() => {
                                handleChooseCompany(comp.id)
                              }}
                              key={index}
                              className="border-0 shadow-sm rounded-4 mb-2 p-3 list-group-hover"
                              style={{ transition: '0.2s', background: 'blur' }}
                            >
                              <CRow>
                                <CCol
                                  sm={2}
                                  className="d-flex align-items-center justify-content-center"
                                >
                                  <CAvatar
                                    shape="rounded"
                                    size="xl"
                                    src={
                                      comp.photo == null
                                        ? defaultPhoto
                                        : `${import.meta.env.VITE_API_DOMAIN}/storage/image/${comp.photo}`
                                    }
                                  />
                                </CCol>
                                <CCol sm={10}>
                                  <div className="d-flex align-items-center justify-content-between">
                                    <h5 className="mb-0 text-dark fw-semibold">{comp.name}</h5>
                                    <CBadge color="danger" className="ms-2">
                                      <i className="fa fa-user-shield me-1"></i> {comp.role}
                                    </CBadge>
                                  </div>
                                  <small className="text-muted">{comp.address}</small>
                                </CCol>
                              </CRow>
                            </CListGroupItem>
                          ))}
                        </CListGroup>
                      )}
                    </>
                  )}

                  <div className="d-grid mt-2">
                    <Link to={'/choose-company/create'}>
                      <CButton type="button" color="primary">
                        {t('createCompany')}
                      </CButton>
                    </Link>
                  </div>
                </CCardBody>
              </CCard>
            </CCol>
          </CRow>
        </CContainer>
      </div>
    </>
  )
}

export default ChooseCompany
