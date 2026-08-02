import React from 'react'
import { APP_URL } from '@/constants/strings'
import { MARKET_STATUS_LIVE, MARKET_TYPE_ASK } from '@/constants/market'
import { VERIFIED_INVENTORY_VERIFIED } from '@/constants/verified'
import { catalog as getCatalog } from '@/service/api'
import ItemDetails from '@/components/ItemDetails'
import type { ItemDetailsProps } from '@/components/ItemDetails'
import ErrorPage from './404'
import type { Item } from '@/lib/types'
import type { MarketDatatable } from '@/components/MarketList'
import type { QueryFilter } from '@/service/http'

export default function DynamicPage(props: Partial<ItemDetailsProps>) {
  const { error } = props
  if (error) {
    return <ErrorPage />
  }

  const { item } = props
  if (!item) {
    return null
  }

  return <ItemDetails {...(props as ItemDetailsProps)} />
}

const marketSearchFilter = {
  page: 1,
  type: MARKET_TYPE_ASK,
  status: MARKET_STATUS_LIVE,
  inventory_status: VERIFIED_INVENTORY_VERIFIED,
  sort: 'lowest',
}

interface CatalogItem extends Item {
  quantity?: number
  bid_count?: number
  asks?: MarketDatatable['data']
  id: string | number
}

// This gets called on every request
export async function getServerSideProps(props: {
  params: { slugs: string[] }
  query: { page?: string }
}) {
  const { params, query } = props
  const { slugs } = params

  // NOTE: this is weird routing bug. maybe happening during page transition.
  if (slugs.indexOf('undefined') !== -1) {
    return {
      props: {},
    }
  }

  const [itemSlug, marketTypeParam, sortParam] = slugs

  // Hotfix backward compatible spelling support of ES arcana
  if (itemSlug === 'intergalactic-orbliterator-earthshaker') {
    return {
      redirect: {
        destination: 'intergalactic-obliterator-earthshaker',
        permanent: false,
      },
    }
  }

  let catalog = {} as CatalogItem
  let error: string | null = null

  const sort = sortParam || marketSearchFilter.sort
  const page = Number(query.page || marketSearchFilter.page)
  const filter: QueryFilter & { item_id?: string | number } = {
    ...marketSearchFilter,
    sort,
    page,
  }

  try {
    catalog = (await getCatalog(itemSlug, filter)) as CatalogItem
    filter.item_id = catalog.id
  } catch (e) {
    error = `catalog get error: ${(e as Error).message}`
  }

  if (!catalog.id) {
    return {
      notFound: true,
    }
  }

  const askData = catalog.asks || []
  const initialAsks: MarketDatatable = {
    data: askData,
    total_count: catalog.quantity || 0,
  }
  const initialBids: MarketDatatable = {
    data: [],
    total_count: catalog.bid_count || 0,
  }

  const canonicalURL = `${APP_URL}/${itemSlug}`
  const marketType = marketTypeParam || 'offers'

  // reduce large page data by nullifying un-needed data
  // nextjs.org/docs/messages/large-page-data
  // this can be reduce on API server level
  delete catalog.asks
  for (let i = 0; i < initialAsks.data.length; i++) {
    const market = initialAsks.data[i] as Partial<typeof initialAsks.data[number]> & {
      inventory?: { steam_assets?: unknown }
    }
    delete market.item
    if (market.inventory) {
      delete market.inventory.steam_assets
    }
  }

  return {
    props: {
      item: catalog,
      canonicalURL,
      filter,
      marketType,
      sortParam: sortParam || 'price',
      initialAsks,
      initialBids,
      error,
    },
  }
}
