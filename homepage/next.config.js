const isProduction = process.env.NODE_ENV == "production";
const isDevelopment = process.env.NODE_ENV == "development";

const prodConfig = {
  output: "export",
};

const devConfig = {
  async rewrites() {
    return [
      {
        source: "/api/poc/:path*",
        destination:
          "https://waf-ce.chaitin.cn/api/poc/:path*",
      },
    ];
  },
};

/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: false,
  images: { unoptimized: true },
};

Object.assign(
  nextConfig,
  isProduction && prodConfig,
  isDevelopment && devConfig
);

module.exports = nextConfig;
