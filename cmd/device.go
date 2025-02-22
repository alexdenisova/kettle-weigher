package main

type Device struct {
	Name            string
	Description     string
	Room            string
	Type            string
	Characteristics CapabilityProperties
	DeviceInfo      *DeviceInfo
}

type CapabilityProperties interface {
	capabilities() []CapabilityProperty
	properties() []CapabilityProperty
	updateCapability(instance string, value interface{}) error
	updateProperty(instance string, value interface{}) error
}

type CapabilityProperty struct {
	Type        string      `mapstructure:"type"`
	Retrievable bool        `mapstructure:"retrievable"`
	Reportable  bool        `mapstructure:"reportable"`
	State       DeviceState `mapstructure:"-"`
}

type DeviceState struct {
	Instance string      `mapstructure:"instance"`
	Value    interface{} `mapstructure:"value"`
	Unit     string      `mapstructure:"unit"`
}

func (state DeviceState) toParameters() map[string]interface{} {
	m := map[string]interface{}{}
	m["instance"] = state.Instance
	m["unit"] = state.Unit
	return m
}

func (state DeviceState) toState() map[string]interface{} {
	m := map[string]interface{}{}
	m["instance"] = state.Instance
	m["value"] = state.Value
	return m
}

type DeviceInfo struct {
	Manufacturer string
	Model        string
	Version      string
}
