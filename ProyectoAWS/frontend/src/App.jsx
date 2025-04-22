import { useState } from "react";
import MonacoEditor from "@monaco-editor/react";
import './App.css'; 

function App() {
  //almacenar texto de la entrada
  const [inputText, setInputText] = useState("");
  //almacena texto despues de ser procesado
  const [outputText, setOutputText] = useState("");

  // Se ejecuta al presionar un botón para procesar
  const handleExecute = async () => {
    try {
      // Envia una solicitud POST a la API en localhost:8080
      const response = await fetch("http://localhost:8080/api/analyze", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ input: inputText }), // Envía el texto ingresado
      });

      // Convierte la respuesta en JSON
      const data = await response.json();
      // Actualiza el estado del servidor
      setOutputText(data.response);
    } catch (error) {
      console.error("Error:", error);
    }
  };

  return (
    <div className="container">
      <h1 className="title">Analizador</h1>
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
