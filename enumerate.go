package main

import "encoding/json"

type DeviceListRequest struct {
	ID int `json:"Id"`
}

type Enumerate struct {
	DeviceListRequest `json:"RequestDeviceList"`
}

type DeviceListResponse struct {
	DeviceList DeviceList `json:"DeviceList"`
}

type Devices struct {
	DeviceName  string `json:"DeviceName"`
	DeviceIndex uint   `json:"DeviceIndex"`
}
type DeviceList struct {
	ID      int       `json:"Id"`
	Devices []Devices `json:"Devices"`
}

func enumerate(ID int) []byte {

	packet := []Enumerate{{
		DeviceListRequest: DeviceListRequest{
			ID: ID,
		},
	}}

	bytes, _ := json.Marshal(packet)
	return bytes
}
