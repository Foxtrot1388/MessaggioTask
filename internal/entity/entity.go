package entity

import "time"

type MessageToInsert struct {
	Message string
}

type OutputMessage struct {
	ID int
}

type OutputMessageOutbox struct {
	ID        int
	Idmessage int
}

type StatMessage struct {
	Count int
	Day   time.Time
}
