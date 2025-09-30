package handlers

import (
	"encoding/json"
	"net/http"
	"projekat/model"
	"projekat/services"
	"strconv"

	"github.com/gorilla/mux"
)

type ConfigGroupHandler struct {
	service services.ConfigGroupService
}

func NewConfigGroupHandler(service services.ConfigGroupService) ConfigGroupHandler {
	return ConfigGroupHandler{service: service}
}

// POST /config-groups
func (h ConfigGroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var grp model.ConfigGroup
	if err := json.NewDecoder(r.Body).Decode(&grp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if grp.Name == "" || grp.Version == 0 {
		http.Error(w, "name and version are required", http.StatusBadRequest)
		return
	}
	if err := h.service.Add(grp); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	resp, _ := json.Marshal(grp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

// GET /config-groups/{name}/{version}
func (h ConfigGroupHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	grp, err := h.service.Get(name, versionInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	resp, _ := json.Marshal(grp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// DELETE /config-groups/{name}/{version}
func (h ConfigGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.Delete(name, versionInt); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// POST /config-groups/{name}/{version}/configs
func (h ConfigGroupHandler) AddConfig(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var cfg model.GroupConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if cfg.Name == "" {
		http.Error(w, "config name is required", http.StatusBadRequest)
		return
	}
	if err := h.service.AddConfigToGroup(name, versionInt, cfg); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	grp, _ := h.service.Get(name, versionInt)
	resp, _ := json.Marshal(grp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// DELETE /config-groups/{name}/{version}/configs/{configName}
func (h ConfigGroupHandler) RemoveConfig(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	configName := mux.Vars(r)["configName"]
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.RemoveConfigFromGroup(name, versionInt, configName); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}


