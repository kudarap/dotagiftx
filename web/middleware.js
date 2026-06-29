export function middleware(request) {
  const country = request.geo?.country || 'US'
  const region = request.geo?.region
  const city = request.geo?.city

  console.log(`User from: ${city}, ${region}, ${country}`)
  console.log(request.url)
}

// export const config = {
//   matcher: '/((?!api|_next/static|_next/image|favicon.ico).*)',
// }
