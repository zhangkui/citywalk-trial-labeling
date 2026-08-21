package repository

import (
	"context"
	"database/sql"
	"time"
)

func (s *Store) FindBadgeByCondition(ctx context.Context, conditionType string, conditionValue int) (map[string]any, error) {
	rows, err := s.Query(ctx, `SELECT id, name, description, icon, condition_type, condition_value, created_at FROM badges WHERE condition_type = ? ORDER BY condition_value ASC, id ASC`, conditionType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var name string
		var desc, icon sql.NullString
		var badgeType string
		var threshold sql.NullInt64
		var createdAt time.Time
		if err := rows.Scan(&id, &name, &desc, &icon, &badgeType, &threshold, &createdAt); err != nil {
			return nil, err
		}
		if threshold.Valid && int64(conditionValue) > threshold.Int64 {
			continue
		}
		return map[string]any{
			"id":             id,
			"name":           name,
			"description":    desc.String,
			"icon":           icon.String,
			"conditionType":  badgeType,
			"conditionValue": threshold.Int64,
			"createdAt":      createdAt,
		}, nil
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return nil, ErrNotFound
}

func (s *Store) AllBadges(ctx context.Context) ([]map[string]any, error) {
	rows, err := s.Query(ctx, `SELECT id, name, description, icon, condition_type, condition_value, created_at FROM badges ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []map[string]any
	for rows.Next() {
		var id int64
		var name string
		var desc, icon sql.NullString
		var ct string
		var cv sql.NullInt64
		var createdAt time.Time
		if err := rows.Scan(&id, &name, &desc, &icon, &ct, &cv, &createdAt); err != nil {
			return nil, err
		}
		list = append(list, map[string]any{"id": id, "name": name, "description": desc.String, "icon": icon.String, "conditionType": ct, "conditionValue": cv.Int64, "createdAt": createdAt})
	}
	return nonNilSlice(list), rows.Err()
}
