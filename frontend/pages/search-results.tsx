import { useRouter } from "next/router";
import { useEffect, useState } from "react";

export default function SearchResults() {
  const router = useRouter();
  const [cards, setCards] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!router.query.cards) {
      setLoading(false);
      return;
    }

    try {
      const parsedCards = JSON.parse(router.query.cards as string);
      setCards(parsedCards);
    } catch (error) {
      console.error("Error parsing card data:", error);
    } finally {
      setLoading(false);
    }
  }, [router.query.cards]);

  if (loading) return <p>Loading...</p>;

  return (
    <div>
      <h1>Search Results</h1>
      {cards.length === 0 ? (
        <p>No cards found.</p>
      ) : (
        <ul style={{ listStyle: "none", padding: 0 }}>
          {cards.map((card) => (
            <li key={card.id} style={{ marginBottom: "15px", cursor: "pointer" }} onClick={() => router.push(`/card/${encodeURIComponent(card.name)}`)}>
              <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
                {card.image_uris?.small && <img src={card.image_uris.small} alt={card.name} width="50" />}
                <span>{card.name}</span>
              </div>
            </li>
          ))}
        </ul>
      )}
      <button onClick={() => router.push("/")}>Back to Home</button>
    </div>
  );
}