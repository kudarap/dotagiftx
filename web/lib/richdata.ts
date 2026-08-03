import { CDN_URL } from '@/service/api'

function schemaOrgProduct(
  canonicalURL: string,
  item: { id?: string | number; name?: string; image?: string; lowest_ask?: number } = {},
  other: object
) {
  const data: {
    '@context': string
    '@type': string
    productID?: string | number
    name?: string
    image?: string
    offers: {
      '@type': string
      priceCurrency: string
      url: string
      availability?: string
      price?: string
    }
  } = {
    '@context': 'https://schema.org',
    '@type': 'Product',
    productID: item.id,
    name: item.name,
    image: `${CDN_URL}/${item.image}`,
    offers: {
      '@type': 'Offer',
      priceCurrency: 'USD',
      url: canonicalURL,
    },
    ...other,
  }

  if (item.lowest_ask) {
    data.offers.availability = 'https://schema.org/InStock'
    data.offers.price = item.lowest_ask.toFixed(2)
  } else {
    data.offers.availability = 'https://schema.org/OutOfStock'
    data.offers.price = '0'
  }

  return data
}

export default schemaOrgProduct
