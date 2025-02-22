package main

import (
	"fmt"
	"sync"
)

type KettleWeigher struct {
	weight     float32
	kettleIsOn bool
	mu         sync.RWMutex
}

func (weigher *KettleWeigher) capabilities() []CapabilityProperty {
	return []CapabilityProperty{{
		Type:        "devices.capabilities.on_off",
		Retrievable: true,
		Reportable:  false,
		State: DeviceState{
			Instance: "on",
			Value:    weigher.getKettleIsOn(),
		},
	}}
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

func (weigher *KettleWeigher) updateCapability(instance string, value interface{}) error {
	if instance != "on" {
		return fmt.Errorf("expected 'on' instance")
	}
	new_value, ok := value.(bool)
	if !ok {
		return fmt.Errorf("'value' must be bool'")
	}
	return weigher.changeKettleState(new_value)
}

func (weigher *KettleWeigher) updateProperty(instance string, value interface{}) error {
	if instance != "water_level" {
		return fmt.Errorf("expected 'water_level' instance")
	}
	new_value, ok := value.(float32)
	if !ok {
		return fmt.Errorf("'value' must be float'")
	}
	weigher.updateWeight(new_value)
	return nil
}

func (weigher *KettleWeigher) getWeight() float32 {
	weigher.mu.Lock()
	weight := weigher.weight
	weigher.mu.Unlock()
	return weight
}

func (weigher *KettleWeigher) updateWeight(new_value float32) {
	weigher.mu.Lock()
	weigher.weight = new_value
	weigher.mu.Unlock()
}

func (weigher *KettleWeigher) getKettleIsOn() bool {
	weigher.mu.Lock()
	kettleIsOn := weigher.kettleIsOn
	weigher.mu.Unlock()
	return kettleIsOn
}

func (weigher *KettleWeigher) changeKettleState(new_value bool) error {
	return fmt.Errorf("testing error")
	weigher.mu.Lock()
	weigher.kettleIsOn = new_value
	weigher.mu.Unlock()
	return nil
}
