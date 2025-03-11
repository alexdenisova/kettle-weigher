package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/mitchellh/mapstructure"
)

func hashPassword(password string) []byte {
	hashed_pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error hashing password: %s", err)
	}
	return hashed_pass
}

func parsePassword(basic_auth string) (string, error) {
	basic_auth = strings.TrimPrefix(basic_auth, "Basic ")
	basic_auth = strings.TrimSpace(basic_auth)
	decoded, err := base64.StdEncoding.DecodeString(basic_auth)
	password := strings.TrimPrefix(string(decoded), ":")
	return password, err
}

func writeError(w *http.ResponseWriter, msg string) {
	err_msg := ErrorMessage{
		Message: msg,
	}
	log.Printf("Error: %s", msg)
	jsonResp, _ := json.Marshal(err_msg)
	(*w).Write(jsonResp)
}

func CPListToMapList(cp_list []CapabilityProperty) []map[string]interface{} {
	m := []map[string]interface{}{}
	for _, cp := range cp_list {
		result := map[string]interface{}{}
		mapstructure.Decode(cp, &result)
		if strings.Contains(cp.Type, "properties") {
			result["parameters"] = cp.State.toParameters()
		}
		result["state"] = cp.State.toState()
		m = append(m, result)
	}
	return m
}

func CPToActionRequest(cp CapabilityProperty) ActionRequest {
	return ActionRequest{
		Type: cp.Type,
		State: StateResponse{
			Instance: cp.State.Instance,
			Value:    cp.State.Value,
		},
	}
}
