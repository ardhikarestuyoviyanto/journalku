import React, { useEffect, useState } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import {
  CButton,
  CCard,
  CCardBody,
  CCardGroup,
  CCardText,
  CCol,
  CContainer,
  CForm,
  CFormCheck,
  CFormFeedback,
  CFormInput,
  CFormLabel,
  CImage,
  CInputGroup,
  CInputGroupText,
  CRow,
} from '@coreui/react'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { useTranslation } from 'react-i18next'
import { faEye, faEyeSlash } from '@fortawesome/free-solid-svg-icons'
import { useFormik } from 'formik'
import * as Yup from 'yup'
import { toast } from 'react-toastify'
import { useDispatch, useSelector } from 'react-redux'

const SignIn = () => {
  const { t, i18n } = useTranslation()
  const [showPassword, setShowPassword] = useState(false)
  const [rememberMe, setRememberMe] = useState(false)
  const [loading, setLoading] = useState(false)
  const dispatch = useDispatch()
  const auth = useSelector((state) => state.auth)
  const navigate = useNavigate()
  const location = useLocation()
  const unauthorized = location.state?.unauthorized

  const formik = useFormik({
    initialValues: {
      email: '',
      password: '',
    },
    validationSchema: Yup.object({
      email: Yup.string().required(t('emailRequired')).email(t('emailNotValid')),
      password: Yup.string().required(t('passwordRequired')),
    }),
    onSubmit: (values) => {
      setLoading(true)
      const formData = new FormData()
      formData.append('email', values.email)
      formData.append('password', values.password)
      formData.append('rememberMe', rememberMe)

      fetch(`${import.meta.env.VITE_API_BASE_URL}/auth/signIn`, {
        method: 'POST',
        body: formData,
        headers: {
          'Accept-Lang': i18n.language,
        },
      })
        .then((res) => res.json())
        .then((res) => {
          if (res.success) {
            dispatch({
              type: 'set',
              auth: res.data,
            })
            navigate('/choose-company')
          } else {
            toast.error(res.error)
          }
        })
        .finally(() => {
          setLoading(false)
        })
    },
  })

  const handlePopUpSignInGoogle = () => {
    window.open(
      `${import.meta.env.VITE_API_BASE_URL}/auth/google/signIn`,
      'Login with Google',
      'width=500,height=600',
    )
  }

  useEffect(() => {
    if (auth != undefined) {
      navigate('/choose-company')
    }

    if (unauthorized) {
      toast.error(t('unAuthorized'))
    }

    const handleMessage = (event) => {
      if (event.origin !== import.meta.env.VITE_API_DOMAIN) return
      try {
        const res = event.data
        if (res.success) {
          dispatch({
            type: 'set',
            auth: res.data,
          })
          navigate('/choose-company')
        } else {
          toast.error(res.error)
        }
      } catch (err) {
        console.error('Invalid JSON from popup', err)
      }
    }

    window.addEventListener('message', handleMessage)
    return () => window.removeEventListener('message', handleMessage)
  }, [])

  return (
    <div className="bg-body-tertiary min-vh-100 d-flex flex-row align-items-center">
      <CContainer>
        <CRow className="justify-content-center">
          <CCol md={8}>
            <CCardGroup>
              <CCard className="p-4">
                <CCardBody>
                  <CForm onSubmit={formik.handleSubmit}>
                    <h1>{t('signIn')}</h1>
                    <p className="text-body-secondary">{t('signInAccount')}</p>
                    <CFormInput
                      className="mb-2"
                      label={t('email')}
                      placeholder={t('email')}
                      onBlur={formik.handleBlur}
                      onChange={formik.handleChange}
                      value={formik.values.email}
                      invalid={formik.touched.email && !!formik.errors.email}
                      feedback={formik.touched.email && formik.errors.email}
                      name="email"
                    />
                    <CFormLabel>{t('password')}</CFormLabel>
                    <CInputGroup className="mb-2">
                      <CFormInput
                        type={showPassword ? 'text' : 'password'}
                        placeholder={t('password')}
                        name="password"
                        onBlur={formik.handleBlur}
                        onChange={formik.handleChange}
                        value={formik.values.password}
                        invalid={formik.touched.password && !!formik.errors.password}
                      />
                      <CInputGroupText
                        onClick={() => setShowPassword(!showPassword)}
                        style={{ cursor: 'pointer' }}
                      >
                        <FontAwesomeIcon icon={showPassword ? faEye : faEyeSlash} />
                      </CInputGroupText>
                    </CInputGroup>
                    {formik.touched.password && formik.errors.password && (
                      <CFormFeedback className="d-block text-danger" type="invalid">
                        {formik.errors.password}
                      </CFormFeedback>
                    )}
                    <div className="mb-4 mt-3">
                      <CFormCheck
                        id="rememberMe"
                        name="rememberMe"
                        checked={rememberMe}
                        onChange={() => setRememberMe(!rememberMe)}
                        label={t('rememberMe')}
                        htmlFor="rememberMe"
                      />
                    </div>
                    <CRow>
                      <CCol xs={6}>
                        <CButton disabled={loading} color="primary" type="submit" className="px-4">
                          {t('signIn')}
                        </CButton>
                      </CCol>
                      <CCol xs={6} className="text-right" style={{ textAlign: 'right' }}>
                        <CButton color="link" className="px-0">
                          {t('forgotPassword')}
                        </CButton>
                      </CCol>
                    </CRow>
                  </CForm>
                  <CCardText className="mt-3 mb-3 text-muted text-center">
                    <small>{t('orLoginWith')}</small>
                  </CCardText>
                  <div className="text-center">
                    <CButton color="light" onClick={() => handlePopUpSignInGoogle()}>
                      <CImage
                        className="mb-1 mt-1"
                        width={20}
                        src={'src/assets/images/google.webp'}
                      />
                    </CButton>
                  </div>
                </CCardBody>
              </CCard>

              <CCard className="text-white bg-primary py-5" style={{ width: '44%' }}>
                <CCardBody className="text-center">
                  <div>
                    <blockquote className="blockquote">{t('slogan1')}</blockquote>
                    <p className="mt-5 mb-5 lead">{t('slogan2')}</p>
                    <Link to="/signup">
                      <CButton color="primary" className="mt-3" active tabIndex={-1}>
                        {t('signUpNow')}
                      </CButton>
                    </Link>
                  </div>
                </CCardBody>
              </CCard>
            </CCardGroup>
          </CCol>
        </CRow>
      </CContainer>
    </div>
  )
}

export default SignIn
