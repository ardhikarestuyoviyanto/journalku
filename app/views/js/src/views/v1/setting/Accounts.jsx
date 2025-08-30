import React, { useEffect, useState } from 'react'
import { AppSidebar, AppFooter, AppHeader } from '../../../components/index'
import {
  CBadge,
  CButton,
  CCard,
  CCardBody,
  CCardHeader,
  CCol,
  CContainer,
  CForm,
  CFormInput,
  CFormLabel,
  CFormSelect,
  CFormTextarea,
  CInputGroup,
  CInputGroupText,
  CModal,
  CModalBody,
  CModalFooter,
  CModalHeader,
  CModalTitle,
  CRow,
  CSpinner,
  CTable,
  CTableBody,
  CTableDataCell,
  CTableHead,
  CTableHeaderCell,
  CTableRow,
  CTooltip,
} from '@coreui/react'
import { Helmet } from 'react-helmet'
import { useTranslation } from 'react-i18next'
import fetchWithAuth from '../../../helpers/fetch'
import checkPermission from '../../../helpers/globalHelper'
import CIcon from '@coreui/icons-react'
import {
  cilArrowBottom,
  cilArrowTop,
  cilPencil,
  cilPlus,
  cilSave,
  cilTrash,
  cilX,
} from '@coreui/icons'
import Select from 'react-select'
import PaginationTable from '../../../components/PaginationTable'
import { useFormik } from 'formik'
import * as Yup from 'yup'
import Select2CustomTheme from '../../../components/Select2CustomTheme'

const Accounts = () => {
  const { t, i18n } = useTranslation()
  const [data, setData] = useState([])
  const [totalPage, setTotalPage] = useState(0)
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [totalRows, setTotalRows] = useState(0)
  const [search, setSearch] = useState('')
  const [sortBy, setSortBy] = useState('account.number_account')
  const [order, setOrder] = useState('asc')
  const [perPage] = useState(25)
  const params = new URLSearchParams({
    page: page,
    perPage: perPage,
    order: order,
    sortBy: sortBy,
    search: search,
  })
  const [showModalForm, setShowModalForm] = useState(false)
  const [categoryAccountData, setCategoryAccountData] = useState([])
  const [codeCategoryAccount, setCodeCategoryAccount] = useState('')
  const [loadingForm, setLoadingForm] = useState(false)
  const statusArchive = [
    {
      label: t('active'),
      value: 0,
    },
    {
      label: t('archive'),
      value: 1,
    },
  ]

  const formik = useFormik({
    initialValues: {
      id: '',
      categoryAccountId: '',
      numberAccount: '',
      name: '',
      statusArchive: '',
      description: '',
    },
    validationSchema: Yup.object({
      categoryAccountId: Yup.string().required(t('categoryAccountRequired')),
      numberAccount: Yup.string()
        .required(t('numberAccountRequired'))
        .max(3, t('numberAccountMax3Digit')),
      name: Yup.string().required(t('nameAccountRequired')),
      statusArchive: Yup.number().required(t('statusArchiveRequired')),
      description: Yup.mixed().nullable().optional(),
    }),
    onSubmit: async (values) => {
      const formData = new FormData()
      formData.append('categoryAccountId', values.categoryAccountId)
      formData.append('numberAccount', `${codeCategoryAccount}${values.numberAccount}`)
      formData.append('name', values.name)
      formData.append('statusArchive', values.statusArchive)
      formData.append('description', values.description)

      setLoadingForm(true)
      if (values.id === '') {
        // Create Account
        const response = await fetchWithAuth(`${import.meta.env.VITE_API_BASE_URL}/accounts`, {
          method: 'POST',
          body: formData,
        })
      } else {
        // Update Account
      }
      setLoadingForm(false)
    },
  })

  const handleGetData = async function () {
    setLoading(true)
    const data = await fetchWithAuth(
      `${import.meta.env.VITE_API_BASE_URL}/accounts?${params.toString()}`,
      {
        method: 'GET',
      },
    )

    if (data.success) {
      setData(data.data.accounts)
      setTotalPage(data.data.lastPage)
      setTotalRows(data.data.totalRows)
    } else {
      toast.error(response.error)
    }
    setLoading(false)
  }

  const handleSort = (field) => {
    const orderRes = sortBy === field && order === 'asc' ? 'desc' : 'asc'
    setSortBy(field)
    setOrder(orderRes)
  }

  const handlePageChange = (page) => {
    setPage(page)
  }

  const handleGetCategoryAccount = async function () {
    const data = await fetchWithAuth(
      `${import.meta.env.VITE_API_BASE_URL}/dropdown/categoryAccounts`,
      {
        method: 'GET',
      },
    )

    if (data.success) {
      const categoryAccountList = []
      data.data.forEach((item, index) => {
        categoryAccountList.push({
          value: item.id,
          label: item.category,
          code: item.code,
        })
      })
      setCategoryAccountData(categoryAccountList)
    } else {
      toast.error(response.error)
    }
  }

  const handleSetCodeCategoryAccount = (index) => {
    const codeCategoryAccount = categoryAccountData[index]?.code
    setCodeCategoryAccount(codeCategoryAccount)
  }

  const breadCrumbs = [
    {
      name: t('accounts'),
      url: '/setting/accounts',
      isActive: true,
    },
  ]

  useEffect(() => {
    handleGetData()
    handleGetCategoryAccount()
  }, [page, search, sortBy, order])

  return (
    <div>
      <Helmet>
        <title>{`${t('accounts')} | ${import.meta.env.VITE_APP_NAME}`}</title>
      </Helmet>

      <AppSidebar />
      <div className="wrapper d-flex flex-column min-vh-100">
        <AppHeader breadCrumbs={breadCrumbs} />
        <div className="body flex-grow-1">
          <CContainer className="px-4" lg>
            <CCard>
              {checkPermission('createAccount') && (
                <CCardHeader>
                  <CButton
                    onClick={() => {
                      setShowModalForm(true)
                      formik.resetForm()
                    }}
                    color="primary"
                    size="sm"
                    style={{ float: 'right' }}
                  >
                    <CIcon icon={cilPlus} /> {t('create')}
                  </CButton>
                </CCardHeader>
              )}

              <CCardBody>
                <CRow className="justify-content-end mb-3">
                  <CCol sm={3}>
                    <CFormInput
                      onChange={(e) => setSearch(e.target.value)}
                      placeholder={t('search')}
                    />
                  </CCol>
                </CRow>
                <CTable striped responsive>
                  <CTableHead>
                    <CTableRow>
                      <CTableHeaderCell
                        onClick={() => {
                          handleSort('account.created_at')
                        }}
                        scope="col"
                        style={{ cursor: 'pointer', width: 70 }}
                      >
                        {t('no')}{' '}
                        <CIcon
                          icon={
                            order == 'desc' && sortBy == 'account.created_at'
                              ? cilArrowBottom
                              : cilArrowTop
                          }
                        />
                      </CTableHeaderCell>
                      <CTableHeaderCell
                        scope="col"
                        onClick={() => {
                          handleSort('account.number_account')
                        }}
                        style={{ cursor: 'pointer' }}
                      >
                        {t('code')}{' '}
                        <CIcon
                          icon={
                            order == 'desc' && sortBy == 'account.number_account'
                              ? cilArrowBottom
                              : cilArrowTop
                          }
                        />
                      </CTableHeaderCell>
                      <CTableHeaderCell
                        scope="col"
                        onClick={() => {
                          handleSort('account.name')
                        }}
                        style={{ cursor: 'pointer' }}
                      >
                        {t('name')}{' '}
                        <CIcon
                          icon={
                            order == 'desc' && sortBy == 'account.name'
                              ? cilArrowBottom
                              : cilArrowTop
                          }
                        />
                      </CTableHeaderCell>
                      <CTableHeaderCell
                        scope="col"
                        onClick={() => {
                          handleSort('account.category_account_id')
                        }}
                        style={{ cursor: 'pointer' }}
                      >
                        {t('category')}{' '}
                        <CIcon
                          icon={
                            order == 'desc' && sortBy == 'account.category_account_id'
                              ? cilArrowBottom
                              : cilArrowTop
                          }
                        />
                      </CTableHeaderCell>
                      <CTableHeaderCell
                        scope="col"
                        onClick={() => {
                          handleSort('account.description')
                        }}
                        style={{ cursor: 'pointer' }}
                      >
                        {t('description')}{' '}
                        <CIcon
                          icon={
                            order == 'desc' && sortBy == 'account.description'
                              ? cilArrowBottom
                              : cilArrowTop
                          }
                        />
                      </CTableHeaderCell>
                      <CTableHeaderCell
                        scope="col"
                        onClick={() => {
                          handleSort('account.status_archive')
                        }}
                        style={{ cursor: 'pointer' }}
                      >
                        {t('status')}{' '}
                        <CIcon
                          icon={
                            order == 'desc' && sortBy == 'account.status_archive'
                              ? cilArrowBottom
                              : cilArrowTop
                          }
                        />
                      </CTableHeaderCell>
                      <CTableHeaderCell scope="col">{t('status')}</CTableHeaderCell>
                    </CTableRow>
                  </CTableHead>
                  <CTableBody>
                    {loading ? (
                      <tr>
                        <td colSpan={7}>
                          <div className="d-flex align-items-center">
                            <div role="status">Loading...</div>
                            <CSpinner color="primary" size="sm" className="ms-auto" />
                          </div>
                        </td>
                      </tr>
                    ) : data.length == 0 ? (
                      <tr>
                        <td colSpan={7}>{t('emptyData')}</td>
                      </tr>
                    ) : (
                      data.map((item, index) => (
                        <tr key={item.id}>
                          <td>{(page - 1) * perPage + index + 1}</td>
                          <td>{item.number_account}</td>
                          <td>{item.name}</td>
                          <td>{item.category_account}</td>
                          <td>{item.description}</td>
                          <td>
                            {item.is_archive == 0 ? (
                              <CBadge color="success" shape="rounded-pill">
                                {t('active')}
                              </CBadge>
                            ) : (
                              <CBadge color="danger" shape="rounded-pill">
                                {t('archive')}
                              </CBadge>
                            )}
                          </td>
                          <td>
                            {checkPermission('updateAccount') && item.is_primary == 0 && (
                              <CTooltip content={t('update')}>
                                <CButton color="primary" size="sm">
                                  <CIcon icon={cilPencil} />
                                </CButton>
                              </CTooltip>
                            )}{' '}
                            {checkPermission('deleteAccount') && item.is_primary == 0 && (
                              <CTooltip content={t('delete')}>
                                <CButton color="danger" className="text-white" size="sm">
                                  <CIcon icon={cilTrash} />
                                </CButton>
                              </CTooltip>
                            )}
                          </td>
                        </tr>
                      ))
                    )}
                  </CTableBody>
                </CTable>
                <PaginationTable
                  onPageChange={handlePageChange}
                  totalPages={totalPage}
                  currentPage={page}
                  totalRows={totalRows}
                />
              </CCardBody>
            </CCard>
          </CContainer>
        </div>
        <AppFooter />
      </div>

      <CModal size="lg" visible={showModalForm} onClose={() => setShowModalForm(false)}>
        <CModalHeader>
          <CModalTitle id="LiveDemoExampleLabel">
            {t('create')} {t('accounts')}
          </CModalTitle>
        </CModalHeader>
        <CForm onSubmit={formik.handleSubmit}>
          <CModalBody>
            <CRow>
              <CCol sm={6}>
                <CFormLabel className="mt-1">
                  {t('categoryAccount')} <span className="text-danger">*</span>
                </CFormLabel>
                <Select2CustomTheme
                  name="categoryAccountId"
                  options={categoryAccountData}
                  placeholder={t('categoryAccount')}
                  onChange={(e) => {
                    const selectedIndex = categoryAccountData.findIndex(
                      (item) => item.value === e.value,
                    )
                    handleSetCodeCategoryAccount(selectedIndex)
                    formik.setFieldValue('categoryAccountId', e.value)
                  }}
                />
                {formik.touched.categoryAccountId && formik.errors.categoryAccountId ? (
                  <div className="text-danger">{formik.errors.categoryAccountId}</div>
                ) : null}
              </CCol>
              <CCol sm={6}>
                <CFormLabel className="mt-1">
                  {t('numberAccount')} <span className="text-danger">*</span>
                </CFormLabel>
                <CInputGroup className="mb-3">
                  <CInputGroupText>{codeCategoryAccount}</CInputGroupText>
                  <CFormInput
                    name="numberAccount"
                    placeholder={t('numberAccount')}
                    onChange={formik.handleChange}
                    onBlur={formik.handleBlur}
                    value={formik.values.numberAccount}
                    invalid={formik.touched.numberAccount && !!formik.errors.numberAccount}
                    feedback={formik.touched.numberAccount && formik.errors.numberAccount}
                  />
                </CInputGroup>
              </CCol>
              <CCol sm={6}>
                <CFormLabel className="mt-1">
                  {t('accounName')} <span className="text-danger">*</span>
                </CFormLabel>
                <CFormInput
                  name="name"
                  placeholder={t('accounName')}
                  onChange={formik.handleChange}
                  onBlur={formik.handleBlur}
                  value={formik.values.name}
                  invalid={formik.touched.name && !!formik.errors.name}
                  feedback={formik.touched.name && formik.errors.name}
                />
              </CCol>
              <CCol sm={6}>
                <CFormLabel className="mt-1">
                  {t('status')} <span className="text-danger">*</span>
                </CFormLabel>
                <Select2CustomTheme
                  name="statusArchive"
                  options={statusArchive}
                  onChange={(e) => {
                    formik.setFieldValue('statusArchive', e.value)
                  }}
                  placeholder={t('status')}
                />
                {formik.touched.statusArchive && formik.errors.statusArchive ? (
                  <div className="text-danger">{formik.errors.statusArchive}</div>
                ) : null}
              </CCol>
              <CCol sm={6}>
                <CFormLabel className="mt-3">{t('description')}</CFormLabel>
                <CFormTextarea
                  name="description"
                  onChange={formik.handleChange}
                  onBlur={formik.handleBlur}
                  value={formik.values.description}
                  invalid={formik.touched.description && !!formik.errors.description}
                  feedback={formik.touched.description && formik.errors.description}
                  placeholder={t('description')}
                  rows={3}
                />
              </CCol>
            </CRow>
          </CModalBody>
          <CModalFooter>
            <CButton color="secondary" onClick={() => setShowModalForm(false)}>
              <CIcon icon={cilX} /> {t('close')}
            </CButton>
            <CButton type="submit" disabled={loadingForm} color="primary">
              <CIcon icon={cilSave} /> {t('submit')}
            </CButton>
          </CModalFooter>
        </CForm>
      </CModal>
    </div>
  )
}

export default Accounts
