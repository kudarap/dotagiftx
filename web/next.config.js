module.exports = {
  // NextJS v12 Rust compiler for minification
  // swcMinify: true,
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
        destination: '/banned-users',
        permanent: true,
      },
    ]
  },
}
