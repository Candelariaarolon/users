import { useState } from "react";
import { useNavigate } from "react-router-dom";
import "./style_login.css";  // Asegúrate de que este archivo CSS exista

const Login = () => {
  const [rightPanelActive] = useState(false);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();

  const handleLogin = async (e) => {
    e.preventDefault();
    setError("");

    try {
      const response = await fetch("http://localhost:8080/users/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });

      if (!response.ok) {
        const err = await response.json();
        throw new Error(err.message || "Login failed");
      }

      const data = await response.json();
      document.cookie = `token=${data.token}; path=/; SameSite=Strict`;
      localStorage.setItem("userName", data.name);
      localStorage.setItem("surname", data.surname);
      navigate("/activities");
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div className={`container ${rightPanelActive ? "right-panel-active" : ""}`} id="container">

      {/* SIGN IN */}
      <div className="form-container sign-in-container">
        <form onSubmit={handleLogin}>
          <h1>Ingresar</h1>
          {error && <p className="error">{error}</p>}
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={e => setEmail(e.target.value)}
            required
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={e => setPassword(e.target.value)}
            required
          />
          <button type="submit">Ingresar</button>
        </form>
      </div>


    </div>
  );
};

export default Login;
