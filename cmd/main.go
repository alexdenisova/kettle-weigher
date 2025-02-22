package main

import (
	"log"
	"net/http"
)

func main() {
	user_id := "alex"
	devices := map[string]*Device{}

	kettle_weigher := KettleWeigher{}
	kettle_weigher_info := Device{
		Name:            "Весы для чайника",
		Description:     "Измеряет вес чайника с водой",
		Room:            "кухня",
		Type:            "devices.types.other",
		Characteristics: &kettle_weigher,
		DeviceInfo: &DeviceInfo{
			Manufacturer: "Alex Denisova",
			Model:        "kettle-weigher",
			Version:      "1.0",
		},
	}
	devices["kettle-weigher"] = &kettle_weigher_info

	app_state := AppState{
		UserId:  user_id,
		Devices: devices,
	}
	log.Printf("Starting server")
	http.HandleFunc("GET /health", healthHandler)
	http.HandleFunc("PATCH /v1.0/user/device/state", app_state.patchDeviceState)
	http.HandleFunc("GET /v1.0/user/devices", app_state.getDevicesHandle)
	http.HandleFunc("POST /v1.0/user/devices/query", app_state.queryDevicesHandle)
	http.ListenAndServe(":8080", nil)
}
