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
	router.HandleFunc("/api/cluster/{clustername}/connect", ConnectCluster).Methods("POST")
	// Create - add a cluster
	router.HandleFunc("/api/cluster", CreateCluster).Methods("POST")
	// Read
	router.HandleFunc("/api/cluster/{clustername}", GetCluster).Methods("GET")
	// Read-all
	router.HandleFunc("/api/cluster", GetClusters).Methods("GET")
	// Update
	router.HandleFunc("/api/cluster/{clustername}", UpdateCluster).Methods("PUT")
	// Delete
	router.HandleFunc("/api/cluster/{clustername}", DeleteCluster).Methods("DELETE")

	//
	// Service
	//

	// Create
	router.HandleFunc("/api/service", CreateService).Methods("POST")
	// Read
	router.HandleFunc("/api/service/{name}", GetService).Methods("GET")
	// Read-all
	router.HandleFunc("/api/service", GetServices).Methods("GET")
	// Update
	router.HandleFunc("/api/service/{name}", UpdateService).Methods("PUT")
	// Delete
	router.HandleFunc("/api/service/{name}", DeleteService).Methods("DELETE")

	//
	// Factory
	//

	// Create
	router.HandleFunc("/api/factory", CreateFactory).Methods("POST")
	// Read
	router.HandleFunc("/api/factory/{name}", GetFactory).Methods("GET")
	// Read-all
	router.HandleFunc("/api/factory", GetFactories).Methods("GET")
	// Update
	router.HandleFunc("/api/factory/{name}", UpdateFactory).Methods("PUT")
	// Delete
	router.HandleFunc("/api/factory/{name}", DeleteFactory).Methods("DELETE")

	//
	// Group
	//

	// Create
	//	router.HandleFunc("/api/group", CreateGroup).Methods("POST")
	// Add Member
	//	router.HandleFunc("/api/group/{name}/member", AddGroupMember).Methods("POST")
	// Read
	//	router.HandleFunc("/api/group/{name}", GetGroup).Methods("GET")
	// Read-all
	//	router.HandleFunc("/api/group", GetGroup).Methods("GET")
	// Update
	//	router.HandleFunc("/api/group/{name}", UpdateGroup).Methods("PUT")
	// Delete
	//	router.HandleFunc("/api/group/{name}", DeleteGroup).Methods("DELETE")

	return router
}
