import { CDN_URL } from '@/service/api'

export function schemaOrgProduct(canonicalURL, item = {}, other) {
  const data = {
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
