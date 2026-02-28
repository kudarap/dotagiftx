// @ts-check

const config = (phase, { defaultConfig }) => {
  /**
   * @type {import('next').NextConfig}
   */
  const nextConfig = {
    /* config options here */
    reactStrictMode: true,

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
      remotePatterns: [
        {
          protocol: 'https',
          hostname: '*.dotagiftx.com', // Matches cdn-1, cdn-2, etc.
          pathname: '/images/**',
        },
      ],
      minimumCacheTTL: 31556952, // 1 year
      deviceSizes: [640, 750, 828, 1080, 1200, 1920, 2048, 3840],
    },
  }
  return nextConfig
}

export default config
