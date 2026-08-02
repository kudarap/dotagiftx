// Domain types mirroring the dotagiftx.com API response shapes.

export interface InventoryAsset {
  asset_id: string
  name: string
  qty: number
  type?: string
  gift_once?: boolean
  date_received?: string
  contains?: string
  gift_from?: string
}

export interface VerificationData {
  status: number
  verified_by?: string
  elapsed_ms?: number
  updated_at: string
  created_at?: string
  bundle_count?: number
  steam_assets?: InventoryAsset[]
}

export interface UserSummary {
  id: number | string
  steam_id: string
  name: string
  avatar: string
  boons?: string[]
  donated_at?: string
  donation?: string
  created_at?: string
}

export interface Market {
  id: string
  type?: number
  status: number
  inventory_status?: number
  delivery_status?: number
  price?: number
  currency?: string
  notes?: string
  qty?: number
  resell?: boolean
  created_at?: string
  updated_at?: string
  user: UserSummary
  partner_steam_id?: string
  seller_steam_id?: string
  item?: Item
  inventory?: VerificationData
  delivery?: VerificationData
}

export interface SearchResponse<T> {
  data: T[]
  total_count: number
  result_count?: number
  count?: number
  limit?: number
}

export interface MarketSummary {
  live: string | number
  reserved: string | number
  sold: string | number
  bids: {
    bid_live: string | number
  }
}

export interface Item {
  id: string | number
  slug?: string
  name: string
  hero?: string
  image?: string
  category?: string
  rarity?: string
  quality?: string
  origin?: string
  lowest_ask?: number
  highest_bid?: number
  reserved_count?: number
  sold_count?: number
}

export interface Profile extends UserSummary {
  subscription?: number
  status?: number
  notes?: string
  is_registered?: boolean
  market_stats?: {
    live: number
    reserved: number
    sold: number
    bid_completed: number
  }
  stats?: {
    live: number
    reserved: number
    sold: number
    bid_completed: number
  }
  created_at?: string
  updated_at?: string
  donated_at?: string
  donation?: string
}

export interface Treasure {
  id: string | number
  slug: string
  name: string
  image?: string
  price?: number
  lowest_ask?: number
  boons?: string[]
  avg_days_to_sell?: number
}

export interface Hero {
  id: string | number
  name: string
  image?: string
  markets?: SearchResponse<Market>
}

export interface StatsSummary {
  [key: string]: number
}
