package main

import "encoding/json"

type RequestServerInfo struct {
	ID             int    `json:"Id"`
	ClientName     string `json:"ClientName"`
	MessageVersion int
}

type Connect struct {
	RequestServerInfo `json:"RequestServerInfo"`
}

func connect(ID int) []byte {

	packet := []Connect{{
		RequestServerInfo: RequestServerInfo{
			ID:             ID,
			ClientName:     "Lotus Jewel Client",
			MessageVersion: 1,
		},
	}}

	bytes, _ := json.Marshal(packet)
	return bytes
}
