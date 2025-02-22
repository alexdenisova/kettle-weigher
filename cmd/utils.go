package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/mitchellh/mapstructure"
)

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
