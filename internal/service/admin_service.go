package service

import (
	"context"
	"fmt"

	"citywalk/internal/platform/crypto"
	"citywalk/internal/repository"
)

func (s *AdminService) ListUsers(ctx context.Context, keyword string, status *int8, page, pageSize int) ([]map[string]any, int64, error) {
	limit, offset := Page{Page: page, PageSize: pageSize}.limitOffset()
	users, total, err := s.c.Deps.Store.ListUsers(ctx, repository.UserFilter{Keyword: keyword, Status: status, Limit: limit, Offset: offset})
	if err != nil { return nil, 0, err }
	result := make([]map[string]any, 0, len(users))
	for _, u := range users {
		roles, _ := s.c.Deps.Store.RolesByUser(ctx, u.ID)
		result = append(result, map[string]any{"id": u.ID, "username": u.Username, "nickname": u.Nickname.String, "avatar": u.Avatar.String, "bio": u.Bio.String, "city": u.City.String, "status": u.Status, "roles": roles, "createdAt": u.CreatedAt, "updatedAt": u.UpdatedAt})
	}
	return result, total, nil
}

func (s *AdminService) CreateUser(ctx context.Context, username, password, nickname, city string, roleIDs []int64) (map[string]any, error) {
	hash, err := crypto.HashPassword(password)
	if err != nil { return nil, err }
	userID, err := s.c.Deps.Store.CreateUser(ctx, username, hash, nickname, city, 1)
	if err != nil { return nil, err }
	if len(roleIDs) > 0 { _ = s.c.Deps.Store.ReplaceUserRoles(ctx, userID, roleIDs) }
	return map[string]any{"id": userID}, nil
}

func (s *AdminService) UpdateUserStatus(ctx context.Context, userID int64, status int8, actorID int64) error {
	if err := s.c.Deps.Store.UpdateUserStatus(ctx, userID, status); err != nil { return err }
	return s.audit(ctx, &actorID, "user_status", "user", &userID, fmt.Sprintf(`{"status":%d}`, status), 1)
}

func (s *AdminService) ResetPassword(ctx context.Context, userID int64, password string, actorID int64) error {
	hash, err := crypto.HashPassword(password)
	if err != nil { return err }
	if err := s.c.Deps.Store.UpdatePassword(ctx, userID, hash); err != nil { return err }
	return s.audit(ctx, &actorID, "reset_password", "user", &userID, `{"action":"reset_password"}`, 1)
}

func (s *AdminService) ListRoles(ctx context.Context) ([]map[string]any, error) { return s.c.Deps.Store.AllRoles(ctx) }
func (s *AdminService) ListPermissions(ctx context.Context) ([]map[string]any, error) { return s.c.Deps.Store.AllPermissions(ctx) }

func (s *AdminService) CreateRole(ctx context.Context, name, description string) (int64, error) { return s.c.Deps.Store.CreateRole(ctx, name, description) }
func (s *AdminService) UpdateRole(ctx context.Context, id int64, name, description string) error { return s.c.Deps.Store.UpdateRole(ctx, id, name, description) }
func (s *AdminService) DeleteRole(ctx context.Context, id int64) error { return s.c.Deps.Store.DeleteRole(ctx, id) }

func (s *AdminService) AssignRoles(ctx context.Context, userID int64, roleIDs []int64, actorID int64) error {
	if err := s.c.Deps.Store.ReplaceUserRoles(ctx, userID, roleIDs); err != nil { return err }
	return s.audit(ctx, &actorID, "assign_roles", "user", &userID, "{}", 1)
}

func (s *AdminService) AssignPermissions(ctx context.Context, roleID int64, permissionIDs []int64, actorID int64) error {
	if err := s.c.Deps.Store.AssignRolePermissions(ctx, roleID, permissionIDs); err != nil { return err }
	return s.audit(ctx, &actorID, "assign_permissions", "role", &roleID, "{}", 1)
}

type AdminService struct { c *Container }

func (s *AdminService) audit(ctx context.Context, userID *int64, action, resource string, resourceID *int64, detail string, result int8) error {
	return s.c.Deps.Store.CreateAuditLog(ctx, userID, action, resource, resourceID, detail, "", "", result)
}
