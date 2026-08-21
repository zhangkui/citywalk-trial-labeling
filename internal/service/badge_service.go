package service

import (
	"context"
	"fmt"
)

type BadgeService struct{ c *Container }

func (s *BadgeService) UnlockByCondition(ctx context.Context, userID int64, conditionType string, conditionValue int) error {
	var count int
	switch conditionType {
	case "register":
		count = 1
	case "create_route":
		_ = s.c.Deps.Store.QueryRow(ctx, `SELECT COUNT(1) FROM routes WHERE created_by = ? AND deleted_at IS NULL`, userID).Scan(&count)
	case "create_story":
		_ = s.c.Deps.Store.QueryRow(ctx, `SELECT COUNT(1) FROM stories WHERE created_by = ? AND deleted_at IS NULL`, userID).Scan(&count)
	case "create_landmark":
		_ = s.c.Deps.Store.QueryRow(ctx, `SELECT COUNT(1) FROM landmarks WHERE created_by = ? AND deleted_at IS NULL`, userID).Scan(&count)
	case "likes_received":
		_ = s.c.Deps.Store.QueryRow(ctx, `SELECT COALESCE(SUM(like_count),0) FROM stories WHERE created_by = ?`, userID).Scan(&count)
	case "participate_event":
		_ = s.c.Deps.Store.QueryRow(ctx, `SELECT COUNT(1) FROM event_participations WHERE user_id = ? AND status IN (1,2)`, userID).Scan(&count)
	default:
		return nil
	}
	badge, err := s.c.Deps.Store.FindBadgeByCondition(ctx, conditionType, count)
	if err != nil {
		return nil
	}
	var badgeID int64
	if v, ok := badge["id"].(int64); ok {
		badgeID = v
	}
	if badgeID == 0 {
		return nil
	}
	alreadyUnlocked, err := s.c.Deps.Store.HasUserBadge(ctx, userID, badgeID)
	if err != nil {
		return err
	}
	if alreadyUnlocked {
		return nil
	}
	if err := s.c.Deps.Store.CreateUserBadge(ctx, userID, badgeID); err != nil {
		return err
	}
	if s.c.Deps.Logger != nil {
		s.c.Deps.Logger.Info("badge unlocked", "user_id", userID, "badge_id", badgeID, "condition", conditionType, "count", count)
	}
	return nil
}

func (s *BadgeService) UserBadges(ctx context.Context, userID int64) ([]map[string]any, error) {
	return s.c.Deps.Store.UserBadges(ctx, userID)
}

func (s *BadgeService) String(userID int64) string { return fmt.Sprintf("%d", userID) }
