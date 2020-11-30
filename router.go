package main

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	// Login
	router.HandleFunc("/api/session", Login).Methods("POST")

	// Session
	router.HandleFunc("/api/currentuser", GetSession).Methods("GET")

	//
	// Cluster
	//

	// Connect
	router.HandleFunc("/api/cluster/{id}/connect", ConnectCluster).Methods("POST")
	// Create
	router.HandleFunc("/api/cluster", CreateCluster).Methods("POST")
	// Read
	router.HandleFunc("/api/cluster/{id}", GetCluster).Methods("GET")
	// Read-all
	router.HandleFunc("/api/cluster", GetClusters).Methods("GET")
	// Update
	router.HandleFunc("/api/cluster/{id}", UpdateCluster).Methods("PUT")
	// Delete
	router.HandleFunc("/api/cluster/{id}", DeleteCluster).Methods("DELETE")

	//
	// Service
	//

	// Create
	router.HandleFunc("/api/service", CreateService).Methods("POST")
	// Read
	router.HandleFunc("/api/service/{id}", GetService).Methods("GET")
	// Read-all
	router.HandleFunc("/api/service", GetServices).Methods("GET")
	// Update
	router.HandleFunc("/api/service/{id}", UpdateService).Methods("PUT")
	// Delete
	router.HandleFunc("/api/service/{id}", DeleteService).Methods("DELETE")

	//
	// Factory
	//

	// Create
	router.HandleFunc("/api/factory", CreateFactory).Methods("POST")
	// Read
	router.HandleFunc("/api/factory/{id}", GetFactory).Methods("GET")
	// Read-all
	router.HandleFunc("/api/factory", GetFactories).Methods("GET")
	// Update
	router.HandleFunc("/api/factory/{id}", UpdateFactory).Methods("PUT")
	// Delete
	router.HandleFunc("/api/factory/{id}", DeleteFactory).Methods("DELETE")

	return router
}
