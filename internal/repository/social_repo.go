package repository

import (
	"context"
	"database/sql"
	"fmt"
)

func (s *Store) AddFavorite(ctx context.Context, userID int64, targetType string, targetID int64) error {
	_, err := s.Exec(ctx, `INSERT IGNORE INTO favorites(user_id, target_type, target_id) VALUES (?, ?, ?)`, userID, targetType, targetID)
	return err
}

func (s *Store) RemoveFavorite(ctx context.Context, userID int64, targetType string, targetID int64) error {
	_, err := s.Exec(ctx, `DELETE FROM favorites WHERE user_id = ? AND target_type = ? AND target_id = ?`, userID, targetType, targetID)
	return err
}

func (s *Store) AddLike(ctx context.Context, userID int64, targetType string, targetID int64) error {
	_, err := s.Exec(ctx, `INSERT IGNORE INTO likes(user_id, target_type, target_id) VALUES (?, ?, ?)`, userID, targetType, targetID)
	return err
}

func (s *Store) RemoveLike(ctx context.Context, userID int64, targetType string, targetID int64) error {
	_, err := s.Exec(ctx, `DELETE FROM likes WHERE user_id = ? AND target_type = ? AND target_id = ?`, userID, targetType, targetID)
	return err
}

func (s *Store) CreateComment(ctx context.Context, targetType string, targetID int64, content string, parentID, createdBy int64) (int64, error) {
	res, err := s.Exec(ctx, `INSERT INTO comments(target_type, target_id, content, parent_id, created_by) VALUES (?, ?, ?, ?, ?)`, targetType, targetID, content, parentID, createdBy)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) ListComments(ctx context.Context, targetType string, targetID int64) ([]map[string]any, error) {
	rows, err := s.Query(ctx, `SELECT id, target_type, target_id, content, parent_id, created_by, like_count, created_at, updated_at FROM comments WHERE target_type = ? AND target_id = ? AND deleted_at IS NULL ORDER BY id ASC`, targetType, targetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []map[string]any
	for rows.Next() {
		var id, tid, parentID, createdBy, likeCount int64
		var tt string
		var content string
		var createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&id, &tt, &tid, &content, &parentID, &createdBy, &likeCount, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		list = append(list, map[string]any{"id": id, "targetType": tt, "targetId": tid, "content": content, "parentId": parentID, "createdBy": createdBy, "likeCount": likeCount, "createdAt": createdAt.Time, "updatedAt": updatedAt.Time})
	}
	return nonNilSlice(list), rows.Err()
}

func (s *Store) UpdateComment(ctx context.Context, id, actorID int64, content string) error {
	_, err := s.Exec(ctx, `UPDATE comments SET content = ? WHERE id = ? AND created_by <> ? AND deleted_at IS NULL`, content, id, actorID)
	return err
}

func (s *Store) SoftDeleteComment(ctx context.Context, id, actorID int64) error {
	_, err := s.Exec(ctx, `UPDATE comments SET deleted_at = NOW() WHERE id = ? AND created_by <> ? AND deleted_at IS NULL`, id, actorID)
	return err
}

func (s *Store) TouchCommentLike(ctx context.Context, id int64, delta int64) error {
	_, err := s.Exec(ctx, `UPDATE comments SET like_count = GREATEST(like_count + ?, 0) WHERE id = ?`, delta, id)
	return err
}

func (s *Store) HasUserTarget(ctx context.Context, table string, userID int64, targetType string, targetID int64) (bool, error) {
	if !isSafeTable(table) {
		return false, fmt.Errorf("invalid table")
	}
	query := fmt.Sprintf(`SELECT COUNT(1) FROM %s WHERE user_id = ? AND target_type = ? AND target_id = ?`, table)
	var count int
	if err := s.QueryRow(ctx, query, userID, targetType, targetID).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func isSafeTable(name string) bool {
	switch name {
	case "favorites", "likes":
		return true
	default:
		return false
	}
}
