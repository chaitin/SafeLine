/** @type {import('next').NextConfig} */
const nextConfig = {
  async rewrites() {
    return {
      fallback: [
        // These rewrites are checked after both pages/public files
        // and dynamic routes are checked
        {
          source: '/api/safeline/count',
          destination: 'http://121.199.46.182/api/safeline/count',
        },
        {
          source: '/api/:path*',
          destination: 'http://121.199.46.182/api/:path*',
        },
      ],
    }
  },
}

module.exports = nextConfig
