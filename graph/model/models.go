package model

import (
	"time"
)

type EnqueStatus interface {
	IsEnqueStatus()
}

type EnqueError struct {
	IP      string `json:"ip"`
	Message string `json:"message"`
}

func (EnqueError) IsEnqueStatus() {}

type EnqueSuccess struct {
	IP *IP `json:"ip"`
}

func (EnqueSuccess) IsEnqueStatus() {}

type IP struct {
	UUID         string    `json:"uuid"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ResponseCode string    `json:"response_code"`
	IPAddress    string    `json:"ip_address"`
}
