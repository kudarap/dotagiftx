import React from 'react'
import PropTypes from 'prop-types'
import { MARKET_STATUS_LIVE } from '@/constants/market'
import { catalog, marketSearch } from '@/service/api'
import { APP_URL } from '@/constants/strings'
import ItemDetails from '@/components/ItemDetails'
import ErrorPage from './404'

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
  status: MARKET_STATUS_LIVE,
  sort: 'price',
}

// This gets called on every request
export async function getServerSideProps(props) {
  const { params, query } = props
  const { slug } = params

  let item = {}
  let error = null

  // Handles no market entry on item
  try {
    item = await catalog(slug)
  } catch (e) {
    error = `catalog get error: ${e.message}`
  }

  if (!item.id) {
    return {
      props: {
        item,
        filter: {},
        error: 'catalog not found',
      },
    }
  }

  const filter = { ...marketSearchFilter, item_id: item.id }
  if (query.page) {
    filter.page = Number(query.page)
  }

  const markets = {}
  // try {
  //   markets = await marketSearch(filter)
  // } catch (e) {
  //   console.log(`market search error: ${e.message}`)
  //   error = e.message
  // }

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
