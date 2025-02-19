package main

func getDevicesResponse() GetDevicesResponse {
	return GetDevicesResponse{
		RequestID: "12345",
		Payload: Payload{
			UserID: "user_6789",
			Devices: []Device{
				{
					ID:          "device_001",
					Name:        "Living Room Light",
					Description: "Smart LED light in the living room",
					Room:        "Living Room",
					Type:        "light",
					CustomData: map[string]interface{}{
						"color":      "white",
						"brightness": 75,
					},
					Capabilities: map[string]interface{}{
						"on_off": map[string]interface{}{
							"supported": true,
						},
						"brightness": map[string]interface{}{
							"min": 0,
							"max": 100,
						},
					},
					Properties: map[string]interface{}{
						"power":            "on",
						"brightness_level": 75,
					},
					DeviceInfo: DeviceInfo{
						Manufacturer: "SmartHome Inc.",
						Model:        "LED123",
						HWVersion:    "1.0",
						SWVersion:    "2.1.3",
					},
				},
			},
		},
	}

}

func queryDevicesResponse() QueryDeviceResponse {
	return QueryDeviceResponse{
		RequestID: "12345",
		Payload: Payload{
			Devices: []Device{
				{
					ID: "device_001",
					Capabilities: map[string]interface{}{
						"on_off": map[string]interface{}{
							"supported": true,
						},
						"brightness": map[string]interface{}{
							"min": 0,
							"max": 100,
						},
					},
					Properties: map[string]interface{}{
						"power":            "on",
						"brightness_level": 75,
					},
				},
			},
		},
	}
}
