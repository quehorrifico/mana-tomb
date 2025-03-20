/** @type {import('next').NextConfig} */
const nextConfig = {
    reactStrictMode: true,
    async rewrites() {
      return [
        {
          source: "/api/card/:cardName",
          destination: "http://localhost:8080/card/:cardName", // Ensure this matches backend routes
        },
      ];
    },
};

module.exports = nextConfig;