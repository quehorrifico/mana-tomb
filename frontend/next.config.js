/** @type {import('next').NextConfig} */
const nextConfig = {
    reactStrictMode: true,
    async rewrites() {
      return [
        {
          source: "/api/card/:cardName",
          destination: "http://localhost:8080/card/:cardName", // Updated to match backend route
        },
        {
          source: "/api/me",
          destination: "http://localhost:8080/me", // Updated to match backend route
        },
        {
          source: "/api/login",
          destination: "http://localhost:8080/login", // Updated to match backend route
        },
        {
          source: "/api/logout",
          destination: "http://localhost:8080/logout", // Updated to match backend route
        },
        {
          source: "/api/register",
          destination: "http://localhost:8080/register", // Updated to match backend route
        },
        {
          source: "/api/decks",
          destination: "http://localhost:8080/decks", // Updated to match backend route
        },
        {
          source: "/api/decks/create",
          destination: "http://localhost:8080/decks/create", // Updated to match backend route
        },
        {
          source: "/api/decks/:deckID",
          destination: "http://localhost:8080/decks/:deckID",
        },
      ];
    },
};

module.exports = nextConfig;