package main

import (
	"fmt"
	"net/http"
)

// /health
func healthHandler(w http.ResponseWriter, r *http.Request) {
	if serverIsHealthy() {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Server is healthy")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Server is not healthy")
	}
}

func serverIsHealthy() bool {
	return true
}
