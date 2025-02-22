package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
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

func (state *AppState) queryDevicesHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	log.Printf("headers: %v", r.Header)

	resp := state.toQueryDeviceResponse()
	resp.RequestID = r.Header.Get("X-Request-Id")
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

// func (state *AppState) changeDevicesHandle(w http.ResponseWriter, r *http.Request) {
// 	decoder := json.NewDecoder(r.Body)
// 	var payload ChangeDevicesRequest
// 	err := decoder.Decode(&payload)
// 	if err != nil {
// 		writeError(&w, fmt.Sprintf("Error parsing body: %s", err))
// 		w.WriteHeader(http.StatusUnprocessableEntity)
// 		return
// 	}

// 	for _, device := range payload.Payload.Devices {
// 		state_device, found := state.Devices[device.ID]
// 		if !found {
// 			writeError(&w, fmt.Sprintf("Device ID %s not found", device.ID))
// 			w.WriteHeader(http.StatusNotFound)
// 			return
// 		}
// 		for _, cap := range device.Capabilities {
// 			state_device.Characteristics.updateCapability(cap.)
// 		}
// 	}
// }

func (state *AppState) toQueryDeviceResponse() QueryDevicesResponse {
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

	return QueryDevicesResponse{
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
