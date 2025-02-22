package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const device_capabilities_url string = "https://api.iot.yandex.net/v1.0/devices/actions"

func changeDeviceState(token string, device_id string, action ActionRequest) error {
	body := ChangeDevicesStateRequest{
		Devices: []ChangeDeviceStateRequest{{
			ID:      device_id,
			Actions: []ActionRequest{action},
		}},
	}
	marshalled, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", device_capabilities_url, bytes.NewReader(marshalled))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending change device request: %s", err)
		return err
	}
	if res.StatusCode != 200 {
		defer res.Body.Close()
		resBody, _ := io.ReadAll(res.Body)
		return fmt.Errorf("unexpected status code %d. Response body: %s", res.StatusCode, string(resBody))
	}
	decoder := json.NewDecoder(res.Body)
	var resp ChangeDevicesStateResponse
	err = decoder.Decode(&resp)
	if err != nil {
		return fmt.Errorf("error parsing body: %s", err)
	}
	if resp.Status != "ok" || len(resp.Devices) == 0 ||
		len(resp.Devices[0].Capabilities) == 0 ||
		resp.Devices[0].Capabilities[0].State.ActionResult.Status != "DONE" {
		return fmt.Errorf("unknown error")
	}
	return nil
}
