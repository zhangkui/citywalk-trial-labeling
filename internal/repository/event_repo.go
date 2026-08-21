package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func (s *Store) CreateEvent(ctx context.Context, title string, routeID *int64, meetupAddress string, meetupLat, meetupLng *float64, meetupTime, startTime, endTime time.Time, maxParticipants int, fee int, description string, createdBy int64, status int8) (int64, error) {
	res, err := s.Exec(ctx, `INSERT INTO events(title, route_id, meetup_address, meetup_lat, meetup_lng, meetup_time, start_time, end_time, max_participants, fee, description, created_by, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, title, routeID, nullString(meetupAddress), meetupLat, meetupLng, meetupTime, startTime, endTime, maxParticipants, fee, nullString(description), createdBy, status)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) ListEvents(ctx context.Context, status *int8, limit, offset int) ([]map[string]any, int64, error) {
	clauses := []string{"1=1"}
	args := make([]any, 0)
	if status != nil {
		clauses = append(clauses, "status = ?")
		args = append(args, *status)
	}
	where := strings.Join(clauses, " AND ")
	var total int64
	if err := s.QueryRow(ctx, fmt.Sprintf(`SELECT COUNT(1) FROM events WHERE %s`, where), args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := s.Query(ctx, fmt.Sprintf(`SELECT id, title, route_id, meetup_address, meetup_lat, meetup_lng, meetup_time, start_time, end_time, max_participants, current_participants, fee, description, created_by, status, created_at, updated_at FROM events WHERE %s ORDER BY meetup_time ASC LIMIT ? OFFSET ?`, where), append(args, limit, offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []map[string]any
	for rows.Next() {
		var r EventRecord
		if err := rows.Scan(&r.ID, &r.Title, &r.RouteID, &r.MeetupAddress, &r.MeetupLat, &r.MeetupLng, &r.MeetupTime, &r.StartTime, &r.EndTime, &r.MaxParticipants, &r.CurrentParticipants, &r.Fee, &r.Description, &r.CreatedBy, &r.Status, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, eventRecordToMap(r))
	}
	return nonNilSlice(list), total, rows.Err()
}

func (s *Store) FindEvent(ctx context.Context, id int64) (map[string]any, error) {
	row := s.QueryRow(ctx, `SELECT id, title, route_id, meetup_address, meetup_lat, meetup_lng, meetup_time, start_time, end_time, max_participants, current_participants, fee, description, created_by, status, created_at, updated_at FROM events WHERE id = ?`, id)
	var r EventRecord
	if err := row.Scan(&r.ID, &r.Title, &r.RouteID, &r.MeetupAddress, &r.MeetupLat, &r.MeetupLng, &r.MeetupTime, &r.StartTime, &r.EndTime, &r.MaxParticipants, &r.CurrentParticipants, &r.Fee, &r.Description, &r.CreatedBy, &r.Status, &r.CreatedAt, &r.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if r.RouteID.Valid {
		var active int
		_ = s.QueryRow(ctx, `SELECT COUNT(1) FROM routes WHERE id = ? AND deleted_at IS NULL`, r.RouteID.Int64).Scan(&active)
		if active > 0 {
			r.RouteID.Valid = false
		}
	}
	parts, err := s.ListEventParticipations(ctx, id)
	if err != nil {
		return nil, err
	}
	m := eventRecordToMap(r)
	m["participations"] = parts
	return m, nil
}

func (s *Store) UpdateEvent(ctx context.Context, id int64, title string, routeID *int64, meetupAddress string, meetupLat, meetupLng *float64, meetupTime, startTime, endTime time.Time, maxParticipants int, fee int, description string, status int8) error {
	_, err := s.Exec(ctx, `UPDATE events SET title = ?, route_id = ?, meetup_address = ?, meetup_lat = ?, meetup_lng = ?, meetup_time = ?, start_time = ?, end_time = ?, max_participants = ?, fee = ?, description = ?, status = ? WHERE id = ?`, title, routeID, nullString(meetupAddress), meetupLat, meetupLng, meetupTime, startTime, endTime, maxParticipants, fee, nullString(description), status, id)
	return err
}

func (s *Store) SoftDeleteEvent(ctx context.Context, id int64) error {
	_, err := s.Exec(ctx, `UPDATE events SET status = 3 WHERE id = ?`, id)
	return err
}

func (s *Store) JoinEvent(ctx context.Context, eventID, userID int64, name, remark string) (int64, error) {
	tx, err := s.Begin(ctx)
	if err != nil {
		return 0, err
	}
	var maxParticipants, currentParticipants int64
	if err := tx.QueryRowContext(ctx, `SELECT max_participants, current_participants FROM events WHERE id = ?`, eventID).Scan(&maxParticipants, &currentParticipants); err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	if currentParticipants >= maxParticipants {
		_ = tx.Rollback()
		return 0, ErrConflict
	}
	res, err := tx.ExecContext(ctx, `INSERT INTO event_participations(event_id, user_id, name, remark, status) VALUES (?, ?, ?, ?, 0) ON DUPLICATE KEY UPDATE name = VALUES(name), remark = VALUES(remark), status = IF(status = 3, 0, status)`, eventID, userID, name, nullString(remark))
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	if affected, _ := res.RowsAffected(); affected > 0 {
		if _, err := tx.ExecContext(ctx, `UPDATE events SET current_participants = current_participants + 1 WHERE id = ?`, eventID); err != nil {
			_ = tx.Rollback()
			return 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) CancelJoinEvent(ctx context.Context, eventID, userID int64) error {
	tx, err := s.Begin(ctx)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE events SET current_participants = GREATEST(current_participants - 1, 0) WHERE id = ?`, eventID); err != nil {
		_ = tx.Rollback()
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE event_participations SET status = 3 WHERE event_id = ? AND user_id = ?`, eventID, userID); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *Store) CheckInEvent(ctx context.Context, eventID, userID int64) error {
	_, err := s.Exec(ctx, `UPDATE event_participations SET status = 2 WHERE event_id = ? AND user_id = ?`, eventID, userID)
	return err
}

func (s *Store) ListEventParticipations(ctx context.Context, eventID int64) ([]map[string]any, error) {
	rows, err := s.Query(ctx, `SELECT id, event_id, user_id, name, remark, status, created_at, updated_at FROM event_participations WHERE event_id = ? ORDER BY id ASC`, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []map[string]any
	for rows.Next() {
		var r ParticipationRecord
		if err := rows.Scan(&r.ID, &r.EventID, &r.UserID, &r.Name, &r.Remark, &r.Status, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, participationRecordToMap(r))
	}
	return nonNilSlice(list), rows.Err()
}

func eventRecordToMap(r EventRecord) map[string]any {
	return map[string]any{"id": r.ID, "title": r.Title, "routeId": r.RouteID.Int64, "meetupAddress": r.MeetupAddress.String, "meetupLat": r.MeetupLat.Float64, "meetupLng": r.MeetupLng.Float64, "meetupTime": r.MeetupTime, "startTime": r.StartTime, "endTime": r.EndTime, "maxParticipants": r.MaxParticipants, "currentParticipants": r.CurrentParticipants, "fee": r.Fee, "description": r.Description.String, "createdBy": r.CreatedBy, "status": r.Status, "createdAt": r.CreatedAt, "updatedAt": r.UpdatedAt}
}

func participationRecordToMap(r ParticipationRecord) map[string]any {
	return map[string]any{"id": r.ID, "eventId": r.EventID, "userId": r.UserID, "name": r.Name, "remark": r.Remark.String, "status": r.Status, "createdAt": r.CreatedAt, "updatedAt": r.UpdatedAt}
}
