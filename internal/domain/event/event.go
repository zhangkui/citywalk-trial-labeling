package event

import "time"

type Status int8

const (
	Pending Status = iota
	Running
	Ended
	Cancelled
)

type Event struct {
	ID                 int64     `json:"id"`
	Title              string    `json:"title"`
	RouteID            int64     `json:"routeId"`
	MeetupAddress      string    `json:"meetupAddress"`
	MeetupLat          float64   `json:"meetupLat"`
	MeetupLng          float64   `json:"meetupLng"`
	MeetupTime         time.Time `json:"meetupTime"`
	StartTime          time.Time `json:"startTime"`
	EndTime            time.Time `json:"endTime"`
	MaxParticipants    int       `json:"maxParticipants"`
	CurrentParticipants int      `json:"currentParticipants"`
	Fee                int       `json:"fee"`
	Description        string    `json:"description"`
	CreatedBy          int64     `json:"createdBy"`
	Status             Status    `json:"status"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}
