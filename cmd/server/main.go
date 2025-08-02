package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/kireeti-28/mini-redis/internal/storage"
)

const filename = "storage.log"

type apiConfig struct {
	port   string
	store  *storage.Storage
	logger *slog.Logger
}

func main() {
	storage, err := storage.NewStorage(filename)
	if err != nil {
		fmt.Printf("Error initializing storage: %v\n", err)
		return
	}

	cfg := apiConfig{
		port:   "9686",
		store:  storage,
		logger: slog.Default(),
	}

	http.HandleFunc("GET /kv", cfg.getAllHandler)
	http.HandleFunc("GET /kv/{key}", cfg.getHandler)
	http.HandleFunc("POST /kv/{key}", cfg.setHandler)
	http.HandleFunc("DELETE /kv/{key}", cfg.deleteHandler)

	fmt.Println(http.ListenAndServe(":"+cfg.port, nil))
}

func (cfg *apiConfig) getAllHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, cfg.store.GetAll())
}

func (cfg *apiConfig) getHandler(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	value, err := cfg.store.Get(key)

	if err != nil {
		cfg.logger.Error("Failed to get value", "key", key, "error", err)
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"Key": key, "Value": value})
}

func (cfg *apiConfig) setHandler(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	data, err := io.ReadAll(r.Body)
	if err != nil {
		cfg.logger.Error("Failed to read request body", "error", err)
		respondWithError(w, http.StatusInternalServerError, "unable to read body", err)
		return
	}

	value := string(data)
	err = cfg.store.Set(key, value)
	if err != nil {
		cfg.logger.Error("Failed to set value", "key", key, "value", value, "error", err)
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"Key": key, "Value": value})
}

func (cfg *apiConfig) deleteHandler(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	err := cfg.store.Delete(key)
	if err != nil {
		cfg.logger.Error("Failed to delete key", "key", key, "error", err)
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
