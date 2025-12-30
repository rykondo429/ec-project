/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: 'standalone',
  env: {
    NEXT_PUBLIC_API_BASE_URL: process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8001',
    NEXT_PUBLIC_CART_API_URL: process.env.NEXT_PUBLIC_CART_API_URL || 'http://localhost:8002',
    NEXT_PUBLIC_ORDER_API_URL: process.env.NEXT_PUBLIC_ORDER_API_URL || 'http://localhost:8003',
    NEXT_PUBLIC_POINT_API_URL: process.env.NEXT_PUBLIC_POINT_API_URL || 'http://localhost:8004',
    NEXT_PUBLIC_PROMOTION_API_URL: process.env.NEXT_PUBLIC_PROMOTION_API_URL || 'http://localhost:8005',
  },
};

module.exports = nextConfig;
