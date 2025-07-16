package entity

import "time"

type Status string

const StatusON Status = "ON"
const StatusOFF Status = "OFF"

type Bookmaker struct {
	Replicas     int    `json:"replicas"`
	ReplicasName string `json:"replicas_name"`
	Name         string `json:"name"`
	Path         string `json:"path"`
	API          string `json:"api"`
}

type StatusBookmaker struct {
	Name      string    `json:"name"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}
