import React from 'react'
import PropTypes from 'prop-types'
import { myMarketSearch } from '@/service/api'
import TablePagination from '@/components/TablePagination'

const initialDatatable = {
  data: [],
  result_count: 0,
  total_count: 0,
  loading: false,
  error: null,
}

const datatableBaseFilter = {
  sort: 'updated_at:desc',
  page: 1,
}

function withDataFetch(Component, initFilter, searchFn = myMarketSearch) {
  function wrapped(props) {
    const { filter: propFilter, onReload } = props

    const [data, setData] = React.useState(initialDatatable)
    const [filter, setFilter] = React.useState({
      ...datatableBaseFilter,
      ...initFilter,
      ...propFilter,
    })
    const [tick, setTick] = React.useState(false)

    React.useEffect(() => {
      ;(async () => {
        setData({ ...data, loading: true, error: null })
        try {
          const res = await searchFn(filter)
          setData({ ...data, loading: false, ...res })
        } catch (e) {
          setData({ ...data, loading: false, error: e.message })
        }
      })()
    }, [filter, tick])

    const handleSearchInput = value => {
      setFilter({ ...filter, loading: true, page: 1, q: value })
    }
    const handlePageChange = (e, page) => {
      setFilter({ ...filter, page })
    }
    const handleReloadToggle = () => {
      setTick(!tick)
      onReload()
    }

    return (
      <>
        <Component
          {...props}
          datatable={data}
          loading={data.loading}
          error={data.error}
          onSearchInput={handleSearchInput}
          onReload={handleReloadToggle}
        />
        <TablePagination
          style={{ textAlign: 'right' }}
          count={data.total_count || 0}
          page={filter.page}
          rowsPerPage={filter.limit}
          onPageChange={handlePageChange}
        />
      </>
    )
  }
  wrapped.prototype = {
    filter: PropTypes.object,
    onReload: PropTypes.func,
  }
  wrapped.defaultProps = {
    filter: {},
    onReload: () => {},
  }

  return wrapped
}

withDataFetch.propTypes = {
  Component: PropTypes.element,
}

export default withDataFetch
