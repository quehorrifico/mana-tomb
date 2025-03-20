/** @type {import('next').NextConfig} */
const nextConfig = {
    reactStrictMode: true,
    async rewrites() {
      return [
        {
          source: "/api/card/:cardName",
          destination: "http://localhost:8080/card/:cardName", // Adjust this if your backend URL is different
        },
      ];
    },
  };
  
  module.exports = nextConfig;