import React from 'react'
import PropTypes from 'prop-types'
import { MARKET_STATUS_LIVE, MARKET_TYPE_ASK } from '@/constants/market'
import { catalog as getCatalog } from '@/service/api'
import { APP_URL } from '@/constants/strings'
import ItemDetails from '@/components/ItemDetails'
import { VERIFIED_INVENTORY_VERIFIED } from '@/constants/verified'
import ErrorPage from './404'

export default function DynamicPage(props) {
  const { error } = props
  if (error) {
    return <ErrorPage />
  }

  const { item } = props
  if (!item) {
    return null
  }

  return <ItemDetails {...props} />
}
DynamicPage.propTypes = {
  item: PropTypes.object,
  error: PropTypes.string,
}
DynamicPage.defaultProps = {
  error: null,
  item: null,
}

const marketSearchFilter = {
  page: 1,
  type: MARKET_TYPE_ASK,
  status: MARKET_STATUS_LIVE,
  inventory_status: VERIFIED_INVENTORY_VERIFIED,
  sort: 'best',
}

// This gets called on every request
export async function getServerSideProps(props) {
  const { params, query } = props
  const { slug } = params
  console.log('SSR', slug, query)
  // NOTE: this is weird routing bug. maybe happening during page transition.
  if (slug === 'undefined') {
    return {
      props: {},
    }
  }

  let catalog = {}
  let error = null

  const sort = query.sort || marketSearchFilter.sort
  const page = Number(query.page || marketSearchFilter.page)
  const filter = { ...marketSearchFilter, sort, page }

  try {
    catalog = await getCatalog(slug, filter)
    filter.item_id = catalog.id
  } catch (e) {
    error = `catalog get error: ${e.message}`
  }

  if (!catalog.id) {
    return {
      props: {
        item: catalog,
        filter: {},
        error: 'catalog not found x',
      },
    }
  }

  const askData = catalog.asks || []
  const initialAsks = {
    data: askData,
    result_count: askData.length || 0,
    total_count: catalog.quantity,
  }

  const initialBids = {
    data: [],
    result_count: 0,
    total_count: catalog.bid_count,
  }

  const canonicalURL = `${APP_URL}/${slug}`

  return {
    props: {
      item: catalog,
      canonicalURL,
      filter,
      initialAsks,
      initialBids,
      error,
    },
  }
}
