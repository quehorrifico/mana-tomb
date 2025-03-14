import React, { useEffect, useState } from "react";

interface Card {
  name: string;
  mana_cost: string;
  oracle_text: string;
}

export default function RandomCard() {
  const [card, setCard] = useState<Card | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    // fetch("/random-card")
    fetch("http://localhost:8080/random-card")
      .then((res) => res.json())
      .then((data) => {
        setCard(data);
        setLoading(false);
      })
      .catch((err) => {
        setError("Failed to load random card.");
        setLoading(false);
      });
  }, []);

  return (
    <div>
      <h1>Random Magic Card</h1>
      {loading && <p>Loading...</p>}
      {error && <p style={{ color: "red" }}>{error}</p>}
      {card && (
        <div>
          <h2>{card.name}</h2>
          <p><strong>Mana Cost:</strong> {card.mana_cost || "No cost"}</p>
          <p><strong>Oracle Text:</strong> {card.oracle_text || "No text available"}</p>
        </div>
      )}
      <button onClick={() => window.location.reload()}>Get Another Card</button>
      <button onClick={() => (window.location.href = "/")}>Back to Home</button>
    </div>
  );
}