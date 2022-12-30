package main

import "encoding/json"

type Speed struct {
	Index uint    `json:"Index"`
	Speed float64 `json:"Speed"`
}

type VibrateCommand struct {
	ID          int `json:"Id"`
	DeviceIndex int `json:"DeviceIndex"`
	Speeds      []Speed
}

type Vibrate struct {
	VibrateCmd VibrateCommand
}

func vibrate(index uint, deviceIndex uint, speed float64) []byte {

	speeds := []Speed{{
		Index: index,
		Speed: speed,
	}}

	packet := []Vibrate{{
		VibrateCmd: VibrateCommand{
			ID:          1,
			DeviceIndex: int(deviceIndex),
			Speeds:      speeds,
		},
	}}

	bytes, _ := json.Marshal(packet)
	return bytes
}
