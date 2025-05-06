import { useEffect, useState } from "react";
import { useRouter } from "next/router";
import { useAuth } from "./authContext";

interface Deck {
  deck_id: number;
  user_id: number;
  name: string;
  description: string;
  created_at: string;
  commander_id: string | null;
}

export default function DeckBuildingPage() {
  const router = useRouter();
  const { user, loading: authLoading } = useAuth();
  const [decks, setDecks] = useState<Deck[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!authLoading && !user) {
      router.push("/login");
    }
  }, [authLoading, user]);

  useEffect(() => {
    const fetchDecks = async () => {
      try {
        const res = await fetch("/api/decks", {
          method: "GET",
          credentials: "include",
        });
        if (!res.ok) throw new Error("Failed to fetch decks");
        const data = await res.json();
        setDecks(data);
      } catch (err: any) {
        setError(err.message || "An error occurred");
      } finally {
        setLoading(false);
      }
    };

    fetchDecks();
  }, []);

  const handleCreateDeck = () => {
    router.push("/create-deck");
  };

  const handleHome = () => {
    router.push("/");
  };

  if (authLoading || !user) {
    return <p>Loading...</p>;
  }

  const getDeckDetails = (deckId: number) => {
    const encodedDeckId = encodeURIComponent(deckId.toString());
    router.push(`/decks/${encodedDeckId}`);
  }

  return (
    <div style={{ padding: "2rem" }}>
      <h1>Your Commander Decks</h1>
      <button onClick={handleCreateDeck}>Create New Deck</button>
      <button onClick={handleHome} style={{ marginLeft: "10px" }}>
        Back to Home
      </button>

      {loading && <p>Loading decks...</p>}
      {error && <p style={{ color: "red" }}>{error}</p>}

      {Array.isArray(decks) && decks.length > 0 ? (
      <ul>
      {decks.map((deck) => (
        <li
          key={deck.deck_id}
          style={{ marginBottom: "1rem", cursor: "pointer" }}
          onClick={() => getDeckDetails(deck.deck_id)} // Navigate to deckName page
        >
          <h3>{deck.name}</h3>
          <p>{deck.description}</p>
          <p>Deck ID: {deck.deck_id}</p>
          <small>Created at: {new Date(deck.created_at).toLocaleString()}</small>
        </li>
      ))}
    </ul>
    ) : (
      !loading && <p>No decks found.</p>
    )}
    </div>
  );
}