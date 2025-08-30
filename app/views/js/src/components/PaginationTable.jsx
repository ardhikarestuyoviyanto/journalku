import { cilArrowLeft, cilArrowRight, cilArrowThickLeft, cilArrowThickRight } from '@coreui/icons'
import CIcon from '@coreui/icons-react'
import { CPagination, CPaginationItem } from '@coreui/react'

const PaginationTable = ({ totalPages, currentPage, onPageChange, totalRows = null }) => {
  const generatePageNumbers = () => {
    const pages = []
    if (currentPage > 3) pages.push(1, '...')
    for (let i = Math.max(1, currentPage - 2); i <= Math.min(totalPages, currentPage + 2); i++) {
      pages.push(i)
    }
    if (currentPage < totalPages - 2) pages.push('...', totalPages)
    return pages
  }

  return (
    <div className="d-flex justify-content-between align-items-center flex-wrap mt-2">
      {totalRows !== null && (
        <small className="text-muted">
          Showing {(currentPage - 1) * 25 + 1} - {Math.min(currentPage * 25, totalRows)} of{' '}
          {totalRows} Data
        </small>
      )}

      <CPagination size="sm" aria-label="Table pagination">
        <CPaginationItem disabled={currentPage === 1} onClick={() => onPageChange(1)}>
          <CIcon icon={cilArrowThickLeft} />
        </CPaginationItem>

        <CPaginationItem disabled={currentPage === 1} onClick={() => onPageChange(currentPage - 1)}>
          <CIcon icon={cilArrowLeft} />
        </CPaginationItem>

        {generatePageNumbers().map((page, index) => (
          <CPaginationItem
            key={index}
            active={page === currentPage}
            onClick={() => Number.isInteger(page) && onPageChange(page)}
            disabled={!Number.isInteger(page)}
          >
            {page}
          </CPaginationItem>
        ))}

        <CPaginationItem
          disabled={currentPage === totalPages || totalPages === 0}
          onClick={() => onPageChange(currentPage + 1)}
        >
          <CIcon icon={cilArrowRight} />
        </CPaginationItem>

        <CPaginationItem
          disabled={currentPage === totalPages || totalPages === 0}
          onClick={() => onPageChange(totalPages)}
        >
          <CIcon icon={cilArrowThickRight} />
        </CPaginationItem>
      </CPagination>
    </div>
  )
}

export default PaginationTable
