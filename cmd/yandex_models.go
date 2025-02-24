package main

type GetDevicesResponse struct {
	RequestID string  `json:"request_id"`
	Payload   Payload `json:"payload"`
}

type QueryDevicesResponse struct {
	RequestID string  `json:"request_id"`
	Payload   Payload `json:"payload"`
}

type ChangeDevicesRequest struct {
	Payload Payload `json:"payload"`
}

type ChangeDevicesResponse struct {
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
	Properties   []map[string]interface{} `json:"properties,omitempty"`
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

type StateResponse struct {
	Instance     string               `json:"instance,omitempty"`
	Value        interface{}          `json:"value,omitempty"`
	ActionResult ActionResultResponse `json:"action_result,omitempty"`
}

type ActionResultResponse struct {
	Status       string `json:"status"`
	ErrorCode    string `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type ChangeDevicesStateRequest struct {
	Devices []ChangeDeviceStateRequest `json:"devices"`
}

type ChangeDeviceStateRequest struct {
	ID      string          `json:"id"`
	Actions []ActionRequest `json:"actions"`
}

type ActionRequest struct {
	Type  string        `json:"type"`
	State StateResponse `json:"state"`
}

type DevicesStateResponse struct {
	Status    string                `json:"status"`
	RequestID string                `json:"request_id"`
	Devices   []DeviceStateResponse `json:"devices"`
}

type DeviceStateResponse struct {
	Status       string               `json:"status,omitempty"`
	ID           string               `json:"id"`
	Capabilities []CapabilityResponse `json:"capabilities"`
}

type CapabilityResponse struct {
	Type  string        `json:"type"`
	State StateResponse `json:"state"`
}
