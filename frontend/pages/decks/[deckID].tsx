import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import { useAuth } from "../authContext";

export default function DeckDetailsPage() {
  const router = useRouter();
  const { deckID } = router.query;
  const { user, loading: authLoading } = useAuth();
  const [deck, setDeck] = useState<any | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!authLoading && !user) {
      router.push("/login");
    }
  }, [authLoading, user]);

  useEffect(() => {
    if (!deckID) return;

    const fetchDeck = async () => {
      try {
        const safedeckID = encodeURIComponent(deckID as string);
        const endpoint = `/api/decks/${safedeckID}`;
        const res = await fetch(endpoint, {
          method: "GET",
          credentials: "include",
        });
        if (!res.ok) {
          throw new Error("Deck not found");
        }
        const data = await res.json();

        setDeck(data);
      } catch (err: any) {
        setError(err.message);
      }
    };

    fetchDeck();
  }, [deckID]);

  const handleEdit = () => {
    if (deckID) {
      router.push(`/decks/edit-deck/${deckID}`);
    }
  };

  const handleDelete = async () => {
    if (!deckID) return;
    const confirmed = window.confirm("Are you sure you want to delete this deck?");
    if (!confirmed) return;

    try {
      const res = await fetch(`/api/decks/delete/${deckID}`, {
        method: "DELETE",
        credentials: "include",
      });

      if (!res.ok) {
        throw new Error("Failed to delete deck");
      }

      alert("Deck deleted successfully");
      router.push("/deck-building");
    } catch (err) {
      console.error("Error deleting deck:", err);
      alert("An error occurred while deleting the deck.");
    }
  };

  if (authLoading || !user) {
    return <p>Loading...</p>;
  }
  if (error) return <p>Error: {error}</p>;
  if (!deck) return <p>No card data found.</p>;

  return (
    <div>
      <h1>{deck.name}</h1>
      <p><strong>Description:</strong> {deck.description || "N/A"}</p>
      <p><strong>Commander:</strong> {deck.commander || "N/A"}</p>
      <p><strong>Cards:</strong></p>
      {deck.cards.length === 0 ? (
        <p>No cards found.</p>
      ) : (
        <ul style={{ listStyle: "none", padding: 0 }}>
          {deck.cards.map((cardName: string, index: number) => (
            <li key={`${cardName}-${index}`} style={{ marginBottom: "15px", cursor: "pointer" }} onClick={() => router.push(`/card/${encodeURIComponent(cardName)}`)}>
              <span>{cardName}</span>
            </li>
          ))}
        </ul>
      )}
      <button onClick={() => router.push("/deck-building")}>Back to Deck Building</button>
      <button onClick={() => router.push("/")}>Back to Home</button>
      <button onClick={handleEdit}>Edit Deck</button>
      <button onClick={handleDelete}>Delete Deck</button>
    </div>
  );
}