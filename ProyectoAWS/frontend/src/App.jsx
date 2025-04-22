import { useState } from "react";
import MonacoEditor from "@monaco-editor/react";
import { Link } from 'react-router-dom';
import './App.css'; 

function App() {
  // almacenar texto de entrada
  const [inputText, setInputText] = useState("");
  // almacenar texto después de ser procesado
  const [outputText, setOutputText] = useState("");

  // obtener la URL del backend desde localStorage (o usar localhost por defecto)
  const backendURL = localStorage.getItem("backend_url") || "http://localhost:8080";

  // Ejecutar solicitud al backend
  const handleExecute = async () => {
    try {
      const response = await fetch(`${backendURL}/api/analyze`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ input: inputText }),
      });

      const data = await response.text(); // usas Fprintln así que espera texto
      setOutputText(data);
    } catch (error) {
      console.error("Error:", error);
      setOutputText("Error al conectar con el backend");
    }
  };

  return (
    <div className="container">
      <h1 className="title">Analizador</h1>

      <Link to="/about" style={{ marginBottom: "1rem", display: "inline-block", color: "#61dafb" }}>
        Configurar IP del Backend
      </Link>

      <div className="textarea-container">
        <MonacoEditor
          value={inputText}
          onChange={(value) => setInputText(value)}
          height="300px"
          language="javascript"
          theme="vs-dark"
          options={{
            selectOnLineNumbers: true,
          }}
        />
        <textarea
          value={outputText}
          readOnly
          placeholder="Resultado..."
          className="output-textarea"
        />
      </div>

      <button className="execute-button" onClick={handleExecute}>
        Ejecutar
      </button>
    </div>
  );
}

export default App;
