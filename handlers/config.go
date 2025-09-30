package handlers

import (
	"encoding/json"
	"net/http"
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
