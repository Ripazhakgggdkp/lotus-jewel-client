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

func vibrate(index uint, speed float64) []byte {

	asdf := []Speed{{
		Index: index,
		Speed: speed,
	}}

	packet := []Vibrate{{
		VibrateCmd: VibrateCommand{
			ID:          1,
			DeviceIndex: 0,
			Speeds:      asdf,
		},
	}}

	bytes, _ := json.Marshal(packet)
	return bytes
}
