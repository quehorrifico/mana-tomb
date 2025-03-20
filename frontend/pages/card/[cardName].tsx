import { useRouter } from "next/router";
import { useEffect, useState } from "react";

export default function CardPage() {
  const router = useRouter();
  const { cardName } = router.query;
  const [card, setCard] = useState<any | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!cardName) return;

    const fetchCard = async () => {
      try {
        const endpoint = cardName === "random" ? "/api/card/random" : `/api/card/${cardName}`;
        const res = await fetch(endpoint);
        if (!res.ok) {
          throw new Error("Card not found");
        }
        const data = await res.json();

        if (data.exact_match) {
          setCard(data.exact_match);
        } else if (data.fuzzy_matches?.length === 1) {
          setCard(data.fuzzy_matches[0]); // Pick the single result
        }
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchCard();
  }, [cardName]);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error: {error}</p>;
  if (!card) return <p>No card data found.</p>;

  return (
    <div>
      <h1>{card.name}</h1>
      <p><strong>Mana Cost:</strong> {card.mana_cost || "N/A"}</p>
      <p><strong>Type:</strong> {card.type_line || "N/A"}</p>
      <p><strong>Oracle Text:</strong> {card.oracle_text || "N/A"}</p>
      {card.image_uris?.normal && <img src={card.image_uris.normal} alt={card.name} />}
      <button onClick={() => router.push("/")}>Back to Home</button>
    </div>
  );
}