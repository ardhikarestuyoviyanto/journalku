import { useSelector } from 'react-redux'
import React, { useEffect, useState } from 'react'
import {
  CButton,
  CCard,
  CCardBody,
  CCol,
  CContainer,
  CForm,
  CFormInput,
  CFormLabel,
  CFormTextarea,
  CImage,
  CRow,
} from '@coreui/react'
import { useTranslation } from 'react-i18next'
import { Link, useNavigate } from 'react-router-dom'
import { useFormik } from 'formik'
import * as Yup from 'yup'
import { toast } from 'react-toastify'
import Select from 'react-select'
import fetchWithAuth from '../../../helpers/fetch'
import history from '../../../history'

const defaultPhoto = `${import.meta.env.VITE_API_DOMAIN}/storage/image/empty_image.jpg`
const CreateCompany = () => {
  const { t, i18n } = useTranslation()
  const [photoPreview, setPhotoPreview] = useState(defaultPhoto)
  const [province, setProvince] = useState([])
  const [regency, setRegency] = useState([])
  const [subDistrict, setSubDistrict] = useState([])
  const [loading, setLoading] = useState(false)

  const handleGetProvince = async () => {
    const data = await fetchWithAuth(`${import.meta.env.VITE_API_BASE_URL}/metadata/province`, {
      method: 'GET',
    })
    if (data && data.success) {
      const provinceList = []
      data.data.forEach((item, index) => {
        provinceList.push({
          value: item.id,
          label: item.name,
        })
      })
      setProvince(provinceList)
    } else {
      toast.error(data.error)
    }
  }

  const handleGetRegency = async (provinceId) => {
    const data = await fetchWithAuth(
      `${import.meta.env.VITE_API_BASE_URL}/metadata/regency?provinceId=${provinceId}`,
      {
        method: 'GET',
      },
    )
    if (data && data.success) {
      const regencyList = []
      data.data.forEach((item, index) => {
        regencyList.push({
          value: item.id,
          label: item.name,
        })
      })
      setRegency(regencyList)
    } else {
      toast.error(data.error)
    }
  }

  const handleGetSubDistrict = async (regencyId) => {
    const data = await fetchWithAuth(
      `${import.meta.env.VITE_API_BASE_URL}/metadata/subDistrict?regencyId=${regencyId}`,
      {
        method: 'GET',
      },
    )
    if (data && data.success) {
      const subDistrictList = []
      data.data.forEach((item, index) => {
        subDistrictList.push({
          value: item.id,
          label: item.name,
        })
      })
      setSubDistrict(subDistrictList)
    } else {
      toast.error(data.error)
    }
  }

  const handleChangeProvince = (e) => {
    // Get Regency
    handleGetRegency(e.value)
    // Set Formik
    formik.setFieldValue('provinceId', e.value)
  }

  const handleChangeRegency = (e) => {
    // Get SubDistrict
    handleGetSubDistrict(e.value)
    // Set Formik
    formik.setFieldValue('regencyId', e.value)
  }

  const handleChangeSubDistrict = (e) => {
    // Set Formik
    formik.setFieldValue('subDistrictId', e.value)
  }

  const formik = useFormik({
    initialValues: {
      name: '',
      photo: '',
      provinceId: '',
      regencyId: '',
      subDistrictId: '',
      address: '',
    },
    validationSchema: Yup.object({
      name: Yup.string().required(t('companyRequired')),
      photo: Yup.mixed()
        .nullable()
        .optional()
        .test(
          'fileFormat',
          t('imageUploadRequitment'),
          (value) =>
            !value || (value && ['image/jpeg', 'image/png', 'image/jpg'].includes(value.type)),
        ),
      provinceId: Yup.string().required(t('provinceRequired')),
      regencyId: Yup.string().required(t('regencyRequired')),
      subDistrictId: Yup.string().required(t('subDistrictRequired')),
      address: Yup.string().required(t('addressRequired')),
    }),
    onSubmit: async (values) => {
      const formData = new FormData()
      formData.append('name', values.name)
      formData.append('photo', values.photo)
      formData.append('provinceId', values.provinceId)
      formData.append('regencyId', values.regencyId)
      formData.append('subDistrictId', values.subDistrictId)
      formData.append('address', values.address)

      setLoading(true)
      const response = await fetchWithAuth(`${import.meta.env.VITE_API_BASE_URL}/company`, {
        method: 'POST',
        body: formData,
      })
      setLoading(false)

      if (!response.success) {
        toast.error(response.error)
      } else {
        history.push('/choose-company', { message: response.message })
      }
    },
  })

  useEffect(() => {
    handleGetProvince()
  }, [])

  return (
    <>
      <div className="bg-body-tertiary min-vh-100 d-flex flex-row align-items-center">
        <CContainer>
          <CRow className="justify-content-center">
            <CCol md={9} lg={7} xl={6}>
              <CCard className="mx-4">
                <CCardBody className="p-4">
                  <h1>{t('createCompany')}</h1>
                  <div className="card-text text-muted mb-3">{t('formCompletion')}</div>
                  <div className="text-center">
                    <CImage rounded src={photoPreview} width={100} />
                  </div>
                  <CForm onSubmit={formik.handleSubmit}>
                    <div className="mb-2">
                      <CFormLabel>
                        {t('company')} <span className="text-danger">*</span>
                      </CFormLabel>
                      <CFormInput
                        onChange={formik.handleChange}
                        onBlur={formik.handleBlur}
                        value={formik.values.name}
                        invalid={formik.touched.name && !!formik.errors.name}
                        feedback={formik.touched.name && formik.errors.name}
                        placeholder={t('company')}
                        name="name"
                      />
                    </div>

                    <div className="mb-2">
                      <CFormLabel>{t('companyPhoto')}</CFormLabel>
                      <CFormInput
                        accept="image/*"
                        onChange={(e) => {
                          const file = e.target.files[0]
                          formik.setFieldValue('photo', file)
                          // Show Preview Photo
                          if (file.type.startsWith('image/')) {
                            setPhotoPreview(URL.createObjectURL(file))
                          }
                        }}
                        onBlur={formik.handleBlur}
                        invalid={formik.touched.photo && !!formik.errors.photo}
                        feedback={formik.touched.photo && formik.errors.photo}
                        name="photo"
                        type="file"
                      />
                    </div>

                    <div className="mb-2">
                      <CFormLabel>
                        {t('province')} <span className="text-danger">*</span>
                      </CFormLabel>
                      <Select
                        options={province}
                        name="provinceId"
                        onChange={(e) => {
                          handleChangeProvince(e)
                          setSubDistrict([])
                        }}
                        placeholder={t('province')}
                      />
                      {formik.touched.provinceId && formik.errors.provinceId ? (
                        <div className="text-danger">{formik.errors.provinceId}</div>
                      ) : null}
                    </div>

                    <div className="mb-2">
                      <CFormLabel>
                        {t('regency')} <span className="text-danger">*</span>
                      </CFormLabel>
                      <Select
                        onChange={(e) => {
                          handleChangeRegency(e)
                        }}
                        options={regency}
                        placeholder={t('regency')}
                        name="regencyId"
                      />
                      {formik.touched.regencyId && formik.errors.regencyId ? (
                        <div className="text-danger">{formik.errors.regencyId}</div>
                      ) : null}
                    </div>

                    <div className="mb-2">
                      <CFormLabel>
                        {t('subDistrict')} <span className="text-danger">*</span>
                      </CFormLabel>
                      <Select
                        options={subDistrict}
                        placeholder={t('subDistrict')}
                        name="subDistrictId"
                        onChange={(e) => {
                          handleChangeSubDistrict(e)
                        }}
                      />
                      {formik.touched.subDistrictId && formik.errors.subDistrictId ? (
                        <div className="text-danger">{formik.errors.subDistrictId}</div>
                      ) : null}
                    </div>

                    <div className="mb-2">
                      <CFormLabel>
                        {t('addressFull')} <span className="text-danger">*</span>
                      </CFormLabel>
                      <CFormTextarea
                        name="address"
                        rows={3}
                        placeholder={t('addressFull')}
                        onChange={formik.handleChange}
                        onBlur={formik.handleBlur}
                        value={formik.values.address}
                        invalid={formik.touched.address && !!formik.errors.address}
                        feedback={formik.touched.address && formik.errors.address}
                      ></CFormTextarea>
                    </div>

                    <div className="d-grid mt-3">
                      <CButton disabled={loading} type="submit" color="primary">
                        {t('submit')}
                      </CButton>
                      <Link to="/choose-company" className="text-center mt-3">
                        {t('listCompany')}
                      </Link>
                    </div>
                  </CForm>
                </CCardBody>
              </CCard>
            </CCol>
          </CRow>
        </CContainer>
      </div>
    </>
  )
}

export default CreateCompany
