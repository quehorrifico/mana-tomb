import React from "react";
import { useRouter } from "next/router";

export default function Home() {
  const router = useRouter();

  return (
    <div>
      <h1>Mana Tomb Home</h1>
      <button onClick={() => router.push("/random-card")}>Get a Random Card</button>
    </div>
  );
}