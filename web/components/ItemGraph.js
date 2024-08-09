import React from 'react'
import useSWR from 'swr'
import Typography from '@mui/material/Typography'
import { MARKET_STATUS_RESERVED, MARKET_STATUS_SOLD, MARKET_TYPE_ASK } from '@/constants/market'
import { fetcher, GRAPH_MARKET_SALES, MARKETS } from '@/service/api'
import MarketSalesChart from '@/components/MarketSalesChart'
import MarketActivity from '@/components/MarketActivity'

const swrConfig = [
  fetcher,
  {
    revalidateOnFocus: false,
    revalidateOnMount: true,
  },
]

const marketSalesGraphFilter = {
  type: MARKET_TYPE_ASK,
}

const marketReservedFilter = {
  type: MARKET_TYPE_ASK,
  status: MARKET_STATUS_RESERVED,
  sort: 'updated_at:desc',
  index: 'item_id',
}

const marketDeliveredFilter = {
  type: MARKET_TYPE_ASK,
  status: MARKET_STATUS_SOLD,
  sort: 'updated_at:desc',
  index: 'item_id',
}

const ItemGraph = ({ itemId = '', itemName = '' }) => {
  // Retrieve market sales graph.
  const shouldLoadGraph = true
  marketSalesGraphFilter.item_id = itemId
  const { data: marketGraph, error: marketGraphError } = useSWR(
    shouldLoadGraph ? [GRAPH_MARKET_SALES, marketSalesGraphFilter] : null,
    ...swrConfig
  )

  // Retrieve market sale activity.
  const shouldLoadHistory = Boolean(marketGraph)
  marketReservedFilter.item_id = itemId
  const {
    data: marketReserved,
    error: marketReservedError,
    isValidating: marketReservedLoading,
  } = useSWR(shouldLoadHistory ? [MARKETS, marketReservedFilter] : null, ...swrConfig)

  marketDeliveredFilter.item_id = itemId
  const {
    data: marketDelivered,
    error: marketDeliveredError,
    isValidating: marketDeliveredLoading,
  } = useSWR(marketReserved ? [MARKETS, marketDeliveredFilter] : null, ...swrConfig)

  let historyCount = false
  if (!marketReservedError && marketReserved) {
    historyCount = marketReserved.total_count
  }
  if (!marketDeliveredError && marketDelivered) {
    historyCount += marketDelivered.total_count
  }

  const isHistoryLoading = marketReservedLoading || marketDeliveredLoading
  return (
    <>
      {shouldLoadHistory && isHistoryLoading && <div>Loading {itemName} history...</div>}
      {shouldLoadHistory && !isHistoryLoading && (
        <div>
          <div>{itemName} history</div>
          {historyCount === 0 && (
            <Typography variant="body2" color="textSecondary">
              No history yet
            </Typography>
          )}
          {!marketGraphError && marketGraph && (
            <>
              <br />
              <MarketSalesChart data={marketGraph} />
            </>
          )}

          <div id="reserved">
            {marketReserved && marketReserved.result_count !== 0 && (
              <MarketActivity datatable={marketReserved} disablePrice />
            )}
          </div>
          <div id="delivered">
            {marketDelivered && marketDelivered.result_count !== 0 && (
              <MarketActivity datatable={marketDelivered} disablePrice />
            )}
          </div>
        </div>
      )}
    </>
  )
}

export default ItemGraph
