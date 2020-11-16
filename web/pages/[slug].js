import React from 'react'
import PropTypes from 'prop-types'
import { MARKET_STATUS_LIVE } from '@/constants/market'
import { catalog, marketSearch } from '@/service/api'
import { APP_URL } from '@/constants/strings'
import ItemPage from './item/[slug]'
import ErrorPage from './404'

export default function DynamicPage(props) {
  const { error } = props
  if (error) {
    return <ErrorPage />
  }

  return <ItemPage {...props} />
}
DynamicPage.propTypes = {
  error: PropTypes.string,
}
DynamicPage.defaultProps = {
  error: null,
}

const marketSearchFilter = {
  page: 1,
  status: MARKET_STATUS_LIVE,
  sort: 'price',
}

// This gets called on every request
export async function getServerSideProps(props) {
  const { params, query } = props
  const { slug } = params

  let item = {}

  // Handles no market entry on item
  try {
    item = await catalog(slug)
  } catch (e) {
    console.log(`catalog get error: ${e.message}`)
  }
  if (!item.id) {
    return {
      props: {
        item,
        filter: {},
      },
    }
  }

  const filter = { ...marketSearchFilter, item_id: item.id }
  if (query.page) {
    filter.page = Number(query.page)
  }

  let markets = {}
  let error = null
  try {
    markets = await marketSearch(filter)
  } catch (e) {
    console.log(`market search error: ${e.message}`)
    error = e.message
  }

  const canonicalURL = `${APP_URL}/${slug}`

  return {
    props: {
      item,
      canonicalURL,
      filter,
      markets,
      error,
    },
  }
}
