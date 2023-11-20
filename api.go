package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

// newAPIServer creates a new instance of APIServer with the specified listen address.
// parameters: - listenAddr: The address on which the API server will listen for incoming requests.
// returns: a pointer to the newly created APIServer instance.
func newAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

// run is a method of APIServer that initializes and starts the HTTP server.
// it configures the router, setting up the "/account" endpoint to be handled by the handleAccount method.
// this method does not return anything, as it is expected to run indefinitely, handling incoming HTTP requests.
// create a new instance of the Gorilla Mux router.
func (s *APIServer) Run() {
	router := mux.NewRouter()
	// register the handleAccount method to be called when a request is made to the "/account" endpoint.
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))

	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccountByID))

	log.Println("JSON API SERVER RUNNING ON PORT: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)

}

// handleAccount is a method of APIServer that handles HTTP requests related to account operations.
// parameters: - w: http.ResponseWriter is used to construct the HTTP response.
// - r: *http.Request represents the incoming HTTP request.
// returns:	an error, if any occurred during the handling of the account request; otherwise, it returns nil.
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}

	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("method not allowed! %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.getAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	fmt.Println(id)
	// write a JSON response containing the account information to the HTTP response writer (w).
	return WriteJSON(w, http.StatusOK, &Account{})
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	// Decode the JSON data from the HTTP request body (r.Body) into the createAccountRequest struct.
	// json.NewDecoder creates a new JSON decoder that reads from the request body.
	// Decode method interprets and translates the JSON data into the fields of the createAccountRequest struct.
	// If there is an error during decoding, store the error in the 'err' variable.
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		log.Println("Error decoding JSON:", err)
		return err
	}
	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// writeJSON is a utility function for writing JSON responses in HTTP handlers.
// it sets the HTTP status code, content type, and encodes the provided value (v)
// as JSON to the response writer (w).
// parameters:
//   - w: http.ResponseWriter is used to construct the HTTP response.
//   - status: HTTP status code to be set in the response.
//   - v: any value to be encoded as JSON and included in the response body.

// returns:	An error, if any occurred during the encoding or writing process; otherwise, it returns nil.
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application./json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

// makeHTTPHandlerFunc is a higher-order function that takes an `apiFunc` as a parameter
// and returns an `http.HandlerFunc`. This function is designed to handle HTTP requests
// by calling the provided `apiFunc` and handling any potential errors.
// parameters:  - f: The `apiFunc` to be executed when handling the HTTP request.
// it represents a function that processes the HTTP request and may return an error.
// returns: an `http.HandlerFunc` that encapsulates the functionality of calling `f`
// and handling any errors that may occur during its execution.
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// call the provided `apiFunc` to handle the HTTP request, and check for errors.
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
