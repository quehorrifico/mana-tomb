// pages/index.tsx
import { useRouter } from "next/router";
import { useState } from "react";
import Link from "next/link";
import { useAuth } from "./authContext";

export default function Home() {
  const router = useRouter();
  const { user, logout, loading } = useAuth();
  const [searchQuery, setSearchQuery] = useState("");

  const handleSearch = async (event: React.FormEvent) => {
    event.preventDefault();
    if (!searchQuery.trim()) return;

    const encodedQuery = encodeURIComponent(searchQuery);
    const res = await fetch(`/api/card/${encodedQuery}`, {
      method: "GET",
      credentials: "include",
    });
    const data = await res.json();

    if (data.exact_match) {
      router.push(`/card/${encodeURIComponent(data.exact_match.name)}`);
    } else if (data.fuzzy_matches && data.fuzzy_matches.length > 0) {
      router.push({
        pathname: "/search-results",
        query: { cards: JSON.stringify(data.fuzzy_matches) },
      });
    } else {
      alert("No cards found.");
    }
  };

  const getRandomCard = async () => {
    const res = await fetch(`/api/card/random`, {
      method: "GET",
      credentials: "include",
    });
    const data = await res.json();
    if (data.exact_match) {
      router.push({
        pathname: "/card/[cardName]",
        query: { cardName: data.exact_match.name },
      });
    } else {
      alert("No random card found.");
    }
  };

  if (loading) {
    return <div>Loading user data...</div>;
  }

  return (
    <div style={{ display: "flex", justifyContent: "center", alignItems: "center", minHeight: "80vh", flexDirection: "column" }}>
      <h1>Mana Tomb</h1>
      <nav>
        {user ? (
          <>
            <span>Welcome, {user.username}!</span> |{" "}
            <button onClick={logout} style={{ marginLeft: "0.5rem" }}>Logout</button>
          </>
        ) : (
          <>
            <Link href="/register">Register</Link> |{" "}
            <Link href="/login">Login</Link>
          </>
        )}
      </nav>
      <form onSubmit={handleSearch}>
        <input
          type="text"
          placeholder="Search for a card..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          style={{ marginBottom: "1rem" }}
        />
        <button type="submit">Search</button>
      </form>
      <button onClick={getRandomCard}>Get Random Card</button>
      <div style={{ marginTop: "2rem", display: "flex", flexDirection: "column", gap: "0.5rem" }}>
        <button onClick={() => router.push("/deck-forum")}>Deck Forum</button>
        <button onClick={() => router.push("/rules-game")}>Rules/Game</button>
        <button onClick={() => router.push("/deck-building")}>Deck Building</button>
      </div>
    </div>
  );
}