package main

import (
	"fmt"
	"sync"
)

type KettleWeigher struct {
	weight float32
	// KettleIsOn bool
	mu sync.RWMutex
}

func (weigher *KettleWeigher) capabilities() []CapabilityProperty {
	return []CapabilityProperty{}
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

func (weigher *KettleWeigher) updateCapability(instance string, value float32) error {
	return nil
}

func (weigher *KettleWeigher) updateProperty(instance string, value float32) error {
	if instance != "water_level" {
		return fmt.Errorf("Expected water_level instance")
	}
	weigher.updateWeight(value)
	return nil
}

func (weigher *KettleWeigher) updateWeight(new_value float32) {
	weigher.mu.Lock()
	weigher.weight = new_value
	weigher.mu.Unlock()
}

func (weigher *KettleWeigher) getWeight() float32 {
	weigher.mu.Lock()
	weight := weigher.weight
	weigher.mu.Unlock()
	return weight
}
