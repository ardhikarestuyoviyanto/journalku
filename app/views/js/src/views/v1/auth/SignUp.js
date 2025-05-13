import React, { useEffect, useState } from 'react'
import {
  CButton,
  CCard,
  CCardBody,
  CCol,
  CContainer,
  CForm,
  CFormInput,
  CInputGroup,
  CInputGroupText,
  CRow,
} from '@coreui/react'
import CIcon from '@coreui/icons-react'
import { cilLockLocked, cilPhone, cilUser } from '@coreui/icons'
import { useTranslation } from 'react-i18next'
import { Link, useNavigate } from 'react-router-dom'
import { useFormik } from 'formik'
import * as Yup from 'yup'
import { toast } from 'react-toastify'
import { useSelector } from 'react-redux'

const SignUp = () => {
  const { t, i18n } = useTranslation()
  const [loading, setLoading] = useState(false)
  const auth = useSelector((state) => state.auth)
  const navigate = useNavigate()

  const formik = useFormik({
    initialValues: {
      name: '',
      email: '',
      phoneNumber: '',
      password: '',
      passwordConfirm: '',
    },
    validationSchema: Yup.object({
      name: Yup.string().required(t('nameRequired')),
      email: Yup.string().required(t('emailRequired')).email(t('emailNotValid')),
      phoneNumber: Yup.string()
        .required(t('phoneNumberRequired'))
        .matches(/^\d{8,}$/, t('phoneNumberNumeric')),
      password: Yup.string().required(t('passwordRequired')),
      passwordConfirm: Yup.string()
        .required(t('passwordConfirmRequired'))
        .oneOf([Yup.ref('password')], t('passwordConfirmNotSame')),
    }),
    onSubmit: (values) => {
      setLoading(true)
      const formData = new FormData()
      formData.append('name', values.name)
      formData.append('email', values.email)
      formData.append('password', values.password)
      formData.append('phoneNumber', values.phoneNumber)

      fetch(`${import.meta.env.VITE_API_BASE_URL}/auth/signUp`, {
        method: 'POST',
        body: formData,
        headers: {
          'Accept-Lang': i18n.language,
        },
      })
        .then((res) => res.json())
        .then((data) => {
          if (data.success) {
            toast.success(data.message)
            formik.resetForm()
          } else {
            toast.error(data.error)
          }
        })
        .finally(() => {
          setLoading(false)
        })
    },
  })

  useEffect(() => {
    if (auth != undefined) {
      navigate('/choose-company')
    }
  }, [])

  return (
    <div className="bg-body-tertiary min-vh-100 d-flex flex-row align-items-center">
      <CContainer>
        <CRow className="justify-content-center">
          <CCol md={9} lg={7} xl={6}>
            <CCard className="mx-4">
              <CCardBody className="p-4">
                <CForm onSubmit={formik.handleSubmit}>
                  <h1>{t('registration')}</h1>
                  <p className="text-body-secondary">{t('createAccount')}</p>
                  <CInputGroup className="mb-3">
                    <CInputGroupText>
                      <CIcon icon={cilUser} />
                    </CInputGroupText>
                    <CFormInput
                      placeholder={t('userName')}
                      name="name"
                      onBlur={formik.handleBlur}
                      onChange={formik.handleChange}
                      value={formik.values.name}
                      invalid={formik.touched.name && !!formik.errors.name}
                      feedback={formik.touched.name && formik.errors.name}
                    />
                  </CInputGroup>
                  <CInputGroup className="mb-3">
                    <CInputGroupText>@</CInputGroupText>
                    <CFormInput
                      onBlur={formik.handleBlur}
                      onChange={formik.handleChange}
                      value={formik.values.email}
                      name="email"
                      placeholder={t('email')}
                      invalid={formik.touched.email && !!formik.errors.email}
                      feedback={formik.touched.email && formik.errors.email}
                    />
                  </CInputGroup>
                  <CInputGroup className="mb-3">
                    <CInputGroupText>
                      <CIcon icon={cilPhone} />
                    </CInputGroupText>
                    <CFormInput
                      onBlur={formik.handleBlur}
                      onChange={formik.handleChange}
                      value={formik.values.phoneNumber}
                      name="phoneNumber"
                      placeholder={t('phoneNumber')}
                      invalid={formik.touched.phoneNumber && !!formik.errors.phoneNumber}
                      feedback={formik.touched.phoneNumber && formik.errors.phoneNumber}
                    />
                  </CInputGroup>
                  <CInputGroup className="mb-3">
                    <CInputGroupText>
                      <CIcon icon={cilLockLocked} />
                    </CInputGroupText>
                    <CFormInput
                      type="password"
                      placeholder={t('password')}
                      onBlur={formik.handleBlur}
                      onChange={formik.handleChange}
                      value={formik.values.password}
                      name="password"
                      invalid={formik.touched.password && !!formik.errors.password}
                      feedback={formik.touched.password && formik.errors.password}
                    />
                  </CInputGroup>
                  <CInputGroup className="mb-4">
                    <CInputGroupText>
                      <CIcon icon={cilLockLocked} />
                    </CInputGroupText>
                    <CFormInput
                      type="password"
                      placeholder={t('confirmPassword')}
                      autoComplete="new-password"
                      onBlur={formik.handleBlur}
                      onChange={formik.handleChange}
                      value={formik.values.passwordConfirm}
                      name="passwordConfirm"
                      invalid={formik.touched.passwordConfirm && !!formik.errors.passwordConfirm}
                      feedback={formik.touched.passwordConfirm && formik.errors.passwordConfirm}
                    />
                  </CInputGroup>
                  <div className="d-grid">
                    <CButton disabled={loading} type="submit" color="primary">
                      {t('registration')}
                    </CButton>
                    <Link to="/" className="text-center mt-3">
                      {t('alreadyHaveAccount')}
                    </Link>
                  </div>
                </CForm>
              </CCardBody>
            </CCard>
          </CCol>
        </CRow>
      </CContainer>
    </div>
  )
}

export default SignUp
