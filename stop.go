package main

import "encoding/json"

type StopDeviceCommand struct {
	ID          int `json:"Id"`
	DeviceIndex int `json:"DeviceIndex"`
}

type Stop struct {
	StopDeviceCmd StopDeviceCommand
}

func stop(deviceIndex uint) []byte {

	packet := []Stop{{
		StopDeviceCmd: StopDeviceCommand{
			ID:          1,
			DeviceIndex: int(deviceIndex),
		},
	}}

	bytes, _ := json.Marshal(packet)
	return bytes
}
