import { useRouter } from "next/router";

export default function PlaceholderPage() {
  const router = useRouter();

  return (
    <div style={{ textAlign: "center", padding: "50px" }}>
      <h1>🚧 This Feature is Under Construction 🚧</h1>
      <p>I'm still working on this feature! Check back soon for updates. 😊</p>
      <button onClick={() => router.push("/")}>Back to Home</button>
    </div>
  );
}