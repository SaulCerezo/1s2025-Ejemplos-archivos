import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";

export default function About() {
  const [backendURL, setBackendURL] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    const savedURL = localStorage.getItem("backend_url");
    if (savedURL) setBackendURL(savedURL);
  }, []);

  const handleChange = (e) => {
    setBackendURL(e.target.value);
  };

  const handleSave = () => {
    localStorage.setItem("backend_url", backendURL);
    navigate("/"); // redirigir al home
  };

  return (
    <div style={{ padding: "2rem" }}>
      <h2>Configuración del Backend</h2>
      <p>Ingresa la IP o URL del backend (ej: http://34.123.45.67:8080):</p>
      <input
        type="text"
        value={backendURL}
        onChange={handleChange}
        placeholder="http://tu-ip:8080"
        style={{
          width: "80%",
          padding: "10px",
          fontSize: "1rem",
          marginTop: "1rem",
        }}
      />
      <br />
      <button
        onClick={handleSave}
        style={{
          marginTop: "1rem",
          padding: "10px 20px",
          fontSize: "1rem",
          cursor: "pointer",
          backgroundColor: "#61dafb",
          border: "none",
          borderRadius: "5px",
        }}
      >
        Aceptar
      </button>
      <p style={{ marginTop: "1rem" }}><strong>Actual:</strong> {backendURL}</p>
    </div>
  );
}
