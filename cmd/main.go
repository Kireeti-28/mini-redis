package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/kireeti-28/mini-redis/internal/storage"
)

type apiConfig struct {
	port  string
	store *storage.Storage
}

func main() {
	cfg := apiConfig{
		port:  "9686",
		store: storage.NewStorage(),
	}

	fmt.Printf("Server starting on port: %s\n", cfg.port)

	http.HandleFunc("GET /kv/{key}", cfg.getHandler)
	http.HandleFunc("POST /kv/{key}", cfg.setHandler)
	http.HandleFunc("DELETE /kv/{key}", cfg.deleteHandler)

	fmt.Println(http.ListenAndServe(":"+cfg.port, nil))
}

func (cfg *apiConfig) getHandler(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	value, err := cfg.store.Get(key)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"Key": key, "Value": value})
}

func (cfg *apiConfig) setHandler(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	data, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to read body", err)
		return
	}

	value := string(data)
	cfg.store.Set(key, value)
	respondWithJSON(w, http.StatusOK, map[string]string{"Key": key, "Value": value})
}

func (cfg *apiConfig) deleteHandler(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	err := cfg.store.Delete(key)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
