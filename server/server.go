package server

import (
	"encoding/json"
	"log"
	"net/http"
)

const PORT = ":8080"

type JSON struct {
	Message string `json:"message"`
}

func Serve() error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		res := JSON{
			Message: "Dette er en test",
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Println("Error encoding response:", err)
		}
	})

	return http.ListenAndServe(PORT, mux)
}
