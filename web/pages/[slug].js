import React from 'react'
import PropTypes from 'prop-types'
import { MARKET_STATUS_LIVE, MARKET_TYPE_ASK } from '@/constants/market'
import { catalog as getCatalog } from '@/service/api'
import { APP_URL } from '@/constants/strings'
import ItemDetails from '@/components/ItemDetails'
import ErrorPage from './404'
import { VERIFIED_INVENTORY_VERIFIED } from '@/constants/verified'

export default function DynamicPage(props) {
  const { error } = props
  if (error) {
    console.error(error)
    return <ErrorPage />
  }

  return <ItemDetails {...props} />
}
DynamicPage.propTypes = {
  error: PropTypes.string,
}
DynamicPage.defaultProps = {
  error: null,
}

const marketSearchFilter = {
  page: 1,
  type: MARKET_TYPE_ASK,
  status: MARKET_STATUS_LIVE,
  inventory_status: VERIFIED_INVENTORY_VERIFIED,
  sort: 'price',
}

// This gets called on every request
export async function getServerSideProps(props) {
  const { params, query } = props
  const { slug } = params

  let catalog = {}
  let error = null

  // Handles no market entry on item
  try {
    catalog = await getCatalog(slug)
  } catch (e) {
    error = `catalog get error: ${e.message}`
  }

  if (!catalog.id) {
    return {
      props: {
        item: catalog,
        filter: {},
        error: 'catalog not found',
      },
    }
  }

  const filter = { ...marketSearchFilter, item_id: catalog.id }
  if (query.page) {
    filter.page = Number(query.page)
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
