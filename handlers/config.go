package handlers

import (
	"encoding/json"
	"net/http"
	"projekat/model"
	"projekat/services"
	"strconv"

	"github.com/gorilla/mux"
)

type ConfigHandler struct {
	service services.ConfigService
}

func NewConfigHandler(service services.ConfigService) ConfigHandler {
	return ConfigHandler{
		service: service,
	}
}

// GET /configs/{name}/{version}
func (c ConfigHandler) Get(w http.ResponseWriter, r *http.Request) {
	// get name and version
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// call service method
	config, err := c.service.Get(name, versionInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// return response
	resp, err := json.Marshal(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// POST /configs
func (c ConfigHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var cfg model.Config
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if cfg.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if cfg.Version == 0 {
		http.Error(w, "version is required", http.StatusBadRequest)
		return
	}
	if err := c.service.Add(cfg); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	resp, err := json.Marshal(cfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// DELETE /configs/{name}/{version}
func (c ConfigHandler) Delete(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := c.service.Delete(name, versionInt); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
