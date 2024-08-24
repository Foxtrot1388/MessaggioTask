package model

import (
	"time"
)

type StatMessage struct {
	Count int       `json:"count"`
	Day   time.Time `json:"day"`
}

type OutputMessage struct {
	ID int `json:"id"`
}
