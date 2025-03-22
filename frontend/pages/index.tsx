import { useRouter } from "next/router";
import { useState } from "react";
import Link from "next/link";

export default function Home() {
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState("");

  const handleSearch = async (event: React.FormEvent) => {
    event.preventDefault();
    if (!searchQuery.trim()) return;

    // Encode query to properly handle special characters like apostrophes
    const encodedQuery = encodeURIComponent(searchQuery);

    const res = await fetch(`/api/card/${encodedQuery}`);
    const data = await res.json();

    if (data.exact_match) {
      // Redirect to card details page if an exact match is found
      router.push(`/card/${encodeURIComponent(data.exact_match.name)}`);
    } else if (data.fuzzy_matches && data.fuzzy_matches.length > 0) {
      // Redirect to search results if multiple matches exist
      router.push({
        pathname: "/search-results",
        query: { cards: JSON.stringify(data.fuzzy_matches) },
      });
    } else {
      alert("No cards found.");
    }
  };

  // Function to fetch a random card and redirect
  const getRandomCard = async () => {
    const res = await fetch(`/api/card/random`);
    const data = await res.json();

    if (data.exact_match) {
      router.push(`/card/${encodeURIComponent(data.exact_match.name)}`);
    } else {
      alert("No random card found.");
    }
  };

  return (
    <div>
      <h1>Mana Tomb</h1>
      <Link href="/register">Register</Link> | <Link href="/login">Login</Link>
      <form onSubmit={handleSearch}>
        <input
          type="text"
          placeholder="Search for a card..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
        <button type="submit">Search</button>
      </form>

      <button onClick={getRandomCard}>Get Random Card</button>

      {/* Buttons for other pages */}
      <button onClick={() => router.push("/deck-forum")}>Deck Forum</button>
      <button onClick={() => router.push("/rules-game")}>Rules/Game</button>
      <button onClick={() => router.push("/deck-building")}>Deck Building</button>
    </div>
  );
}