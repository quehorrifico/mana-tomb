import { useRouter } from "next/router";
import { useEffect, useState } from "react";

export default function DeckDetailsPage() {
  const router = useRouter();
  const { deckID } = router.query;
  const [deck, setDeck] = useState<any | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

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
      } finally {
        setLoading(false);
      }
    };

    console.log("new Fetched deck:", deck); // Debugging line
    fetchDeck();
  }, [deckID]);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error: {error}</p>;
  if (!deck) return <p>No card data found.</p>;

  return (
    <div>
      <h1>{deck.name}</h1>
      <p><strong>Description:</strong> {deck.description || "N/A"}</p>
      <p><strong>Commander:</strong> {deck.commander_id || "N/A"}</p>
      <p><strong>Cards:</strong></p>
      <ul style={{ listStyle: "none", padding: 0 }}>
        {deck.cards.map((card: any) => (
          <li key={card.id} style={{ marginBottom: "15px", cursor: "pointer" }} onClick={() => router.push(`/card/${encodeURIComponent(card.name)}`)}>
            <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
              {card.image_uris?.small && <img src={card.image_uris.small} alt={card.name} width="50" />}
              <span>{card.name}</span>
            </div>
          </li>
        ))}
      </ul>
      <button onClick={() => router.push("/deck-building")}>Back to Deck Building</button>
      <button onClick={() => router.push("/")}>Back to Home</button>
    </div>
  );
}