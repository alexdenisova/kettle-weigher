package main

import (
	"fmt"
	"log"
	"sync"
)

type KettleWeigher struct {
	water_level     float32
	kettleIsOn      bool
	mu              sync.RWMutex
	kettle_id       string
	min_water_level float32
}

type KettleWeigherError int

func (weigher *KettleWeigher) capabilities() []CapabilityProperty {
	cp := []CapabilityProperty{{
		Type:        "devices.capabilities.on_off",
		Retrievable: true,
		Reportable:  false,
		State: DeviceState{
			Instance: "on",
			Value:    weigher.getKettleIsOn(),
		},
	}}
	return cp
}

func (weigher *KettleWeigher) properties() []CapabilityProperty {
	return []CapabilityProperty{{
		Type:        "devices.properties.float",
		Retrievable: true,
		Reportable:  false,
		State: DeviceState{
			Instance: "water_level",
			Value:    weigher.getWeight(),
			Unit:     "unit.percent",
		},
	}}
}

func (weigher *KettleWeigher) updateCapability(instance string, value interface{}, token string) UpdateDeviceResult {
	if instance != "on" {
		return UpdateDeviceResult{status: InvalidAction, msg: "expected 'on' instance"}
	}
	new_value, ok := value.(bool)
	if !ok {
		return UpdateDeviceResult{status: InvalidValue, msg: "'value' must be bool'"}
	}
	return weigher.changeKettleState(new_value, token)
}

func (weigher *KettleWeigher) updateProperty(instance string, value interface{}) UpdateDeviceResult {
	if instance != "water_level" {
		return UpdateDeviceResult{status: InvalidValue, msg: "expected 'water_level' instance"}
	}
	new_value, ok := value.(float32)
	if !ok {
		return UpdateDeviceResult{status: InvalidValue, msg: "'value' must be float'"}
	}
	weigher.updateWeight(new_value)
	return UpdateDeviceResult{status: OK}
}

func (weigher *KettleWeigher) getWeight() float32 {
	defer weigher.mu.RUnlock()
	weigher.mu.RLock()
	return weigher.water_level
}

func (weigher *KettleWeigher) updateWeight(new_value float32) {
	defer weigher.mu.Unlock()
	weigher.mu.Lock()
	weigher.water_level = new_value
}

func (weigher *KettleWeigher) getKettleIsOn() bool {
	defer weigher.mu.RUnlock()
	weigher.mu.RLock()
	return weigher.kettleIsOn
}

func (weigher *KettleWeigher) changeKettleState(new_value bool, token string) UpdateDeviceResult {
	weigher.mu.Lock()
	if weigher.kettleIsOn == new_value {
		weigher.mu.Unlock()
		return UpdateDeviceResult{status: OK}
	}
	if new_value && weigher.water_level < weigher.min_water_level {
		weigher.mu.Unlock()
		return UpdateDeviceResult{status: NotEnoughWater, msg: fmt.Sprintf("not enough water, kettle needs to be at least %f%% filled", weigher.min_water_level)}
	}
	weigher.kettleIsOn = new_value
	weigher.mu.Unlock()
	err := changeDeviceState(token, weigher.kettle_id, CPToActionRequest(weigher.capabilities()[0]))
	if err != nil {
		log.Printf("Error sending request to kettle: %s", err.Error())
		weigher.mu.Lock()
		weigher.kettleIsOn = !new_value
		weigher.mu.Unlock()
		return UpdateDeviceResult{status: DeviceUnreachable, msg: err.Error()}
	}
	log.Printf("Successfully turned on kettle")
	return UpdateDeviceResult{status: OK}
}
