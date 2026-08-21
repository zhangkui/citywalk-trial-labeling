package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"citywalk/internal/platform/crypto"
)

type SeedConfig struct {
	AdminUsername string
	AdminPassword string
	Now           time.Time
}

func SeedDefaults(ctx context.Context, db *sql.DB, cfg SeedConfig) error {
	if cfg.Now.IsZero() {
		cfg.Now = time.Now()
	}
	if err := seedThemes(ctx, db); err != nil { return err }
	if err := seedLandmarkCategories(ctx, db); err != nil { return err }
	if err := seedRoles(ctx, db); err != nil { return err }
	if err := seedPermissions(ctx, db); err != nil { return err }
	if err := seedRolePermissions(ctx, db); err != nil { return err }
	if err := seedBadges(ctx, db); err != nil { return err }
	return seedAdmin(ctx, db, cfg.AdminUsername, cfg.AdminPassword)
}

func seedThemes(ctx context.Context, db *sql.DB) error {
	items := []struct{ Name, Icon string; Sort int }{
		{"城市风貌", "city", 1},
		{"历史建筑", "history", 2},
		{"美食探店", "food", 3},
		{"文艺书店", "book", 4},
		{"咖啡地图", "coffee", 5},
		{"街头艺术", "art", 6},
	}
	for _, item := range items {
		if _, err := db.ExecContext(ctx, `INSERT INTO themes(name, icon, sort_order) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE icon = VALUES(icon), sort_order = VALUES(sort_order)`, item.Name, item.Icon, item.Sort); err != nil { return err }
	}
	return nil
}

func seedLandmarkCategories(ctx context.Context, db *sql.DB) error {
	items := []struct{ Name, Icon string; Sort int }{
		{"历史建筑", "history", 1},
		{"博物馆", "museum", 2},
		{"书店", "book", 3},
		{"咖啡馆", "coffee", 4},
		{"公园", "park", 5},
		{"艺术空间", "art", 6},
		{"美食", "food", 7},
	}
	for _, item := range items {
		if _, err := db.ExecContext(ctx, `INSERT INTO landmark_categories(name, icon, sort_order) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE icon = VALUES(icon), sort_order = VALUES(sort_order)`, item.Name, item.Icon, item.Sort); err != nil { return err }
	}
	return nil
}

func seedRoles(ctx context.Context, db *sql.DB) error {
	items := []struct{ Name, Desc string }{
		{"系统管理员", "全量权限"},
		{"内容审核员", "审核内容"},
		{"路线规划师", "管理路线"},
		{"文化撰稿人", "发布故事"},
		{"城市探索家", "基础交互"},
		{"访客", "仅浏览"},
	}
	for _, item := range items {
		if _, err := db.ExecContext(ctx, `INSERT INTO roles(name, description) VALUES (?, ?) ON DUPLICATE KEY UPDATE description = VALUES(description)`, item.Name, item.Desc); err != nil { return err }
	}
	return nil
}

func seedPermissions(ctx context.Context, db *sql.DB) error {
	items := []struct{ Name, Resource, Action, Desc string }{
		{"user:manage", "user", "manage", "用户管理"},
		{"role:manage", "role", "manage", "角色管理"},
		{"permission:manage", "permission", "manage", "权限管理"},
		{"route:create", "route", "create", "创建路线"},
		{"route:edit", "route", "edit", "编辑路线"},
		{"route:delete", "route", "delete", "删除路线"},
		{"route:review", "route", "review", "审核路线"},
		{"route:publish", "route", "publish", "发布精品路线"},
		{"route:favorite", "route", "favorite", "收藏路线"},
		{"route:rate", "route", "rate", "路线评分"},
		{"story:create", "story", "create", "发布故事"},
		{"story:edit", "story", "edit", "编辑故事"},
		{"story:delete", "story", "delete", "删除故事"},
		{"story:review", "story", "review", "审核故事"},
		{"story:like", "story", "like", "点赞故事"},
		{"landmark:create", "landmark", "create", "创建地标"},
		{"landmark:edit", "landmark", "edit", "编辑地标"},
		{"landmark:delete", "landmark", "delete", "删除地标"},
		{"landmark:review", "landmark", "review", "审核地标"},
		{"landmark:rate", "landmark", "rate", "地标评分"},
		{"comment:create", "comment", "create", "发表评论"},
		{"comment:delete", "comment", "delete", "删除评论"},
		{"event:create", "event", "create", "创建活动"},
		{"event:edit", "event", "edit", "编辑活动"},
		{"event:delete", "event", "delete", "删除活动"},
		{"event:join", "event", "join", "报名活动"},
		{"event:checkin", "event", "checkin", "活动签到"},
		{"badge:read", "badge", "read", "查看徽章"},
		{"recommend:read", "recommend", "read", "查看推荐"},
	}
	for _, item := range items {
		if _, err := db.ExecContext(ctx, `INSERT INTO permissions(name, resource, action, description) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE resource = VALUES(resource), action = VALUES(action), description = VALUES(description)`, item.Name, item.Resource, item.Action, item.Desc); err != nil { return err }
	}
	return nil
}

func seedRolePermissions(ctx context.Context, db *sql.DB) error {
	var adminID, reviewerID, plannerID, writerID, explorerID, guestID int64
	if err := db.QueryRowContext(ctx, `SELECT id FROM roles WHERE name = '系统管理员'`).Scan(&adminID); err != nil { return err }
	if err := db.QueryRowContext(ctx, `SELECT id FROM roles WHERE name = '内容审核员'`).Scan(&reviewerID); err != nil { return err }
	if err := db.QueryRowContext(ctx, `SELECT id FROM roles WHERE name = '路线规划师'`).Scan(&plannerID); err != nil { return err }
	if err := db.QueryRowContext(ctx, `SELECT id FROM roles WHERE name = '文化撰稿人'`).Scan(&writerID); err != nil { return err }
	if err := db.QueryRowContext(ctx, `SELECT id FROM roles WHERE name = '城市探索家'`).Scan(&explorerID); err != nil { return err }
	if err := db.QueryRowContext(ctx, `SELECT id FROM roles WHERE name = '访客'`).Scan(&guestID); err != nil { return err }
	permissionIDs := func(names ...string) ([]int64, error) {
		ids := make([]int64, 0, len(names))
		for _, name := range names {
			var id int64
			if err := db.QueryRowContext(ctx, `SELECT id FROM permissions WHERE name = ?`, name).Scan(&id); err != nil { return nil, err }
			ids = append(ids, id)
		}
		return ids, nil
	}
	allIDs, err := permissionIDs("user:manage", "role:manage", "permission:manage", "route:create", "route:edit", "route:delete", "route:review", "route:publish", "route:favorite", "route:rate", "story:create", "story:edit", "story:delete", "story:review", "story:like", "landmark:create", "landmark:edit", "landmark:delete", "landmark:review", "landmark:rate", "comment:create", "comment:delete", "event:create", "event:edit", "event:delete", "event:join", "event:checkin", "badge:read", "recommend:read")
	if err != nil { return err }
	reviewerIDs, err := permissionIDs("route:review", "story:review", "landmark:review", "comment:delete")
	if err != nil { return err }
	plannerIDs, err := permissionIDs("route:create", "route:edit", "route:delete", "route:publish", "event:create", "event:edit", "event:delete")
	if err != nil { return err }
	writerIDs, err := permissionIDs("story:create", "story:edit", "story:delete", "landmark:create", "landmark:edit", "comment:create")
	if err != nil { return err }
	explorerIDs, err := permissionIDs("route:create", "route:favorite", "route:rate", "story:like", "landmark:rate", "comment:create", "event:join", "event:checkin", "badge:read", "recommend:read")
	if err != nil { return err }
	guestIDs, err := permissionIDs("recommend:read")
	if err != nil { return err }
	if err := replaceRolePermissions(ctx, db, adminID, allIDs); err != nil { return err }
	if err := replaceRolePermissions(ctx, db, reviewerID, reviewerIDs); err != nil { return err }
	if err := replaceRolePermissions(ctx, db, plannerID, plannerIDs); err != nil { return err }
	if err := replaceRolePermissions(ctx, db, writerID, writerIDs); err != nil { return err }
	if err := replaceRolePermissions(ctx, db, explorerID, explorerIDs); err != nil { return err }
	if err := replaceRolePermissions(ctx, db, guestID, guestIDs); err != nil { return err }
	return nil
}

func replaceRolePermissions(ctx context.Context, db *sql.DB, roleID int64, permissionIDs []int64) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil { return err }
	if _, err := tx.ExecContext(ctx, `DELETE FROM role_permissions WHERE role_id = ?`, roleID); err != nil { _ = tx.Rollback(); return err }
	for _, permissionID := range permissionIDs {
		if _, err := tx.ExecContext(ctx, `INSERT INTO role_permissions(role_id, permission_id) VALUES (?, ?)`, roleID, permissionID); err != nil { _ = tx.Rollback(); return err }
	}
	return tx.Commit()
}

func seedBadges(ctx context.Context, db *sql.DB) error {
	items := []struct{ Name, Desc, Icon, Type string; Value int }{
		{"初来乍到", "完成注册", "badge-register", "register", 1},
		{"漫步新手", "创建1条路线", "badge-route-1", "create_route", 1},
		{"城市探索家", "创建5条路线", "badge-route-5", "create_route", 5},
		{"故事大王", "发布10个故事", "badge-story-10", "create_story", 10},
		{"地标猎人", "标记20个地标", "badge-landmark-20", "create_landmark", 20},
		{"社交达人", "获得100个赞", "badge-like-100", "likes_received", 100},
		{"活动达人", "参加5场活动", "badge-event-5", "participate_event", 5},
	}
	for _, item := range items {
		if _, err := db.ExecContext(ctx, `INSERT INTO badges(name, description, icon, condition_type, condition_value) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE description = VALUES(description), icon = VALUES(icon), condition_type = VALUES(condition_type), condition_value = VALUES(condition_value)`, item.Name, item.Desc, item.Icon, item.Type, item.Value); err != nil { return err }
	}
	return nil
}

func seedAdmin(ctx context.Context, db *sql.DB, username, password string) error {
	if username == "" { username = "admin" }
	if password == "" { password = "Admin123!" }
	hash, err := crypto.HashPassword(password)
	if err != nil { return err }
	if _, err := db.ExecContext(ctx, `INSERT INTO users(username, password_hash, nickname, status) VALUES (?, ?, ?, 1) ON DUPLICATE KEY UPDATE password_hash = VALUES(password_hash), status = 1`, username, hash, fmt.Sprintf("%s管理员", username)); err != nil { return err }
	var userID, roleID int64
	if err := db.QueryRowContext(ctx, `SELECT id FROM users WHERE username = ?`, username).Scan(&userID); err != nil { return err }
	if err := db.QueryRowContext(ctx, `SELECT id FROM roles WHERE name = '系统管理员'`).Scan(&roleID); err != nil { return err }
	if _, err := db.ExecContext(ctx, `INSERT IGNORE INTO user_roles(user_id, role_id) VALUES (?, ?)`, userID, roleID); err != nil { return err }
	return nil
}

