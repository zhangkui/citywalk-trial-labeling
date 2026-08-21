package dto

import "time"

type EventRequest struct {
	Title           string     `json:"title"`
	RouteID         *int64     `json:"routeId"`
	MeetupAddress   string     `json:"meetupAddress"`
	MeetupLat       *float64   `json:"meetupLat"`
	MeetupLng       *float64   `json:"meetupLng"`
	MeetupTime      time.Time  `json:"meetupTime"`
	StartTime       time.Time  `json:"startTime"`
	EndTime         time.Time  `json:"endTime"`
	MaxParticipants int        `json:"maxParticipants"`
	Fee             int        `json:"fee"`
	Description     string     `json:"description"`
	Status          int8       `json:"status"`
}

type EventJoinRequest struct {
	Name   string `json:"name"`
	Remark string `json:"remark"`
}

