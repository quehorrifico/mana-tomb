import { useRouter } from "next/router";
import { useEffect, useState } from "react";

interface ImageURIs {
  normal: string;
}

interface Card {
  name: string;
  mana_cost: string;
  type_line: string;
  oracle_text: string;
  image_uris: ImageURIs;
}

export default function CardResults() {
  const router = useRouter();
  const { cardName } = router.query;
  const [cards, setCards] = useState<Card[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!cardName) return;

    fetch(`/api/card/${cardName}`)
      .then((res) => {
        if (!res.ok) throw new Error("No matching cards found");
        return res.json();
      })
      .then((data) => {
        setCards(data);
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message);
        setLoading(false);
      });
  }, [cardName]);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>{error}</p>;

  return (
    <div style={{ textAlign: "center", padding: "20px" }}>
      <h1>Search Results for: {cardName}</h1>
      {cards.length === 0 ? (
        <p>No matching cards found.</p>
      ) : (
        <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(200px, 1fr))", gap: "20px" }}>
          {cards.map((card) => (
            <div key={card.name} style={{ border: "1px solid #ccc", padding: "10px", borderRadius: "8px" }}>
              <h2>{card.name}</h2>
              <p><strong>Mana Cost:</strong> {card.mana_cost || "N/A"}</p>
              <p><strong>Type:</strong> {card.type_line}</p>
              <p><strong>Oracle Text:</strong> {card.oracle_text}</p>
              {card.image_uris?.normal && (
                <img src={card.image_uris.normal} alt={card.name} style={{ width: "100%" }} />
              )}
            </div>
          ))}
        </div>
      )}
      <button onClick={() => router.push("/")}>Back to Home</button>
    </div>
  );
}