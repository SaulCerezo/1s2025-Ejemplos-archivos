package main

import (
	"ejemplo2/Analizador"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permitir peticiones desde cualquier origen
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// estructura para manejar los json
type RequestData struct {
	Input string `json:"input"`
}

// Handler para procesar las solicitudes a /api/analyze
func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Cuerpo de la solicitud - lectura
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Decodificar JSON
	var requestData RequestData
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	// Obtener el texto de la entrada
	input := requestData.Input

	fmt.Println("Recibido en el backend:", input)

	// Llamar al analizador con el input limpio
	Analizador.Analizador(input)

	// Responder al frontend
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Procesado correctamente")
}

func main() {
	r := mux.NewRouter()
	r.Use(enableCORS)
	// Definimos la mismar ruta del front /api/analyze para manejar POST y OPTIONS
	r.HandleFunc("/api/analyze", analyzeHandler).Methods("POST", "OPTIONS")

	fmt.Println("Servidor en ejecución en http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
