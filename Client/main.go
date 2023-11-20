package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Your existing chaincode functions
// ...
func handleCreateLand(w http.ResponseWriter, r *http.Request) {
	// Extract data from the request if needed
	// ...
	params := mux.Vars(r)
	landID := params["landID"]
    Ownername := params["ownerName"]
	// Use the provided chaincode function
	result, err := submitTxnFn(
		"manufacturer",
		"autochannel",
		"KBA-Automobile",
		"LandContract",
		"invoke",
		make(map[string][]byte),
		"CreateLand",
		landID, Ownername,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the result in JSON format
	response := map[string]string{"result": result}
	json.NewEncoder(w).Encode(response)
}

// Handler function for updating land ownership
func handleUpdateLandOwnership(w http.ResponseWriter, r *http.Request) {
	// Extract data from the request if needed
	// ...

	// Use the provided chaincode function
	result, err := submitTxnFn(
		"manufacturer",
		"autochannel",
		"KBA-Automobile",
		"LandContract",
		"invoke",
		make(map[string][]byte),
		"UpdateLandOwnership",
		"Land01", "NewOwnerName",
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the result in JSON format
	response := map[string]string{"result": result}
	json.NewEncoder(w).Encode(response)
}

// Handler function for querying land by ID
func handleReadLand(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	landID := params["landID"]

	// Use the provided chaincode function
	result, err := submitTxnFn(
		"manufacturer",
		"autochannel",
		"KBA-Automobile",
		"LandContract",
		"query",
		make(map[string][]byte),
		"ReadLand",
		landID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the result in JSON format
	response := map[string]string{"result": result}
	json.NewEncoder(w).Encode(response)
}

// Handler function for deleting a land
func handleDeleteLand(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	landID := params["landID"]

	// Use the provided chaincode function
	result, err := submitTxnFn(
		"manufacturer",
		"autochannel",
		"KBA-Automobile",
		"LandContract",
		"invoke",
		make(map[string][]byte),
		"DeleteLand",
		landID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the result in JSON format
	response := map[string]string{"result": result}
	json.NewEncoder(w).Encode(response)
}

// Handler function for getting lands by range
func handleGetLandsByRange(w http.ResponseWriter, r *http.Request) {
	// Extract data from the request if needed
	// ...

	// Use the provided chaincode function
	result, err := submitTxnFn(
		"manufacturer",
		"autochannel",
		"KBA-Automobile",
		"LandContract",
		"query",
		make(map[string][]byte),
		"GetLandsByRange",
		"startKey", "endKey",
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the result in JSON format
	response := map[string]string{"result": result}
	json.NewEncoder(w).Encode(response)
}

// Add more handler functions for other missing chaincode functions as needed

func main() {
	// Create a new Gorilla Mux router
	router := mux.NewRouter()

	// Define your API endpoints
	router.HandleFunc("/api/land/create/{landID}/{ownerName}", handleCreateLand).Methods("POST")
	router.HandleFunc("/api/land/update/{landID}/{newOwner}", handleUpdateLandOwnership).Methods("PUT")
	router.HandleFunc("/api/land/read/{landID}", handleReadLand).Methods("GET")
	router.HandleFunc("/api/land/delete/{landID}", handleDeleteLand).Methods("DELETE")
	router.HandleFunc("/api/land/range", handleGetLandsByRange).Methods("GET")

	// Start the HTTP server
	fmt.Println("API server listening on port 8080")
	http.ListenAndServe(":8500", router)
}