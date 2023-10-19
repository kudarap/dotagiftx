module.exports = {
  // NextJS v12 Rust compiler for minification
  swcMinify: true,
  async rewrites() {
    return [
      {
        source: '/id/:slug',
        destination: '/profiles/:slug?vanity=:slug',
      },
    ]
  },
  async redirects() {
    return [
      {
        source: '/item/:slug',
        destination: '/:slug',
        permanent: true,
      },
      {
        source: '/user/:id',
        destination: '/profiles/:id',
        permanent: true,
      },
      {
        source: '/blacklist',
        destination: '/bans',
        permanent: true,
      },
      {
        source: '/banned-users',
        destination: '/bans',
        permanent: true,
      },
    ]
  },
  images: {
    domains: ['localhost', 'api.dotagiftx.com', 'd2gapi.chiligarlic.com'],
    minimumCacheTTL: 31556952, // 1 year
    deviceSizes: [640, 750, 828, 1080, 1200, 1920, 2048, 3840],
  },
}
