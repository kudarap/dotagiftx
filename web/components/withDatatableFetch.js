import React from 'react'
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
  sort: 'created_at:desc',
  page: 1,
}

const withDataFetch = (Component, initFilter) => props => {
  const [data, setData] = React.useState(initialDatatable)
  const [filter, setFilter] = React.useState({ ...datatableBaseFilter, ...initFilter })
  const [tick, setTick] = React.useState(false)

  React.useEffect(() => {
    ;(async () => {
      setData({ ...data, loading: true, error: null })
      try {
        const res = await myMarketSearch(filter)
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
  }

  return (
    <>
      <Component
        datatable={data}
        loading={data.loading}
        error={data.error}
        onSearchInput={handleSearchInput}
        onReload={handleReloadToggle}
        {...props}
      />
      <TablePagination
        style={{ textAlign: 'right' }}
        count={data.total_count || 0}
        page={filter.page}
        rowsPerPage={filter.limit}
        onChangePage={handlePageChange}
      />
    </>
  )
}

export default withDataFetch
