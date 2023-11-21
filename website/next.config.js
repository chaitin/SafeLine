/** @type {import('next').NextConfig} */
const nextConfig = {
  async rewrites() {
    return {
      fallback: [
        // These rewrites are checked after both pages/public files
        // and dynamic routes are checked
        {
          source: '/api/count',
          destination: 'https://waf-ce.chaitin.cn/api/count',
        },
        {
          source: '/api/:path*',
          destination: 'http://10.10.4.142:8080/api/:path*',
        },
      ],
    }
  },
}

module.exports = nextConfig
