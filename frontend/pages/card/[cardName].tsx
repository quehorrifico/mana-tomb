import { useRouter } from "next/router";
import { useEffect, useState } from "react";

export default function CardPage() {
  const router = useRouter();
  const { cardName } = router.query;
  const [cardData, setCardData] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!cardName) return;

    const fetchCard = async () => {
      try {
        const response = await fetch(`http://localhost:8080/card/${encodeURIComponent(cardName as string)}`);
        if (!response.ok) throw new Error("Card not found");

        const data = await response.json();
        setCardData(data);
        setError(null);
      } catch (err) {
        setError("Card not found");
        setCardData(null);
      }
    };

    fetchCard();
  }, [cardName]);

  return (
    <div>
      <h1>Card Details</h1>
      {error && <p style={{ color: "red" }}>{error}</p>}
      {cardData && (
        <div style={{ border: "1px solid #ccc", padding: "10px", marginTop: "20px" }}>
          <h2>{cardData.name}</h2>
          <p><strong>Mana Cost:</strong> {cardData.mana_cost || "N/A"}</p>
          <p><strong>Oracle Text:</strong> {cardData.oracle_text}</p>
          {cardData.image_uris?.normal && <img src={cardData.image_uris.normal} alt={cardData.name} />}
        </div>
      )}
      <button onClick={() => router.push("/")}>Back to Home</button>
    </div>
  );
}