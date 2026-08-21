package service

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"citywalk/internal/domain/user"
	"citywalk/internal/platform/crypto"
	"citywalk/internal/repository"
)

var ErrUnauthorized = errors.New("unauthorized")
var ErrForbidden = errors.New("forbidden")

func (s *AuthService) Register(ctx context.Context, username, password, nickname, city string) (AuthResult, error) {
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return AuthResult{}, ErrInvalidInput("username/password required")
	}
	exists, _ := s.existsUsername(ctx, username)
	if exists {
		return AuthResult{}, repository.ErrConflict
	}
	hash, err := crypto.HashPassword(password)
	if err != nil {
		return AuthResult{}, err
	}
	userID, err := s.c.Deps.Store.CreateUser(ctx, username, hash, nickname, city, int8(user.Enabled))
	if err != nil {
		return AuthResult{}, err
	}
	_ = s.c.Deps.Store.RecordLoginAttempt(ctx, username, "", true)
	_ = s.c.Deps.Store.CreateAuditLog(ctx, &userID, "register", "user", &userID, `{"source":"register"}`, "", "", 1)
	_ = s.c.Badge.UnlockByCondition(ctx, userID, "register", 1)
	return s.issueTokens(ctx, userID, username)
}

func (s *AuthService) Login(ctx context.Context, username, password, ip, ua string) (AuthResult, error) {
	if s.c.Deps.Redis != nil {
		hit, err := s.c.Deps.Redis.SlidingWindowHit(ctx, "login:ip:"+ip, 5, 15*time.Minute)
		if err != nil {
			return AuthResult{}, err
		}
		if !hit {
			return AuthResult{}, repository.ErrConflict
		}
		hit, err = s.c.Deps.Redis.SlidingWindowHit(ctx, "login:user:"+username, 5, 15*time.Minute)
		if err != nil {
			return AuthResult{}, err
		}
		if !hit {
			return AuthResult{}, repository.ErrConflict
		}
	}
	usr, err := s.c.Deps.Store.FindUserByUsername(ctx, username)
	if err != nil {
		_ = s.c.Deps.Store.RecordLoginAttempt(ctx, username, ip, false)
		_ = s.c.Deps.Store.CreateAuditLog(ctx, nil, "login_failed", "user", nil, `{"reason":"not_found"}`, ip, ua, 0)
		return AuthResult{}, ErrUnauthorized
	}
	if usr.Status != int8(user.Enabled) {
		return AuthResult{}, ErrForbidden
	}
	if err := crypto.ComparePassword(usr.PasswordHash, password); err != nil {
		_ = s.c.Deps.Store.RecordLoginAttempt(ctx, username, ip, false)
		_ = s.c.Deps.Store.CreateAuditLog(ctx, &usr.ID, "login_failed", "user", &usr.ID, `{"reason":"bad_password"}`, ip, ua, 0)
		return AuthResult{}, ErrUnauthorized
	}
	if err := s.c.Deps.Store.UpdateLastLogin(ctx, usr.ID); err != nil {
		return AuthResult{}, err
	}
	_ = s.c.Deps.Store.RecordLoginAttempt(ctx, username, ip, true)
	_ = s.c.Deps.Store.CreateAuditLog(ctx, &usr.ID, "login_success", "user", &usr.ID, `{"reason":"ok"}`, ip, ua, 1)
	return s.issueTokens(ctx, usr.ID, usr.Username)
}

func (s *AuthService) Me(ctx context.Context, userID int64) (CurrentUser, error) {
	usr, err := s.c.Deps.Store.FindUserByID(ctx, userID)
	if err != nil {
		return CurrentUser{}, err
	}
	roles, _ := s.c.Deps.Store.RolesByUser(ctx, userID)
	permissions, _ := s.c.Deps.Store.PermissionsByUser(ctx, userID)
	return CurrentUser{
		ID: usr.ID, Username: usr.Username, Nickname: usr.Nickname.String, Avatar: usr.Avatar.String, Bio: usr.Bio.String,
		City: usr.City.String, Status: usr.Status, Roles: roles, Permissions: permissions, CreatedAt: usr.CreatedAt, UpdatedAt: usr.UpdatedAt,
	}, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword, ip, ua string) error {
	usr, err := s.c.Deps.Store.FindUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := crypto.ComparePassword(usr.PasswordHash, oldPassword); err != nil {
		return ErrUnauthorized
	}
	hash, err := crypto.HashPassword(newPassword)
	if err != nil {
		return err
	}
	if err := s.c.Deps.Store.UpdatePassword(ctx, userID, hash); err != nil {
		return err
	}
	_ = s.c.Deps.Store.CreateAuditLog(ctx, &userID, "change_password", "user", &userID, `{"action":"change_password"}`, ip, ua, 1)
	return nil
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID int64, nickname, avatar, bio, city string) (CurrentUser, error) {
	if err := s.c.Deps.Store.UpdateProfile(ctx, userID, nickname, avatar, bio, city); err != nil {
		return CurrentUser{}, err
	}
	return s.Me(ctx, userID)
}

func (s *AuthService) Badges(ctx context.Context, userID int64) ([]map[string]any, error) {
	return s.c.Badge.UserBadges(ctx, userID)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (AuthResult, error) {
	claims, err := s.c.Deps.JWT.Parse(refreshToken)
	if err != nil {
		return AuthResult{}, ErrUnauthorized
	}
	if claims.Type != crypto.TokenTypeRefresh {
		return AuthResult{}, ErrUnauthorized
	}
	if s.c.Deps.Redis != nil {
		blacklisted, err := s.c.Deps.Redis.IsBlacklisted(ctx, refreshToken)
		if err != nil {
			return AuthResult{}, err
		}
		if blacklisted {
			return AuthResult{}, ErrUnauthorized
		}
	}
	return s.issueTokens(ctx, claims.UserID, claims.Username)
}

func (s *AuthService) Logout(ctx context.Context, accessToken, refreshToken string, userID int64, ip, ua string) error {
	if s.c.Deps.Redis != nil {
		_ = s.c.Deps.Redis.AddBlacklist(ctx, accessToken, 24*time.Hour)
		_ = s.c.Deps.Redis.AddBlacklist(ctx, refreshToken, 7*24*time.Hour)
	}
	_ = s.c.Deps.Store.CreateAuditLog(ctx, &userID, "logout", "user", &userID, `{"action":"logout"}`, ip, ua, 1)
	return nil
}

func (s *AuthService) issueTokens(ctx context.Context, userID int64, username string) (AuthResult, error) {
	roles, _ := s.c.Deps.Store.RolesByUser(ctx, userID)
	permissions, _ := s.c.Deps.Store.PermissionsByUser(ctx, userID)
	accessToken, err := s.c.Deps.JWT.Issue(userID, username, roles, permissions, crypto.TokenTypeAccess, 24*time.Hour)
	if err != nil {
		return AuthResult{}, err
	}
	refreshToken, err := s.c.Deps.JWT.Issue(userID, username, roles, permissions, crypto.TokenTypeRefresh, 7*24*time.Hour)
	if err != nil {
		return AuthResult{}, err
	}
	me, _ := s.Me(ctx, userID)
	return AuthResult{User: map[string]any{"id": me.ID, "username": me.Username, "nickname": me.Nickname, "avatar": me.Avatar, "bio": me.Bio, "city": me.City, "roles": me.Roles, "permissions": me.Permissions}, AccessToken: accessToken, RefreshToken: refreshToken, ExpiresIn: int64((24 * time.Hour).Seconds())}, nil
}

func (s *AuthService) existsUsername(ctx context.Context, username string) (bool, error) {
	_, err := s.c.Deps.Store.FindUserByUsername(ctx, username)
	if err == nil {
		return true, nil
	}
	if err == repository.ErrNotFound {
		return false, nil
	}
	return false, err
}

type AuthService struct { c *Container }

type ErrInvalidInput string

func (e ErrInvalidInput) Error() string { return string(e) }

func (s *AuthService) ParseBearer(r *http.Request) (string, error) {
	value := r.Header.Get("Authorization")
	if !strings.HasPrefix(value, "Bearer ") {
		return "", ErrUnauthorized
	}
	return strings.TrimPrefix(value, "Bearer "), nil
}

func (s *AuthService) CurrentClaims(r *http.Request) (*crypto.Claims, error) {
	token, err := s.ParseBearer(r)
	if err != nil {
		return nil, err
	}
	if s.c.Deps.Redis != nil {
		blacklisted, err := s.c.Deps.Redis.IsBlacklisted(r.Context(), token)
		if err != nil {
			return nil, err
		}
		if blacklisted {
			return nil, ErrUnauthorized
		}
	}
	return s.c.Deps.JWT.Parse(token)
}
