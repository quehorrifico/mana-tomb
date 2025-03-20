import React, { useEffect, useState } from "react";

interface Card {
  name: string;
  mana_cost: string;
  image_uris: {
    small: string;
    normal: string;
    large: string;
    png: string;
    art_crop: string;
    border_crop: string;
  };
  type_line: string;
  oracle_text: string;
  set: string;
  set_name: string;
  set_uri: string;
  set_id: string;
  set_type: string;
  set_search_uri: string;
  scryfall_set_uri: string;
}

export default function RandomCard() {
  const [card, setCard] = useState<Card | null>(null);
  console.log("Card Data:", card);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
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
          <p><strong>Type:</strong> {card.type_line}</p>
          <img src={card.image_uris.normal} alt={card.name} />
          <p><strong>Set:</strong> {card.set}</p>
          <p><strong>Set Name:</strong> {card.set_name}</p>
          <p><strong>Set URI:</strong> {card.set_uri}</p>
          <p><strong>Set ID:</strong> {card.set_id}</p>
          <p><strong>Set Type:</strong> {card.set_type}</p>
          <p><strong>Set Search URI:</strong> {card.set_search_uri}</p>
          <p><strong>Scryfall Set URI:</strong> {card.scryfall_set_uri}</p>
        </div>
      )}
      <button onClick={() => window.location.reload()}>Get Another Card</button>
      <button onClick={() => (window.location.href = "/")}>Back to Home</button>
    </div>
  );
}