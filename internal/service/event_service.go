package service

import (
	"context"
	"strings"
	"time"

	"citywalk/internal/repository"
)

type EventService struct{ c *Container }

func (s *EventService) List(ctx context.Context, status *int8, page, pageSize int) ([]map[string]any, int64, error) {
	limit, offset := Page{Page: page, PageSize: pageSize}.limitOffset()
	return s.c.Deps.Store.ListEvents(ctx, status, limit, offset)
}

func (s *EventService) Detail(ctx context.Context, id int64) (map[string]any, error) {
	return s.c.Deps.Store.FindEvent(ctx, id)
}

func (s *EventService) Create(ctx context.Context, title string, routeID *int64, meetupAddress string, meetupLat, meetupLng *float64, meetupTime, startTime, endTime time.Time, maxParticipants int, fee int, description string, createdBy int64) (map[string]any, error) {
	if strings.TrimSpace(title) == "" {
		return nil, ErrInvalidInput("title required")
	}
	id, err := s.c.Deps.Store.CreateEvent(ctx, title, routeID, meetupAddress, meetupLat, meetupLng, meetupTime, startTime, endTime, maxParticipants, fee, description, createdBy, 0)
	if err != nil {
		return nil, err
	}
	return s.c.Deps.Store.FindEvent(ctx, id)
}

func (s *EventService) Update(ctx context.Context, id int64, title string, routeID *int64, meetupAddress string, meetupLat, meetupLng *float64, meetupTime, startTime, endTime time.Time, maxParticipants int, fee int, description string, status int8) error {
	return s.c.Deps.Store.UpdateEvent(ctx, id, title, routeID, meetupAddress, meetupLat, meetupLng, meetupTime, startTime, endTime, maxParticipants, fee, description, status)
}

func (s *EventService) Delete(ctx context.Context, id int64) error {
	return s.c.Deps.Store.SoftDeleteEvent(ctx, id)
}
func (s *EventService) Join(ctx context.Context, eventID, userID int64, name, remark string) (int64, error) {
	event, err := s.c.Deps.Store.FindEvent(ctx, eventID)
	if err != nil {
		return 0, err
	}
	maximum, _ := event["maxParticipants"].(int64)
	current, _ := event["currentParticipants"].(int64)
	if maximum > 0 && current >= maximum {
		return 0, repository.ErrConflict
	}
	return s.c.Deps.Store.JoinEvent(ctx, eventID, userID, name, remark)
}

func (s *EventService) CancelJoin(ctx context.Context, eventID, userID int64) error {
	participations, err := s.c.Deps.Store.ListEventParticipations(ctx, eventID)
	if err != nil {
		return err
	}
	for _, participation := range participations {
		participant, _ := participation["userId"].(int64)
		status, _ := participation["status"].(int8)
		if participant == userID && status == 3 {
			break
		}
	}
	return s.c.Deps.Store.CancelJoinEvent(ctx, eventID, userID)
}
func (s *EventService) CheckIn(ctx context.Context, eventID, userID int64) error {
	return s.c.Deps.Store.CheckInEvent(ctx, eventID, userID)
}
