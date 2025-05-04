// pages/create-deck.tsx

import { useState, useEffect } from "react";
import { useRouter } from "next/router";
import { useAuth } from "./authContext";

export default function CreateDeck() {
  const router = useRouter();
  const { user } = useAuth();
  const [deckName, setDeckName] = useState("");
  const [description, setDescription] = useState("");
  const [commander, setCommander] = useState("");
  const [cards, setCards] = useState<string[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [searchResults, setSearchResults] = useState<any[]>([]);

  useEffect(() => {
    if (!user) {
      router.push("/login");
    }
  }, [user, router]);

  const handleSearch = async () => {
    const res = await fetch(`/api/card/${encodeURIComponent(searchQuery)}`);
    const data = await res.json();
    if (data.exact_match) {
      setSearchResults([data.exact_match]);
    } else if (data.fuzzy_matches) {
      setSearchResults(data.fuzzy_matches);
    } else {
      setSearchResults([]);
    }
  };

  const handleAddCard = (cardName: string) => {
    if (!cards.includes(cardName)) {
      setCards([...cards, cardName]);
    }
  };

  const handleSaveDeck = async () => {
    const res = await fetch("/api/decks/create", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({
        name: deckName,
        description,
        commander,
        cards,
        user_id: user.id,
      }),
    });

    if (res.ok) {
      router.push("/deck-building");
    } else {
      alert("Failed to save deck.");
    }
  };

  return (
    <div>
      <h1>Create New Commander Deck</h1>

      <input
        type="text"
        placeholder="Deck Name"
        value={deckName}
        onChange={(e) => setDeckName(e.target.value)}
      />
      <br />
      <textarea
        placeholder="Deck Description"
        value={description}
        onChange={(e) => setDescription(e.target.value)}
      />
      <br />
      <input
        type="text"
        placeholder="Commander"
        value={commander}
        onChange={(e) => setCommander(e.target.value)}
      />
      <br />
      <h2>Search for Cards</h2>
      <input
        type="text"
        placeholder="Card name"
        value={searchQuery}
        onChange={(e) => setSearchQuery(e.target.value)}
      />
      <button onClick={handleSearch}>Search</button>

      <ul>
        {searchResults.map((card, i) => (
          <li key={i}>
            {card.name}
            <button onClick={() => handleAddCard(card.name)}>Add</button>
          </li>
        ))}
      </ul>

      <h3>Cards in Deck</h3>
      <ul>
        {cards.map((card, i) => (
          <li key={i}>{card}</li>
        ))}
      </ul>

      <button onClick={handleSaveDeck}>Save Deck</button>
      <button onClick={() => router.push("/deck-builder")}>Back to Decks</button>
      <button onClick={() => router.push("/")}>Home</button>
    </div>
  );
}
