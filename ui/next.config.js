/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // We'll use a local API
  rewrites: async () => {
    return [
      {
        source: '/api/:path*',
        destination: 'http://localhost:8080/api/:path*',
      },
    ]
  },
}

module.exports = nextConfig
