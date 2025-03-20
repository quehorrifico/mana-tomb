import React, { useState } from "react";
import { useRouter } from "next/router";

export default function Home() {
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState("");

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (!searchQuery.trim()) return;

    // Redirect to the new card page with the searched card name
    router.push(`/card/${encodeURIComponent(searchQuery.trim())}`);
  };

  return (
    <div>
      <h1>Mana Tomb Home</h1>

      {/* Search Bar */}
      <form onSubmit={handleSearch}>
        <input
          type="text"
          placeholder="Search for a card..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
        <button type="submit">Search</button>
      </form>

      {/* Navigation Buttons */}
      <div style={{ marginTop: "20px" }}>
        <button onClick={() => router.push("/deck-forum")}>Deck Forum</button>
        <button onClick={() => router.push("/rules")}>Rules/Game</button>
        <button onClick={() => router.push("/deck-building")}>Deck Building</button>
      </div>
    </div>
  );
}