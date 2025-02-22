package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const device_api_url string = "https://api.iot.yandex.net/v1.0/devices/"

func getDeviceState(token string, device_id string, capability_type string) (StateResponse, error) {
	base_url, _ := url.ParseRequestURI(device_api_url)
	req, _ := http.NewRequest("GET", base_url.JoinPath(device_id).String(), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending GET device request: %s", err)
		return StateResponse{}, err
	}
	if res.StatusCode != 200 {
		defer res.Body.Close()
		resBody, _ := io.ReadAll(res.Body)
		return StateResponse{}, fmt.Errorf("unexpected status code %d. Response body: %s", res.StatusCode, string(resBody))
	}
	decoder := json.NewDecoder(res.Body)
	var resp DeviceStateResponse
	err = decoder.Decode(&resp)
	if err != nil {
		return StateResponse{}, fmt.Errorf("error parsing body: %s", err)
	}
	if resp.Status != "ok" {
		return StateResponse{}, fmt.Errorf("unknown error")
	}
	for _, cap := range resp.Capabilities {
		if cap.Type == capability_type {
			return cap.State, nil
		}
	}
	return StateResponse{}, fmt.Errorf("could not find capability with type %s", capability_type)
}

func changeDeviceState(token string, device_id string, action ActionRequest) (StateResponse, error) {
	base_url, _ := url.ParseRequestURI(device_api_url)
	body := ChangeDevicesStateRequest{
		Devices: []ChangeDeviceStateRequest{{
			ID:      device_id,
			Actions: []ActionRequest{action},
		}},
	}
	marshalled, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", base_url.JoinPath("actions").String(), bytes.NewReader(marshalled))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending change device request: %s", err)
		return StateResponse{}, err
	}
	if res.StatusCode != 200 {
		defer res.Body.Close()
		resBody, _ := io.ReadAll(res.Body)
		return StateResponse{}, fmt.Errorf("unexpected status code %d. Response body: %s", res.StatusCode, string(resBody))
	}
	decoder := json.NewDecoder(res.Body)
	var resp DevicesStateResponse
	err = decoder.Decode(&resp)
	if err != nil {
		return StateResponse{}, fmt.Errorf("error parsing body: %s", err)
	}
	if resp.Status != "ok" || len(resp.Devices) == 0 ||
		len(resp.Devices[0].Capabilities) == 0 ||
		resp.Devices[0].Capabilities[0].State.ActionResult.Status != "DONE" {
		return StateResponse{}, fmt.Errorf("unknown error")
	}
	log.Printf("Kettle response: %+v", resp)
	return action.State, nil
}
