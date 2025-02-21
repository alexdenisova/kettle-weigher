package main

type GetDevicesResponse struct {
	RequestID string  `json:"request_id"`
	Payload   Payload `json:"payload"`
}

type QueryDeviceResponse struct {
	RequestID string  `json:"request_id"`
	Payload   Payload `json:"payload"`
}

type Payload struct {
	UserID  string           `json:"user_id,omitempty"`
	Devices []DeviceResponse `json:"devices"`
}

type DeviceResponse struct {
	ID           string                   `json:"id"`
	Name         string                   `json:"name,omitempty"`
	Description  string                   `json:"description,omitempty"`
	Room         string                   `json:"room,omitempty"`
	Type         string                   `json:"type,omitempty"`
	CustomData   map[string]interface{}   `json:"custom_data,omitempty"`
	Capabilities []map[string]interface{} `json:"capabilities,omitempty"`
	Properties   []map[string]interface{} `json:"properties"`
	DeviceInfo   *DeviceInfoResponse      `json:"device_info,omitempty"`
	ErrorCode    string                   `json:"error_code,omitempty"`
	ErrorMessage string                   `json:"error_message,omitempty"`
}

type DeviceInfoResponse struct {
	Manufacturer string `json:"manufacturer,omitempty"`
	Model        string `json:"model,omitempty"`
	HWVersion    string `json:"hw_version,omitempty"`
	SWVersion    string `json:"sw_version,omitempty"`
}

