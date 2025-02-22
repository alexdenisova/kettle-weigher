package main

import (
	"log"
	"net/http"

	"github.com/spf13/viper"
)

func getEnv() *viper.Viper {
	env := viper.New()
	env.SetEnvPrefix("KW_")
	env.BindEnv("min_water_level") // KW__MIN_WATER_LEVEL
	env.SetDefault("min_water_level", "20")
	env.BindEnv("kettle_id") // KW__KETTLE_ID
	return env
}

func main() {
	env := getEnv()

	user_id := "alex"
	devices := map[string]*Device{}

	kettle_weigher := KettleWeigher{
		kettle_id:       env.GetString("kettle_id"),
		min_water_level: float32(env.GetFloat64("min_water_level")),
	}
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
	http.HandleFunc("GET /v1.0", healthHandler)
	http.HandleFunc("POST /v1.0/user/unlink", app_state.unlinkUserHandle)
	http.HandleFunc("PATCH /v1.0/user/device/state", app_state.patchDeviceStateHandle)
	http.HandleFunc("GET /v1.0/user/devices", app_state.getDevicesHandle)
	http.HandleFunc("POST /v1.0/user/devices/query", app_state.queryDevicesHandle)
	http.HandleFunc("POST /v1.0/user/devices/action", app_state.changeDevicesStateHandle)
	http.ListenAndServe(":8080", nil)
}
