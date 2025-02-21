package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

type AppState struct {
	UserId  string
	Devices map[string]*Device
}

func (state *AppState) patchDeviceState(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var payload DeviceStatePayload
	err := decoder.Decode(&payload)
	if err != nil {
		writeError(&w, fmt.Sprintf("Error parsing body: %s", err))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	validate := validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("json")
	})
	err = validate.Struct(payload)
	var validateErrs validator.ValidationErrors
	if errors.As(err, &validateErrs) {
		err_msg := validatorErrorString(validateErrs)
		writeError(&w, err_msg)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if err != nil {
		writeError(&w, fmt.Sprintf("Error parsing body: %s", err))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	device, found := state.Devices[payload.DeviceID]
	if !found {
		writeError(&w, "Device ID not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch payload.Type {
	case "capability":
		err = device.Characteristics.updateCapability(payload.Instance, *payload.Value)
	case "property":
		err = device.Characteristics.updateProperty(payload.Instance, *payload.Value)
	default:
		writeError(&w, "Error parsing body: 'type' must be 'capability' or 'property'")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Device %s successfully updated", payload.DeviceID)
	w.WriteHeader(http.StatusNoContent)
}

func (state *AppState) getDevicesHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	resp := state.toGetDevicesResponse()
	resp.RequestID = r.Header.Get("X-Request-Id")
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func (state *AppState) postDevicesHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	resp := state.toQueryDeviceResponse()
	resp.RequestID = r.Header.Get("X-Request-Id")
	jsonResp, _ := json.Marshal(resp)
	log.Printf("response: %v", string(jsonResp[:]))
	w.Write(jsonResp)
}

func writeError(w *http.ResponseWriter, msg string) {
	err_msg := ErrorMessage{
		Message: msg,
	}
	jsonResp, _ := json.Marshal(err_msg)
	(*w).Write(jsonResp)
}

func (state *AppState) toQueryDeviceResponse() QueryDeviceResponse {
	devices := []DeviceResponse{}
	for device_id, device := range state.Devices {
		capabilities := CPListtoMapList(device.Characteristics.capabilities())
		properties := CPListtoMapList(device.Characteristics.properties())
		devices = append(devices, DeviceResponse{
			ID:           device_id,
			Capabilities: capabilities,
			Properties:   properties,
			// ErrorCode:    "",
			// ErrorMessage: "",
		})
	}

	return QueryDeviceResponse{
		Payload: Payload{
			Devices: devices,
		},
	}
}

func (state *AppState) toGetDevicesResponse() GetDevicesResponse {
	devices := []DeviceResponse{}
	for device_id, device := range state.Devices {
		capabilities := CPListtoMapList(device.Characteristics.capabilities())
		properties := CPListtoMapList(device.Characteristics.properties())
		devices = append(devices, DeviceResponse{
			ID:           device_id,
			Name:         device.Name,
			Description:  device.Description,
			Room:         device.Room,
			Type:         device.Type,
			Capabilities: capabilities,
			Properties:   properties,
			DeviceInfo: &DeviceInfoResponse{
				Manufacturer: device.DeviceInfo.Manufacturer,
				Model:        device.DeviceInfo.Model,
				HWVersion:    device.DeviceInfo.Version,
				SWVersion:    device.DeviceInfo.Version,
			},
			// ErrorCode:    "",
			// ErrorMessage: "",
		})
	}

	return GetDevicesResponse{
		Payload: Payload{
			UserID:  state.UserId,
			Devices: devices,
		},
	}
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
