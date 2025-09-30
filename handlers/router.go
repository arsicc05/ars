package handlers

import (
	"github.com/gorilla/mux"
)

func BuildRouter(configHandler ConfigHandler, groupHandler ConfigGroupHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/configs", configHandler.Create).Methods("POST")
	r.HandleFunc("/configs/{name}/{version}", configHandler.Get).Methods("GET")
	r.HandleFunc("/configs/{name}/{version}", configHandler.Delete).Methods("DELETE")

	// groups
	r.HandleFunc("/config-groups", groupHandler.Create).Methods("POST")
	r.HandleFunc("/config-groups/{name}/{version}", groupHandler.Get).Methods("GET")
	r.HandleFunc("/config-groups/{name}/{version}", groupHandler.Delete).Methods("DELETE")
	r.HandleFunc("/config-groups/{name}/{version}/configs", groupHandler.AddConfig).Methods("POST")
	r.HandleFunc("/config-groups/{name}/{version}/configs/{configName}", groupHandler.RemoveConfig).Methods("DELETE")
	return r
}


