module.exports = {
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
    ]
  },
}
