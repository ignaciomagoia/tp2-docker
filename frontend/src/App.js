import React, { useState } from "react";

function App() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");

  // Registro
  const register = async () => {
    try {
      const res = await fetch("http://localhost:8081/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
      });

      if (!res.ok) {
        throw new Error("Error en el registro");
      }

      const data = await res.json();
      setMessage("✅ Registro exitoso: " + data.message);
    } catch (err) {
      setMessage("❌ " + err.message);
    }
  };

  // Login
  const login = async () => {
    try {
      const res = await fetch("http://localhost:8081/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
      });

      if (!res.ok) {
        throw new Error("Error en el login");
      }

      const data = await res.json();
      setMessage("✅ Login exitoso. Token: " + data.token);
    } catch (err) {
      setMessage("❌ " + err.message);
    }
  };

  return (
    <div className="App" style={{ padding: "20px", fontFamily: "Arial" }}>
      <h1>Frontend Auth</h1>

      <div>
        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          style={{ margin: "5px" }}
        />
        <br />
        <input
          type="password"
          placeholder="Contraseña"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          style={{ margin: "5px" }}
        />
      </div>

      <div style={{ marginTop: "10px" }}>
        <button onClick={register} style={{ marginRight: "10px" }}>
          Registrarse
        </button>
        <button onClick={login}>Login</button>
      </div>

      <p style={{ marginTop: "20px", fontWeight: "bold" }}>{message}</p>
    </div>
  );
}

export default App;
