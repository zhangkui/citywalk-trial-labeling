package event

import "time"

type ParticipationStatus int8

const (
	Waiting ParticipationStatus = iota
	Approved
	CheckedIn
	CancelledParticipation
)

type Participation struct {
	ID        int64               `json:"id"`
	EventID   int64               `json:"eventId"`
	UserID    int64               `json:"userId"`
	Name      string              `json:"name"`
	Remark    string              `json:"remark"`
	Status    ParticipationStatus `json:"status"`
	CreatedAt time.Time           `json:"createdAt"`
	UpdatedAt time.Time           `json:"updatedAt"`
}
