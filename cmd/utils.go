package main

import (
	"encoding/json"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

func writeError(w *http.ResponseWriter, msg string) {
	err_msg := ErrorMessage{
		Message: msg,
	}
	jsonResp, _ := json.Marshal(err_msg)
	(*w).Write(jsonResp)
}

func CPListtoMapList(cp_list []CapabilityProperty) []map[string]interface{} {
	m := []map[string]interface{}{}
	for _, cp := range cp_list {
		result := map[string]interface{}{}
		mapstructure.Decode(cp, &result)
		result["parameters"] = cp.State.toParameters()
		result["state"] = cp.State.toState()
		m = append(m, result)
	}
	return m
}
