import React from 'react'
import { myMarketSearch } from '@/service/api'
import TablePagination from '@/components/TablePagination'
import type { CSSProperties } from 'react'
import type { QueryFilter } from '@/service/http'
import type { ComponentType } from 'react'

interface Datatable<T> {
  data: T[]
  result_count?: number
  total_count?: number
  loading: boolean
  error: string | null
}

const initialDatatable: Datatable<never> = {
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

type SearchFn = (filter: QueryFilter) => Promise<unknown>

export interface DatatableFetchProps {
  datatable: Datatable<never>
  loading: boolean
  error: string | null
  onSearchInput: (value: string) => void
  onReload: () => void
}

interface WrappedProps {
  filter?: QueryFilter
  onReload?: () => void
}

function withDataFetch<P>(
  Component: ComponentType<P>,
  initFilter: QueryFilter = {},
  searchFn: SearchFn = myMarketSearch
) {
  function Wrapped({
    filter: propFilter,
    onReload = () => {},
    ...props
  }: WrappedProps & Omit<P, keyof DatatableFetchProps>) {
    const [data, setData] = React.useState<Datatable<never>>(initialDatatable)
    const [filter, setFilter] = React.useState<QueryFilter & { page?: number; q?: string }>({
      ...datatableBaseFilter,
      ...initFilter,
      ...propFilter,
    })
    const [tick, setTick] = React.useState(false)

    React.useEffect(() => {
      ;(async () => {
        setData(current => ({ ...current, loading: true, error: null }))
        try {
          const res = (await searchFn(filter)) as Partial<Datatable<never>>
          setData(current => ({ ...current, loading: false, ...res }))
        } catch (e) {
          setData(current => ({ ...current, loading: false, error: (e as Error).message }))
        }
      })()
    }, [filter, tick])

    const handleSearchInput = (value: string) => {
      setFilter({ ...filter, loading: true, page: 1, q: value })
    }
    const handlePageChange = (e: unknown, page: number) => {
      setFilter({ ...filter, page })
    }
    const handleReloadToggle = () => {
      setTick(!tick)
      onReload()
    }

    return (
      <>
        <Component
          {...(props as P)}
          datatable={data}
          loading={data.loading}
          error={data.error}
          onSearchInput={handleSearchInput}
          onReload={handleReloadToggle}
        />
        <TablePagination
          style={{ textAlign: 'right' }}
          count={data.total_count || 0}
          page={filter.page || 1}
          rowsPerPage={Number(filter.limit) || 10}
          onPageChange={handlePageChange}
        />
      </>
    )
  }

  return Wrapped
}

export default withDataFetch
