package main

type QueryDeviceResponse struct {
	RequestID string  `json:"request_id"`
	Payload   Payload `json:"payload"`
}

type GetDevicesResponse struct {
	RequestID string  `json:"request_id"`
	Payload   Payload `json:"payload"`
}

type Payload struct {
	UserID  string   `json:"user_id"`
	Devices []Device `json:"devices"`
}

type Device struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Room         string                 `json:"room,omitempty"`
	Type         string                 `json:"type,omitempty"`
	CustomData   map[string]interface{} `json:"custom_data,omitempty"`
	Capabilities map[string]interface{} `json:"capabilities"`
	Properties   map[string]interface{} `json:"properties"`
	DeviceInfo   DeviceInfo             `json:"device_info,omitempty"`
	ErrorCode    string                 `json:"error_code,omitempty"`
	ErrorMessage string                 `json:"error_message,omitempty"`
}

type DeviceInfo struct {
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	HWVersion    string `json:"hw_version"`
	SWVersion    string `json:"sw_version"`
}
