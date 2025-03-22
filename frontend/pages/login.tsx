import { useState } from "react";

export default function LoginPage() {
  const [creds, setCreds] = useState({ email: "", password: "" });
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setMessage("");
    setError("");

    try {
      const res = await fetch("http://localhost:8080/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(creds),
      });

      if (res.ok) {
        const user = await res.json();
        setMessage(`Logged in! Hello, ${user.username}`);
        // Save user info or token somewhere if you plan to maintain session
      } else {
        const text = await res.text();
        setError(`Error: ${text}`);
      }
    } catch (err: any) {
      setError(`Request failed: ${err.message}`);
    }
  };

  return (
    <div style={{ maxWidth: "400px", margin: "0 auto" }}>
      <h1>Login</h1>
      {message && <p style={{ color: "green" }}>{message}</p>}
      {error && <p style={{ color: "red" }}>{error}</p>}
      <form onSubmit={handleSubmit}>
        <div>
          <label>Email:</label>
          <input
            type="email"
            required
            value={creds.email}
            onChange={(e) => setCreds({ ...creds, email: e.target.value })}
          />
        </div>

        <div>
          <label>Password:</label>
          <input
            type="password"
            required
            value={creds.password}
            onChange={(e) => setCreds({ ...creds, password: e.target.value })}
          />
        </div>

        <button type="submit">Login</button>
      </form>
    </div>
  );
}