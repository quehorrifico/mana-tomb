import { useEffect, useState } from "react";
import { useRouter } from "next/router";

interface Deck {
  deck_id: number;
  user_id: number;
  name: string;
  description: string;
  commander: string;
  cards: string[];
}

export default function EditDeck() {
  const router = useRouter();
  const { deckID } = router.query;
  const [deck, setDeck] = useState<Deck | null>(null);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [commander, setCommander] = useState("");
  const [cards, setCards] = useState<string>("");

  useEffect(() => {
    if (!deckID) return;

    fetch(`/api/decks/${deckID}`)
      .then((res) => res.json())
      .then((data) => {
        setDeck(data);
        setName(data.name);
        setDescription(data.description);
        setCommander(data.commander);
        setCards(data.cards.join(", "));
      })
      .catch((err) => console.error("Error fetching deck:", err));
  }, [deckID]);

  const handleUpdate = async () => {
    if (!deck) return;

    const updatedDeck = {
      deck_id: deck.deck_id,
      name,
      description,
      commander,
      cards: cards.split(",").map((card) => card.trim()),
    };

    try {
      const res = await fetch(`/api/decks/update/${deckID}`, { // might need to change to updatedDeck.deck_id
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(updatedDeck),
      });

      if (res.ok) {
        alert("Deck updated!");
        router.push(`/decks/${deck.deck_id}`);
      } else {
        alert("Failed to update deck.");
      }
    } catch (error) {
      console.error("Error updating deck:", error);
      alert("Something went wrong.");
    }
  };

  if (!deck) return <p>Loading deck...</p>;

  return (
    <div>
      <h1>Edit Deck: {deck.name}</h1>
      <label>
        Name:
        <input value={name} onChange={(e) => setName(e.target.value)} />
      </label>
      <br />
      <label>
        Description:
        <input value={description} onChange={(e) => setDescription(e.target.value)} />
      </label>
      <br />
      <label>
        Commander:
        <input value={commander} onChange={(e) => setCommander(e.target.value)} />
      </label>
      <br />
      <label>
        Cards (comma-separated):
        <textarea value={cards} onChange={(e) => setCards(e.target.value)} />
      </label>
      <br />
      <button onClick={handleUpdate}>Save Changes</button>
      <button onClick={() => router.push("/deck-building")}>Back to Decks</button>
    </div>
  );
}