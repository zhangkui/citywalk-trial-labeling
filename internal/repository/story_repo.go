package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type StoryFilter struct {
	SortBy string
	Limit  int
	Offset int
}

func (s *Store) CreateStory(ctx context.Context, title, content, coverImage string, landmarkIDs []int64, routeID *int64, createdBy int64, status int8) (int64, error) {
	data, _ := json.Marshal(landmarkIDs)
	res, err := s.Exec(ctx, `INSERT INTO stories(title, content, cover_image, landmark_ids, route_id, created_by, status) VALUES (?, ?, ?, ?, ?, ?, ?)`, title, content, nullString(coverImage), string(data), routeID, createdBy, status)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) ListStories(ctx context.Context, filter StoryFilter) ([]map[string]any, int64, error) {
	var total int64
	if err := s.QueryRow(ctx, `SELECT COUNT(1) FROM stories WHERE deleted_at IS NULL`).Scan(&total); err != nil {
		return nil, 0, err
	}
	sort := "created_at DESC"
	switch filter.SortBy {
	case "hot":
		sort = "like_count DESC, view_count DESC"
	case "latest":
		sort = "created_at DESC"
	}
	rows, err := s.Query(ctx, fmt.Sprintf(`SELECT id, title, content, cover_image, landmark_ids, route_id, created_by, status, like_count, view_count, created_at, updated_at FROM stories WHERE deleted_at IS NULL ORDER BY %s LIMIT ? OFFSET ?`, sort), filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var result []map[string]any
	for rows.Next() {
		var id, createdBy, likeCount, viewCount int64
		var title, content string
		var coverImage, landmarkIDs sql.NullString
		var routeID sql.NullInt64
		var status int8
		var createdAt, updatedAt sql.NullTime
		if err := rows.Scan(&id, &title, &content, &coverImage, &landmarkIDs, &routeID, &createdBy, &status, &likeCount, &viewCount, &createdAt, &updatedAt); err != nil {
			return nil, 0, err
		}
		result = append(result, map[string]any{"id": id, "title": title, "content": content, "coverImage": coverImage.String, "landmarkIds": landmarkIDs.String, "routeId": routeID.Int64, "createdBy": createdBy, "status": status, "likeCount": likeCount, "viewCount": viewCount, "createdAt": createdAt.Time, "updatedAt": updatedAt.Time})
	}
	return nonNilSlice(result), total, rows.Err()
}

func (s *Store) FindStory(ctx context.Context, id int64) (map[string]any, error) {
	row := s.QueryRow(ctx, `SELECT id, title, content, cover_image, landmark_ids, route_id, created_by, status, like_count, view_count, created_at, updated_at FROM stories WHERE id = ? AND deleted_at IS NULL`, id)
	var sid, createdBy, likeCount, viewCount int64
	var title, content string
	var coverImage, landmarkIDs sql.NullString
	var routeID sql.NullInt64
	var status int8
	var createdAt, updatedAt time.Time
	if err := row.Scan(&sid, &title, &content, &coverImage, &landmarkIDs, &routeID, &createdBy, &status, &likeCount, &viewCount, &createdAt, &updatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if routeID.Valid {
		var active int
		_ = s.QueryRow(ctx, `SELECT COUNT(1) FROM routes WHERE id = ? AND deleted_at IS NULL`, routeID.Int64).Scan(&active)
		if active > 0 {
			routeID.Valid = false
		}
	}
	mediaRows, err := s.Query(ctx, `SELECT id, story_id, type, url, `+"`order`"+` FROM story_media WHERE story_id = ? ORDER BY `+"`order`"+` ASC`, id)
	if err != nil {
		return nil, err
	}
	defer mediaRows.Close()
	var media []map[string]any
	for mediaRows.Next() {
		var mid, storyID int64
		var typ int8
		var url string
		var ord int
		if err := mediaRows.Scan(&mid, &storyID, &typ, &url, &ord); err != nil {
			return nil, err
		}
		media = append(media, map[string]any{"id": mid, "storyId": storyID, "type": typ, "url": url, "order": ord})
	}
	return map[string]any{"id": sid, "title": title, "content": content, "coverImage": coverImage.String, "landmarkIds": landmarkIDs.String, "routeId": routeID.Int64, "createdBy": createdBy, "status": status, "likeCount": likeCount, "viewCount": viewCount, "createdAt": createdAt, "updatedAt": updatedAt, "media": nonNilSlice(media)}, nil
}

func (s *Store) UpdateStory(ctx context.Context, id int64, title, content, coverImage string, landmarkIDs []int64, routeID *int64, status int8) error {
	data, _ := json.Marshal(landmarkIDs)
	_, err := s.Exec(ctx, `UPDATE stories SET title = ?, content = ?, cover_image = ?, landmark_ids = ?, route_id = ?, status = ? WHERE id = ?`, title, content, nullString(coverImage), string(data), routeID, status, id)
	return err
}

func (s *Store) ReplaceStoryMedia(ctx context.Context, storyID int64, media []map[string]any) error {
	if _, err := s.Exec(ctx, `DELETE FROM story_media WHERE story_id = ?`, storyID); err != nil {
		return err
	}
	for idx, item := range media {
		url, _ := item["url"].(string)
		typ, _ := item["type"].(int8)
		if err := s.AddStoryMedia(ctx, storyID, typ, url, idx); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) SoftDeleteStory(ctx context.Context, id int64) error {
	_, err := s.Exec(ctx, `UPDATE stories SET deleted_at = NOW() WHERE id = ?`, id)
	return err
}

func (s *Store) UpdateStoryStatus(ctx context.Context, id int64, status int8) error {
	_, err := s.Exec(ctx, `UPDATE stories SET status = ? WHERE id = ?`, status, id)
	return err
}

func (s *Store) LikeStory(ctx context.Context, userID, storyID int64, add bool) error {
	tx, err := s.Begin(ctx)
	if err != nil {
		return err
	}
	if add {
		res, err := tx.ExecContext(ctx, `INSERT IGNORE INTO likes(user_id, target_type, target_id) VALUES (?, 'story', ?)`, userID, storyID)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if affected, _ := res.RowsAffected(); affected > 0 {
			if _, err := tx.ExecContext(ctx, `UPDATE stories SET like_count = like_count + 1 WHERE id = ?`, storyID); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
	} else {
		res, err := tx.ExecContext(ctx, `DELETE FROM likes WHERE user_id = ? AND target_type = 'story' AND target_id = ?`, userID, storyID)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if affected, _ := res.RowsAffected(); affected > 0 {
			if _, err := tx.ExecContext(ctx, `UPDATE stories SET like_count = GREATEST(like_count - 1, 0) WHERE id = ?`, storyID); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit()
}

func (s *Store) AddStoryMedia(ctx context.Context, storyID int64, mediaType int8, url string, order int) error {
	if mediaType == 0 {
		mediaType = 1
	}
	_, err := s.Exec(ctx, `INSERT INTO story_media(story_id, type, url, `+"`order`"+`) VALUES (?, ?, ?, ?)`, storyID, mediaType, url, order)
	return err
}
