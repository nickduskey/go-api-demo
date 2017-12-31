package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type jwtToken struct {
	Token string `json:"token"`
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if len(vars["username"]) == 0 || len(vars["password"]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please provide username and password to obtain the token"))
		return
	}

	// Query the DB for name and check against password

	// Issue the JWT
}
