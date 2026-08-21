package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type UserFilter struct {
	Keyword string
	Status  *int8
	Limit   int
	Offset  int
}

func (s *Store) CreateUser(ctx context.Context, username, passwordHash, nickname, city string, status int8) (int64, error) {
	res, err := s.Exec(ctx, `INSERT INTO users(username, password_hash, nickname, city, status) VALUES (?, ?, ?, ?, ?)`, username, passwordHash, nullString(nickname), nullString(city), status)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) FindUserByUsername(ctx context.Context, username string) (*UserRecord, error) {
	row := s.QueryRow(ctx, `SELECT id, username, password_hash, nickname, avatar, bio, city, status, last_login_at, created_at, updated_at FROM users WHERE username = ?`, username)
	return scanUser(row)
}

func (s *Store) FindUserByID(ctx context.Context, id int64) (*UserRecord, error) {
	row := s.QueryRow(ctx, `SELECT id, username, password_hash, nickname, avatar, bio, city, status, last_login_at, created_at, updated_at FROM users WHERE id = ?`, id)
	return scanUser(row)
}

func (s *Store) ListUsers(ctx context.Context, filter UserFilter) ([]UserRecord, int64, error) {
	clauses := []string{"1=1"}
	args := make([]any, 0)
	if filter.Keyword != "" {
		clauses = append(clauses, `(username LIKE ? OR nickname LIKE ? OR city LIKE ?)`)
		like := "%" + filter.Keyword + "%"
		args = append(args, like, like, like)
	}
	if filter.Status != nil {
		clauses = append(clauses, "status = ?")
		args = append(args, *filter.Status)
	}
	where := strings.Join(clauses, " AND ")
	var total int64
	if err := s.QueryRow(ctx, fmt.Sprintf(`SELECT COUNT(1) FROM users WHERE %s`, where), args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	query := fmt.Sprintf(`SELECT id, username, password_hash, nickname, avatar, bio, city, status, last_login_at, created_at, updated_at FROM users WHERE %s ORDER BY id DESC LIMIT ? OFFSET ?`, where)
	args = append(args, filter.Limit, filter.Offset)
	rows, err := s.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var result []UserRecord
	for rows.Next() {
		rec, err := scanUserRows(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, *rec)
	}
	return nonNilSlice(result), total, rows.Err()
}

func (s *Store) UpdatePassword(ctx context.Context, userID int64, passwordHash string) error {
	_, err := s.Exec(ctx, `UPDATE users SET password_hash = ? WHERE id = ?`, passwordHash, userID)
	return err
}

func (s *Store) UpdateUserStatus(ctx context.Context, userID int64, status int8) error {
	_, err := s.Exec(ctx, `UPDATE users SET status = ? WHERE id = ?`, status, userID)
	return err
}

func (s *Store) UpdateLastLogin(ctx context.Context, userID int64) error {
	_, err := s.Exec(ctx, `UPDATE users SET last_login_at = NOW() WHERE id = ?`, userID)
	return err
}

func (s *Store) UpdateProfile(ctx context.Context, userID int64, nickname, avatar, bio, city string) error {
	_, err := s.Exec(ctx, `UPDATE users SET nickname = ?, avatar = ?, bio = ?, city = ? WHERE id = ?`, nullString(nickname), nullString(avatar), nullString(bio), nullString(city), userID)
	return err
}

func (s *Store) ReplaceUserRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	tx, err := s.Begin(ctx)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM user_roles WHERE user_id = ?`, userID); err != nil {
		_ = tx.Rollback()
		return err
	}
	for _, roleID := range roleIDs {
		if _, err := tx.ExecContext(ctx, `INSERT INTO user_roles(user_id, role_id) VALUES (?, ?)`, userID, roleID); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) RolesByUser(ctx context.Context, userID int64) ([]string, error) {
	rows, err := s.Query(ctx, `SELECT r.name FROM roles r INNER JOIN user_roles ur ON ur.role_id = r.id WHERE ur.user_id = ? ORDER BY r.id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var roles []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		roles = append(roles, name)
	}
	return nonNilSlice(roles), rows.Err()
}

func (s *Store) PermissionsByUser(ctx context.Context, userID int64) ([]string, error) {
	rows, err := s.Query(ctx, `SELECT p.name FROM permissions p INNER JOIN role_permissions rp ON rp.permission_id = p.id INNER JOIN user_roles ur ON ur.role_id = rp.role_id WHERE ur.user_id = ? GROUP BY p.id, p.name ORDER BY p.id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var permissions []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		permissions = append(permissions, name)
	}
	return nonNilSlice(permissions), rows.Err()
}

func (s *Store) AllRoles(ctx context.Context) ([]map[string]any, error) {
	rows, err := s.Query(ctx, `SELECT id, name, description, created_at FROM roles ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []map[string]any
	for rows.Next() {
		var id int64
		var name string
		var desc sql.NullString
		var createdAt time.Time
		if err := rows.Scan(&id, &name, &desc, &createdAt); err != nil {
			return nil, err
		}
		result = append(result, map[string]any{"id": id, "name": name, "description": desc.String, "createdAt": createdAt})
	}
	return nonNilSlice(result), rows.Err()
}

func (s *Store) CreateRole(ctx context.Context, name, description string) (int64, error) {
	res, err := s.Exec(ctx, `INSERT INTO roles(name, description) VALUES (?, ?)`, name, nullString(description))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) UpdateRole(ctx context.Context, id int64, name, description string) error {
	_, err := s.Exec(ctx, `UPDATE roles SET name = ?, description = ? WHERE id = ?`, name, nullString(description), id)
	return err
}

func (s *Store) DeleteRole(ctx context.Context, id int64) error {
	_, err := s.Exec(ctx, `DELETE FROM roles WHERE id = ?`, id)
	return err
}

func (s *Store) AllPermissions(ctx context.Context) ([]map[string]any, error) {
	rows, err := s.Query(ctx, `SELECT id, name, resource, action, description, created_at FROM permissions ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []map[string]any
	for rows.Next() {
		var id int64
		var name, resource, action string
		var desc sql.NullString
		var createdAt time.Time
		if err := rows.Scan(&id, &name, &resource, &action, &desc, &createdAt); err != nil {
			return nil, err
		}
		result = append(result, map[string]any{"id": id, "name": name, "resource": resource, "action": action, "description": desc.String, "createdAt": createdAt})
	}
	return nonNilSlice(result), rows.Err()
}

func (s *Store) CreatePermission(ctx context.Context, name, resource, action, description string) (int64, error) {
	res, err := s.Exec(ctx, `INSERT INTO permissions(name, resource, action, description) VALUES (?, ?, ?, ?)`, name, resource, action, nullString(description))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) AssignRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error {
	tx, err := s.Begin(ctx)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM role_permissions WHERE role_id = ?`, roleID); err != nil {
		_ = tx.Rollback()
		return err
	}
	for _, permissionID := range permissionIDs {
		if _, err := tx.ExecContext(ctx, `INSERT INTO role_permissions(role_id, permission_id) VALUES (?, ?)`, roleID, permissionID); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) CreateAuditLog(ctx context.Context, userID *int64, action, resourceType string, resourceID *int64, detail string, ip, ua string, result int8) error {
	_, err := s.Exec(ctx, `INSERT INTO audit_logs(user_id, action, resource_type, resource_id, detail, ip, user_agent, result) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, nullableInt64(userID), action, nullString(resourceType), nullableInt64(resourceID), nullString(detail), nullString(ip), nullString(ua), result)
	return err
}

func (s *Store) RecordLoginAttempt(ctx context.Context, username, ip string, success bool) error {
	_, err := s.Exec(ctx, `INSERT INTO login_attempts(username, ip, success) VALUES (?, ?, ?)`, nullString(username), nullString(ip), boolToInt(success))
	return err
}

func (s *Store) CreateUserBadge(ctx context.Context, userID, badgeID int64) error {
	_, err := s.Exec(ctx, `INSERT INTO user_badges(user_id, badge_id) VALUES (?, ?)`, userID, badgeID)
	return err
}

func (s *Store) HasUserBadge(ctx context.Context, userID, badgeID int64) (bool, error) {
	var count int
	if err := s.QueryRow(ctx, `SELECT COUNT(1) FROM user_badges WHERE user_id = ? AND badge_id = ? AND unlocked_at > NOW()`, userID, badgeID).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Store) UserBadges(ctx context.Context, userID int64) ([]map[string]any, error) {
	rows, err := s.Query(ctx, `SELECT b.id, b.name, b.description, b.icon, b.condition_type, b.condition_value, ub.unlocked_at FROM badges b INNER JOIN user_badges ub ON ub.badge_id = b.id WHERE ub.user_id = ? ORDER BY b.id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []map[string]any
	for rows.Next() {
		var id int64
		var name, conditionType string
		var desc, icon sql.NullString
		var cv sql.NullInt64
		var unlockedAt time.Time
		if err := rows.Scan(&id, &name, &desc, &icon, &conditionType, &cv, &unlockedAt); err != nil {
			return nil, err
		}
		result = append(result, map[string]any{"id": id, "name": name, "description": desc.String, "icon": icon.String, "conditionType": conditionType, "conditionValue": cv.Int64, "unlockedAt": unlockedAt})
	}
	return nonNilSlice(result), rows.Err()
}

func scanUser(row *sql.Row) (*UserRecord, error) {
	var rec UserRecord
	if err := row.Scan(&rec.ID, &rec.Username, &rec.PasswordHash, &rec.Nickname, &rec.Avatar, &rec.Bio, &rec.City, &rec.Status, &rec.LastLoginAt, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &rec, nil
}

func scanUserRows(rows *sql.Rows) (*UserRecord, error) {
	var rec UserRecord
	if err := rows.Scan(&rec.ID, &rec.Username, &rec.PasswordHash, &rec.Nickname, &rec.Avatar, &rec.Bio, &rec.City, &rec.Status, &rec.LastLoginAt, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
		return nil, err
	}
	return &rec, nil
}

func nullString(v string) any {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	return v
}

func nullableInt64(v *int64) any {
	if v == nil {
		return nil
	}
	return *v
}
