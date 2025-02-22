package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

type AppState struct {
	UserId  string
	Devices map[string]*Device
}

func (state *AppState) unlinkUserHandle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	resp := make(map[string]string)
	resp["request_id"] = r.Header.Get("X-Request-Id")
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func (state *AppState) patchDeviceStateHandle(w http.ResponseWriter, r *http.Request) {
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

	var result UpdateDeviceResult
	switch payload.Type {
	case "capability":
		// result = device.Characteristics.updateCapability(payload.Instance, *payload.Value)
	case "property":
		result = device.Characteristics.updateProperty(payload.Instance, *payload.Value)
	default:
		writeError(&w, "Error parsing body: 'type' must be 'capability' or 'property'")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if result.status != OK {
		writeError(&w, result.msg)
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

func (state *AppState) changeDevicesStateHandle(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var payload ChangeDevicesRequest
	err := decoder.Decode(&payload)
	if err != nil {
		writeError(&w, fmt.Sprintf("Error parsing body: %s", err))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	devices_response := []DeviceResponse{}
	for _, device := range payload.Payload.Devices {
		state_device, found := state.Devices[device.ID]
		if !found {
			writeError(&w, fmt.Sprintf("Device ID %s not found", device.ID))
			w.WriteHeader(http.StatusNotFound)
			return
		}
		device_response := DeviceResponse{ID: device.ID}
		var capabilities []map[string]interface{}
		for idx, cap := range device.Capabilities {
			state_map, found := cap["state"]
			state, ok := state_map.(map[string]interface{})
			if !found || !ok {
				writeError(&w, fmt.Sprintf("Missing field 'state' in capability %d", idx))
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
			var device_state StateResponse
			err := mapstructure.Decode(state, &device_state)
			if err != nil {
				writeError(&w, fmt.Sprintf("Field 'state' in capability %d has unexpected format", idx))
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}

			result := state_device.Characteristics.updateCapability(device_state.Instance, device_state.Value, token)
			var new_state StateResponse
			new_state.Instance = state["instance"].(string)
			if result.status == OK {
				new_state.ActionResult.Status = "DONE"
			} else {
				new_state.ActionResult.Status = "ERROR"
				switch result.status {
				case InvalidValue:
					new_state.ActionResult.ErrorCode = "INVALID_VALUE"
				case InvalidAction:
					new_state.ActionResult.ErrorCode = "INVALID_ACTION"
				case NotEnoughWater:
					new_state.ActionResult.ErrorCode = "NOT_ENOUGH_WATER"
				case DeviceUnreachable:
					new_state.ActionResult.ErrorCode = "DEVICE_UNREACHABLE"
				case UnknownError:
					new_state.ActionResult.ErrorCode = "INTERNAL_ERROR"
				}
			}

			cap_response := make(map[string]interface{})
			cap_response["type"] = cap["type"]
			cap_response["state"] = new_state
			capabilities = append(capabilities, cap_response)
		}
		device_response.Capabilities = capabilities
		devices_response = append(devices_response, device_response)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	var resp ChangeDevicesResponse
	resp.RequestID = r.Header.Get("X-Request-Id")
	resp.Payload.Devices = devices_response
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func (state *AppState) toQueryDeviceResponse() QueryDevicesResponse {
	devices := []DeviceResponse{}
	for device_id, device := range state.Devices {
		capabilities := CPListToMapList(device.Characteristics.capabilities())
		properties := CPListToMapList(device.Characteristics.properties())
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
		capabilities := CPListToMapList(device.Characteristics.capabilities())
		properties := CPListToMapList(device.Characteristics.properties())
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
