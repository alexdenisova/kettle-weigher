package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type HelloResponse struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("GET /health", healthHandler)
	http.HandleFunc("GET /v1.0/user/devices", getDevicesHandle)
	http.HandleFunc("POST /v1.0/user/devices/query", postDevicesHandle)
	http.HandleFunc("/", other)
	http.ListenAndServe(":8080", nil)
}

func other(w http.ResponseWriter, r *http.Request) {
	log.Printf(r.URL.Path)
	w.WriteHeader(http.StatusOK)
}

func getDevicesHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	test_device := getDevicesResponse()
	test_device.RequestID = r.Header.Get("X-Request-Id")
	jsonResp, err := json.Marshal(test_device)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func postDevicesHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	test_device := queryDevicesResponse()
	test_device.RequestID = r.Header.Get("X-Request-Id")
	jsonResp, err := json.Marshal(test_device)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}
