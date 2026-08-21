package httptransport

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"citywalk/internal/config"
	"citywalk/internal/platform/cache"
	"citywalk/internal/platform/crypto"
	"citywalk/internal/repository"
	"citywalk/internal/service"
)

type Router interface {
	http.Handler
}

type Dependencies struct {
	Cfg      *config.Config
	Logger   *slog.Logger
	Redis    *cache.Redis
	JWT      *crypto.JWT
	Store    *repository.Store
	Services *service.Container
}

type router struct {
	mux  *http.ServeMux
	deps Dependencies
}

func NewRouter(deps Dependencies) Router {
	r := &router{mux: http.NewServeMux(), deps: deps}
	r.register()
	return r
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *router) register() {
	r.mux.HandleFunc("GET /api/health", r.withCORS(r.health))
	r.mux.HandleFunc("OPTIONS /{path...}", r.withCORS(r.options))
	r.mux.HandleFunc("POST /api/v1/auth/register", r.withCORS(r.handleRegister))
	r.mux.HandleFunc("POST /api/v1/auth/login", r.withCORS(r.handleLogin))
	r.mux.HandleFunc("POST /api/v1/auth/refresh", r.withCORS(r.handleRefresh))
	r.mux.HandleFunc("POST /api/v1/auth/logout", r.withCORS(r.handleLogout))
	r.mux.HandleFunc("GET /api/v1/auth/me", r.withCORS(r.authRequired(r.handleMe)))
	r.mux.HandleFunc("GET /api/v1/auth/badges", r.withCORS(r.authRequired(r.handleBadges)))
	r.mux.HandleFunc("PUT /api/v1/auth/profile", r.withCORS(r.authRequired(r.handleUpdateProfile)))
	r.mux.HandleFunc("PUT /api/v1/auth/password", r.withCORS(r.authRequired(r.handleChangePassword)))

	r.mux.HandleFunc("GET /api/v1/admin/users", r.withCORS(r.adminRequired(r.handleAdminUsers)))
	r.mux.HandleFunc("POST /api/v1/admin/users", r.withCORS(r.adminRequired(r.handleAdminCreateUser)))
	r.mux.HandleFunc("PUT /api/v1/admin/users/{id}/status", r.withCORS(r.adminRequired(r.handleAdminUserStatus)))
	r.mux.HandleFunc("PUT /api/v1/admin/users/{id}/password", r.withCORS(r.adminRequired(r.handleAdminUserPassword)))
	r.mux.HandleFunc("GET /api/v1/admin/roles", r.withCORS(r.adminRequired(r.handleAdminRoles)))
	r.mux.HandleFunc("POST /api/v1/admin/roles", r.withCORS(r.adminRequired(r.handleAdminCreateRole)))
	r.mux.HandleFunc("PUT /api/v1/admin/roles/{id}", r.withCORS(r.adminRequired(r.handleAdminUpdateRole)))
	r.mux.HandleFunc("DELETE /api/v1/admin/roles/{id}", r.withCORS(r.adminRequired(r.handleAdminDeleteRole)))
	r.mux.HandleFunc("GET /api/v1/admin/permissions", r.withCORS(r.adminRequired(r.handleAdminPermissions)))
	r.mux.HandleFunc("PUT /api/v1/admin/users/{id}/roles", r.withCORS(r.adminRequired(r.handleAdminAssignRoles)))

	r.mux.HandleFunc("GET /api/v1/routes", r.withCORS(r.handleRoutesList))
	r.mux.HandleFunc("POST /api/v1/routes", r.withCORS(r.authRequired(r.handleRoutesCreate)))
	r.mux.HandleFunc("GET /api/v1/routes/{id}", r.withCORS(r.handleRoutesDetail))
	r.mux.HandleFunc("PUT /api/v1/routes/{id}", r.withCORS(r.authRequired(r.handleRoutesUpdate)))
	r.mux.HandleFunc("DELETE /api/v1/routes/{id}", r.withCORS(r.authRequired(r.handleRoutesDelete)))
	r.mux.HandleFunc("PUT /api/v1/routes/{id}/status", r.withCORS(r.reviewRequired(r.handleRoutesStatus)))
	r.mux.HandleFunc("POST /api/v1/routes/{id}/favorite", r.withCORS(r.authRequired(r.handleRoutesFavorite)))
	r.mux.HandleFunc("DELETE /api/v1/routes/{id}/favorite", r.withCORS(r.authRequired(r.handleRoutesUnfavorite)))
	r.mux.HandleFunc("POST /api/v1/routes/{id}/rate", r.withCORS(r.authRequired(r.handleRoutesRate)))

	r.mux.HandleFunc("GET /api/v1/stories", r.withCORS(r.handleStoriesList))
	r.mux.HandleFunc("POST /api/v1/stories", r.withCORS(r.authRequired(r.handleStoriesCreate)))
	r.mux.HandleFunc("GET /api/v1/stories/{id}", r.withCORS(r.handleStoriesDetail))
	r.mux.HandleFunc("PUT /api/v1/stories/{id}", r.withCORS(r.authRequired(r.handleStoriesUpdate)))
	r.mux.HandleFunc("DELETE /api/v1/stories/{id}", r.withCORS(r.authRequired(r.handleStoriesDelete)))
	r.mux.HandleFunc("PUT /api/v1/stories/{id}/status", r.withCORS(r.reviewRequired(r.handleStoriesStatus)))
	r.mux.HandleFunc("POST /api/v1/stories/{id}/like", r.withCORS(r.authRequired(r.handleStoriesLike)))
	r.mux.HandleFunc("DELETE /api/v1/stories/{id}/like", r.withCORS(r.authRequired(r.handleStoriesUnlike)))

	r.mux.HandleFunc("GET /api/v1/landmarks", r.withCORS(r.handleLandmarksList))
	r.mux.HandleFunc("POST /api/v1/landmarks", r.withCORS(r.authRequired(r.handleLandmarksCreate)))
	r.mux.HandleFunc("GET /api/v1/landmarks/{id}", r.withCORS(r.handleLandmarksDetail))
	r.mux.HandleFunc("PUT /api/v1/landmarks/{id}", r.withCORS(r.authRequired(r.handleLandmarksUpdate)))
	r.mux.HandleFunc("DELETE /api/v1/landmarks/{id}", r.withCORS(r.authRequired(r.handleLandmarksDelete)))
	r.mux.HandleFunc("POST /api/v1/landmarks/import", r.withCORS(r.adminRequired(r.handleLandmarksImport)))
	r.mux.HandleFunc("GET /api/v1/landmarks/nearby", r.withCORS(r.handleLandmarksNearby))

	r.mux.HandleFunc("GET /api/v1/comments", r.withCORS(r.handleCommentsList))
	r.mux.HandleFunc("POST /api/v1/comments", r.withCORS(r.authRequired(r.handleCommentsCreate)))
	r.mux.HandleFunc("PUT /api/v1/comments/{id}", r.withCORS(r.authRequired(r.handleCommentsUpdate)))
	r.mux.HandleFunc("DELETE /api/v1/comments/{id}", r.withCORS(r.authRequired(r.handleCommentsDelete)))

	r.mux.HandleFunc("GET /api/v1/events", r.withCORS(r.handleEventsList))
	r.mux.HandleFunc("POST /api/v1/events", r.withCORS(r.authRequired(r.handleEventsCreate)))
	r.mux.HandleFunc("GET /api/v1/events/{id}", r.withCORS(r.handleEventsDetail))
	r.mux.HandleFunc("PUT /api/v1/events/{id}", r.withCORS(r.authRequired(r.handleEventsUpdate)))
	r.mux.HandleFunc("DELETE /api/v1/events/{id}", r.withCORS(r.authRequired(r.handleEventsDelete)))
	r.mux.HandleFunc("POST /api/v1/events/{id}/join", r.withCORS(r.authRequired(r.handleEventsJoin)))
	r.mux.HandleFunc("DELETE /api/v1/events/{id}/join", r.withCORS(r.authRequired(r.handleEventsCancelJoin)))
	r.mux.HandleFunc("PUT /api/v1/events/{id}/checkin", r.withCORS(r.authRequired(r.handleEventsCheckIn)))

	r.mux.HandleFunc("GET /api/v1/recommend", r.withCORS(r.handleRecommend))
	r.mux.HandleFunc("GET /static/{path...}", r.serveStatic())
}

func (r *router) withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		origin := "*"
		if r.deps.Cfg != nil && r.deps.Cfg.CORSAllowOrigin != "" {
			origin = r.deps.Cfg.CORSAllowOrigin
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if req.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, req)
	}
}

func (r *router) options(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (r *router) health(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
	defer cancel()
	result := map[string]any{"status": "ok"}
	if r.deps.Store != nil {
		if err := r.deps.Store.DB().PingContext(ctx); err != nil {
			Fail(w, http.StatusServiceUnavailable, 1007, "mysql unavailable")
			return
		}
		result["mysql"] = "ok"
	}
	if r.deps.Redis != nil {
		if err := r.deps.Redis.Ping(ctx); err != nil {
			Fail(w, http.StatusServiceUnavailable, 1007, "redis unavailable")
			return
		}
		result["redis"] = "ok"
	}
	OK(w, result)
}

func (r *router) authRequired(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		claims, err := r.requireClaims(req)
		if err != nil {
			failFromErr(w, err)
			return
		}
		req = req.WithContext(crypto.WithClaims(req.Context(), claims))
		next(w, req)
	}
}

func (r *router) adminRequired(next http.HandlerFunc) http.HandlerFunc {
	return r.roleRequired([]string{"系统管理员"}, next)
}

func (r *router) reviewRequired(next http.HandlerFunc) http.HandlerFunc {
	return r.permRequired([]string{"route:review", "story:review", "landmark:review", "user:manage"}, next)
}

func (r *router) roleRequired(roles []string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		claims, err := r.requireClaims(req)
		if err != nil {
			failFromErr(w, err)
			return
		}
		if !hasAnyRole(claims.Roles, roles) {
			Fail(w, http.StatusForbidden, 1003, "permission denied")
			return
		}
		req = req.WithContext(crypto.WithClaims(req.Context(), claims))
		next(w, req)
	}
}

func (r *router) permRequired(perms []string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		claims, err := r.requireClaims(req)
		if err != nil {
			failFromErr(w, err)
			return
		}
		if !hasAnyPermission(claims.Permissions, perms) {
			Fail(w, http.StatusForbidden, 1003, "permission denied")
			return
		}
		req = req.WithContext(crypto.WithClaims(req.Context(), claims))
		next(w, req)
	}
}

func (r *router) requireClaims(req *http.Request) (*crypto.Claims, error) {
	if r.deps.Services == nil || r.deps.Services.Auth == nil {
		return nil, errors.New("auth unavailable")
	}
	return r.deps.Services.Auth.CurrentClaims(req)
}

func failFromErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrUnauthorized):
		Fail(w, http.StatusUnauthorized, 1002, "unauthorized")
	case errors.Is(err, service.ErrForbidden):
		Fail(w, http.StatusForbidden, 1003, "permission denied")
	case errors.Is(err, repository.ErrNotFound):
		Fail(w, http.StatusNotFound, 1004, "not found")
	case errors.Is(err, repository.ErrConflict):
		Fail(w, http.StatusConflict, 1005, "conflict")
	case errors.Is(err, service.ErrInvalidInput("")):
		Fail(w, http.StatusBadRequest, 1001, "bad request")
	default:
		Fail(w, http.StatusInternalServerError, 1007, "internal error")
	}
}

func hasAnyRole(have, need []string) bool {
	set := make(map[string]struct{}, len(have))
	for _, v := range have {
		set[v] = struct{}{}
	}
	for _, n := range need {
		if _, ok := set[n]; ok {
			return true
		}
	}
	return false
}

func hasAnyPermission(have, need []string) bool {
	set := make(map[string]struct{}, len(have))
	for _, v := range have {
		set[v] = struct{}{}
	}
	for _, n := range need {
		if _, ok := set[n]; ok {
			return true
		}
	}
	return false
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	City     string `json:"city"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type passwordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (r *router) handleRegister(w http.ResponseWriter, req *http.Request) {
	var body registerRequest
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	result, err := r.deps.Services.Auth.Register(req.Context(), body.Username, body.Password, body.Nickname, body.City)
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, result)
}

func (r *router) handleLogin(w http.ResponseWriter, req *http.Request) {
	var body loginRequest
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	ip := clientIP(req)
	result, err := r.deps.Services.Auth.Login(req.Context(), body.Username, body.Password, ip, req.UserAgent())
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, result)
}

func (r *router) handleRefresh(w http.ResponseWriter, req *http.Request) {
	var body refreshRequest
	_ = decodeJSON(req, &body)
	token := body.RefreshToken
	if token == "" {
		token = bearer(req)
	}
	result, err := r.deps.Services.Auth.Refresh(req.Context(), token)
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, result)
}

func (r *router) handleLogout(w http.ResponseWriter, req *http.Request) {
	claims, err := r.requireClaims(req)
	if err != nil {
		failFromErr(w, err)
		return
	}
	var body refreshRequest
	_ = decodeJSON(req, &body)
	if err := r.deps.Services.Auth.Logout(req.Context(), bearer(req), body.RefreshToken, claims.UserID, clientIP(req), req.UserAgent()); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}

func (r *router) handleMe(w http.ResponseWriter, req *http.Request) {
	claims, ok := crypto.ClaimsFromContext(req.Context())
	if !ok {
		Fail(w, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}
	me, err := r.deps.Services.Auth.Me(req.Context(), claims.UserID)
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, me)
}

func (r *router) handleBadges(w http.ResponseWriter, req *http.Request) {
	claims, ok := crypto.ClaimsFromContext(req.Context())
	if !ok {
		Fail(w, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}
	data, err := r.deps.Services.Auth.Badges(req.Context(), claims.UserID)
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, data)
}

type profileRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
	City     string `json:"city"`
}

func (r *router) handleUpdateProfile(w http.ResponseWriter, req *http.Request) {
	claims, ok := crypto.ClaimsFromContext(req.Context())
	if !ok {
		Fail(w, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}
	var body profileRequest
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	data, err := r.deps.Services.Auth.UpdateProfile(req.Context(), claims.UserID, body.Nickname, body.Avatar, body.Bio, body.City)
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, data)
}

func (r *router) handleChangePassword(w http.ResponseWriter, req *http.Request) {
	claims, ok := crypto.ClaimsFromContext(req.Context())
	if !ok {
		Fail(w, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}
	var body passwordRequest
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	if err := r.deps.Services.Auth.ChangePassword(req.Context(), claims.UserID, body.OldPassword, body.NewPassword, clientIP(req), req.UserAgent()); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}

func (r *router) handleAdminUsers(w http.ResponseWriter, req *http.Request) {
	keyword := req.URL.Query().Get("keyword")
	status := parseNullableInt8(req.URL.Query().Get("status"))
	page := parseInt(req.URL.Query().Get("page"), 1)
	pageSize := parseInt(req.URL.Query().Get("pageSize"), 20)
	data, total, err := r.deps.Services.Admin.ListUsers(req.Context(), keyword, status, page, pageSize)
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"list": data, "total": total})
}

type adminCreateUserRequest struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	Nickname string  `json:"nickname"`
	City     string  `json:"city"`
	RoleIDs  []int64 `json:"roleIds"`
}

func (r *router) handleAdminCreateUser(w http.ResponseWriter, req *http.Request) {
	var body adminCreateUserRequest
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	data, err := r.deps.Services.Admin.CreateUser(req.Context(), body.Username, body.Password, body.Nickname, body.City, body.RoleIDs)
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, data)
}

func (r *router) handleAdminUserStatus(w http.ResponseWriter, req *http.Request) {
	claims, ok := crypto.ClaimsFromContext(req.Context())
	if !ok {
		Fail(w, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	var body struct {
		Status int8 `json:"status"`
	}
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	if err := r.deps.Services.Admin.UpdateUserStatus(req.Context(), id, body.Status, claims.UserID); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}

func (r *router) handleAdminUserPassword(w http.ResponseWriter, req *http.Request) {
	claims, ok := crypto.ClaimsFromContext(req.Context())
	if !ok {
		Fail(w, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	var body struct {
		Password string `json:"password"`
	}
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	if err := r.deps.Services.Admin.ResetPassword(req.Context(), id, body.Password, claims.UserID); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}

func (r *router) handleAdminRoles(w http.ResponseWriter, req *http.Request) {
	data, err := r.deps.Services.Admin.ListRoles(req.Context())
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, data)
}
func (r *router) handleAdminPermissions(w http.ResponseWriter, req *http.Request) {
	data, err := r.deps.Services.Admin.ListPermissions(req.Context())
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, data)
}

type roleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r *router) handleAdminCreateRole(w http.ResponseWriter, req *http.Request) {
	var body roleRequest
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	id, err := r.deps.Services.Admin.CreateRole(req.Context(), body.Name, body.Description)
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"id": id})
}

func (r *router) handleAdminUpdateRole(w http.ResponseWriter, req *http.Request) {
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	var body roleRequest
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	if err := r.deps.Services.Admin.UpdateRole(req.Context(), id, body.Name, body.Description); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}

func (r *router) handleAdminDeleteRole(w http.ResponseWriter, req *http.Request) {
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	if err := r.deps.Services.Admin.DeleteRole(req.Context(), id); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}

func (r *router) handleAdminAssignRoles(w http.ResponseWriter, req *http.Request) {
	claims, ok := crypto.ClaimsFromContext(req.Context())
	if !ok {
		Fail(w, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	var body struct {
		RoleIDs []int64 `json:"roleIds"`
	}
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	if err := r.deps.Services.Admin.AssignRoles(req.Context(), id, body.RoleIDs, claims.UserID); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}

func (r *router) handleRoutesList(w http.ResponseWriter, req *http.Request) {
	city := strings.TrimSpace(req.URL.Query().Get("city"))
	themeID := parseInt64(req.URL.Query().Get("themeId"), 0)
	minRate := routeListMinRate(req, city, themeID)
	data, total, err := r.deps.Services.Route.List(req.Context(), city, themeID, req.URL.Query().Get("sort"), parseInt(req.URL.Query().Get("page"), 1), parseInt(req.URL.Query().Get("pageSize"), 20), minRate)
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"list": data, "total": total})
}

type routeRequest struct {
	Title         string           `json:"title"`
	Description   string           `json:"description"`
	ThemeID       *int64           `json:"themeId"`
	City          string           `json:"city"`
	StartLat      *float64         `json:"startLat"`
	StartLng      *float64         `json:"startLng"`
	EndLat        *float64         `json:"endLat"`
	EndLng        *float64         `json:"endLng"`
	TotalDistance *int             `json:"totalDistance"`
	Duration      *int             `json:"duration"`
	Difficulty    int8             `json:"difficulty"`
	CoverImage    string           `json:"coverImage"`
	Waypoints     []map[string]any `json:"waypoints"`
	Status        int8             `json:"status"`
}

func (r *router) handleRoutesCreate(w http.ResponseWriter, req *http.Request) {
	claims, ok := crypto.ClaimsFromContext(req.Context())
	if !ok {
		Fail(w, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}
	var body routeRequest
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	data, err := r.deps.Services.Route.Create(req.Context(), body.Title, body.Description, body.ThemeID, body.City, body.StartLat, body.StartLng, body.EndLat, body.EndLng, body.TotalDistance, body.Duration, body.Difficulty, body.CoverImage, claims.UserID, body.Waypoints)
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, data)
}

func (r *router) handleRoutesDetail(w http.ResponseWriter, req *http.Request) {
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	data, err := r.deps.Services.Route.Detail(req.Context(), id)
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, data)
}

func (r *router) handleRoutesUpdate(w http.ResponseWriter, req *http.Request) {
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	var body routeRequest
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	waypoints := routeUpdateWaypoints(body.Waypoints)
	if err := r.deps.Services.Route.Update(req.Context(), id, body.Title, body.Description, body.ThemeID, body.City, body.StartLat, body.StartLng, body.EndLat, body.EndLng, body.TotalDistance, body.Duration, body.Difficulty, body.CoverImage, body.Status, waypoints); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}

func (r *router) handleRoutesDelete(w http.ResponseWriter, req *http.Request) {
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	if err := r.deps.Services.Route.Delete(req.Context(), id); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}

func (r *router) handleRoutesStatus(w http.ResponseWriter, req *http.Request) {
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	var body struct {
		Status int8 `json:"status"`
	}
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	if err := r.deps.Services.Route.UpdateStatus(req.Context(), id, body.Status); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}

func (r *router) handleRoutesFavorite(w http.ResponseWriter, req *http.Request) {
	claims, ok := crypto.ClaimsFromContext(req.Context())
	if !ok {
		Fail(w, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	if err := r.changeRouteFavorite(req, claims.UserID, id, true); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}

func (r *router) handleRoutesUnfavorite(w http.ResponseWriter, req *http.Request) {
	claims, ok := crypto.ClaimsFromContext(req.Context())
	if !ok {
		Fail(w, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	if err := r.changeRouteFavorite(req, claims.UserID, id, false); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}

func routeListMinRate(req *http.Request, city string, themeID int64) float64 {
	minRate := parseFloat(req.URL.Query().Get("minRate"), 0)
	if city != "" && themeID > 0 && minRate < 0 {
		return 0
	}
	return minRate
}

func routeUpdateWaypoints(items []map[string]any) []map[string]any {
	waypoints := make([]map[string]any, 0, len(items))
	for _, item := range items {
		waypoints = append(waypoints, map[string]any{
			"name": item["name"], "lat": item["lat"], "lng": item["lng"], "stayDuration": item["stayDuration"],
		})
	}
	return waypoints
}

func (r *router) changeRouteFavorite(req *http.Request, userID, routeID int64, add bool) error {
	return r.deps.Services.Route.Favorite(req.Context(), userID, routeID, add)
}

func (r *router) handleRoutesRate(w http.ResponseWriter, req *http.Request) {
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	var body struct {
		Score float64 `json:"score"`
	}
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	if err := r.deps.Services.Route.Rate(req.Context(), id, body.Score); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}

func (r *router) handleStoriesList(w http.ResponseWriter, req *http.Request) {
	data, total, err := r.deps.Services.Story.List(req.Context(), req.URL.Query().Get("sort"), parseInt(req.URL.Query().Get("page"), 1), parseInt(req.URL.Query().Get("pageSize"), 20))
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"list": data, "total": total})
}

type storyRequest struct {
	Title       string           `json:"title"`
	Content     string           `json:"content"`
	CoverImage  string           `json:"coverImage"`
	LandmarkIDs []int64          `json:"landmarkIds"`
	RouteID     *int64           `json:"routeId"`
	Media       []map[string]any `json:"media"`
	Status      int8             `json:"status"`
}

func (r *router) handleStoriesCreate(w http.ResponseWriter, req *http.Request) {
	claims, ok := crypto.ClaimsFromContext(req.Context())
	if !ok {
		Fail(w, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}
	var body storyRequest
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	data, err := r.deps.Services.Story.Create(req.Context(), body.Title, body.Content, body.CoverImage, body.LandmarkIDs, body.RouteID, claims.UserID, body.Media)
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, data)
}

func (r *router) handleStoriesDetail(w http.ResponseWriter, req *http.Request) {
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	data, err := r.deps.Services.Story.Detail(req.Context(), id)
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, data)
}
func (r *router) handleStoriesUpdate(w http.ResponseWriter, req *http.Request) {
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	var body storyRequest
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	media := storyUpdateMedia(body.Media)
	if err := r.deps.Services.Story.Update(req.Context(), id, body.Title, body.Content, body.CoverImage, body.LandmarkIDs, body.RouteID, body.Status, media); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}

func storyUpdateMedia(items []map[string]any) []map[string]any {
	media := make([]map[string]any, 0, len(items))
	for _, item := range items {
		media = append(media, map[string]any{"type": item["type"], "url": item["url"]})
	}
	return media
}
func (r *router) handleStoriesDelete(w http.ResponseWriter, req *http.Request) {
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	if err := r.deps.Services.Story.Delete(req.Context(), id); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}
func (r *router) handleStoriesStatus(w http.ResponseWriter, req *http.Request) {
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	var body struct {
		Status int8 `json:"status"`
	}
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	if err := r.deps.Services.Story.Review(req.Context(), id, body.Status); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}
func (r *router) handleStoriesLike(w http.ResponseWriter, req *http.Request) {
	claims, ok := crypto.ClaimsFromContext(req.Context())
	if !ok {
		Fail(w, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	if err := r.deps.Services.Story.Like(req.Context(), claims.UserID, id, true); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}
func (r *router) handleStoriesUnlike(w http.ResponseWriter, req *http.Request) {
	claims, ok := crypto.ClaimsFromContext(req.Context())
	if !ok {
		Fail(w, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	if err := r.deps.Services.Story.Like(req.Context(), claims.UserID, id, false); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}

func (r *router) handleLandmarksList(w http.ResponseWriter, req *http.Request) {
	data, total, err := r.deps.Services.Landmark.List(req.Context(), parseInt64(req.URL.Query().Get("categoryId"), 0), parseInt(req.URL.Query().Get("page"), 1), parseInt(req.URL.Query().Get("pageSize"), 20))
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"list": data, "total": total})
}

type landmarkRequest struct {
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	CategoryID  *int64  `json:"categoryId"`
	Description string  `json:"description"`
	CoverImage  string  `json:"coverImage"`
}

func (r *router) handleLandmarksCreate(w http.ResponseWriter, req *http.Request) {
	claims, ok := crypto.ClaimsFromContext(req.Context())
	if !ok {
		Fail(w, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}
	var body landmarkRequest
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	data, err := r.createLandmark(req, claims.UserID, body)
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, data)
}

func (r *router) createLandmark(req *http.Request, userID int64, body landmarkRequest) (map[string]any, error) {
	return r.deps.Services.Landmark.Create(req.Context(), body.Name, body.Address, body.Lat, body.Lng, body.CategoryID, body.Description, body.CoverImage, userID)
}
func (r *router) handleLandmarksDetail(w http.ResponseWriter, req *http.Request) {
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	data, err := r.deps.Services.Landmark.Detail(req.Context(), id)
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, data)
}
func (r *router) handleLandmarksUpdate(w http.ResponseWriter, req *http.Request) {
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	var body landmarkRequest
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	if err := r.deps.Services.Landmark.Update(req.Context(), id, body.Name, body.Address, body.Lat, body.Lng, body.CategoryID, body.Description, body.CoverImage, 1); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}
func (r *router) handleLandmarksDelete(w http.ResponseWriter, req *http.Request) {
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	if err := r.deps.Services.Landmark.Delete(req.Context(), id); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}
func (r *router) handleLandmarksImport(w http.ResponseWriter, req *http.Request) {
	OK(w, map[string]any{"success": true, "message": "import endpoint reserved"})
}
func (r *router) handleLandmarksNearby(w http.ResponseWriter, req *http.Request) {
	data, err := r.deps.Services.Landmark.Nearby(req.Context(), parseFloat(req.URL.Query().Get("lat"), 0), parseFloat(req.URL.Query().Get("lng"), 0), parseFloat(req.URL.Query().Get("radius"), 5000), parseInt(req.URL.Query().Get("limit"), 20))
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, data)
}

func (r *router) handleCommentsList(w http.ResponseWriter, req *http.Request) {
	data, err := r.deps.Services.Comment.List(req.Context(), req.URL.Query().Get("targetType"), parseInt64(req.URL.Query().Get("targetId"), 0))
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, data)
}
func (r *router) handleCommentsCreate(w http.ResponseWriter, req *http.Request) {
	claims, ok := crypto.ClaimsFromContext(req.Context())
	if !ok {
		Fail(w, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}
	var body struct {
		TargetType string `json:"targetType"`
		TargetID   int64  `json:"targetId"`
		Content    string `json:"content"`
		ParentID   int64  `json:"parentId"`
	}
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	id, err := r.deps.Services.Comment.Create(req.Context(), body.TargetType, body.TargetID, body.Content, body.ParentID, claims.UserID)
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"id": id})
}
func (r *router) handleCommentsUpdate(w http.ResponseWriter, req *http.Request) {
	claims, ok := crypto.ClaimsFromContext(req.Context())
	if !ok {
		Fail(w, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	var body struct {
		Content string `json:"content"`
	}
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	actorID := commentMutationActor(claims.UserID, claims.Permissions)
	if err := r.deps.Services.Comment.Update(req.Context(), id, actorID, body.Content); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}

func commentMutationActor(userID int64, permissions []string) int64 {
	actorID := userID
	for _, permission := range permissions {
		if permission == "comment:create" || permission == "comment:update" {
			actorID = 0
			continue
		}
		if permission == "comment:delete" {
			actorID = userID
		}
	}
	return actorID
}

func (r *router) handleCommentsDelete(w http.ResponseWriter, req *http.Request) {
	claims, ok := crypto.ClaimsFromContext(req.Context())
	if !ok {
		Fail(w, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	if err := r.deps.Services.Comment.Delete(req.Context(), id, claims.UserID); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}

func (r *router) handleEventsList(w http.ResponseWriter, req *http.Request) {
	data, total, err := r.deps.Services.Event.List(req.Context(), parseNullableInt8(req.URL.Query().Get("status")), parseInt(req.URL.Query().Get("page"), 1), parseInt(req.URL.Query().Get("pageSize"), 20))
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"list": data, "total": total})
}
func (r *router) handleEventsDetail(w http.ResponseWriter, req *http.Request) {
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	data, err := r.deps.Services.Event.Detail(req.Context(), id)
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, data)
}
func (r *router) handleEventsCreate(w http.ResponseWriter, req *http.Request) {
	claims, ok := crypto.ClaimsFromContext(req.Context())
	if !ok {
		Fail(w, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}
	var body struct {
		Title           string    `json:"title"`
		RouteID         *int64    `json:"routeId"`
		MeetupAddress   string    `json:"meetupAddress"`
		MeetupLat       *float64  `json:"meetupLat"`
		MeetupLng       *float64  `json:"meetupLng"`
		MeetupTime      time.Time `json:"meetupTime"`
		StartTime       time.Time `json:"startTime"`
		EndTime         time.Time `json:"endTime"`
		MaxParticipants int       `json:"maxParticipants"`
		Fee             int       `json:"fee"`
		Description     string    `json:"description"`
	}
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	data, err := r.deps.Services.Event.Create(req.Context(), body.Title, body.RouteID, body.MeetupAddress, body.MeetupLat, body.MeetupLng, body.MeetupTime, body.StartTime, body.EndTime, body.MaxParticipants, body.Fee, body.Description, claims.UserID)
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, data)
}
func (r *router) handleEventsUpdate(w http.ResponseWriter, req *http.Request) {
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	var body struct {
		Title           string    `json:"title"`
		RouteID         *int64    `json:"routeId"`
		MeetupAddress   string    `json:"meetupAddress"`
		MeetupLat       *float64  `json:"meetupLat"`
		MeetupLng       *float64  `json:"meetupLng"`
		MeetupTime      time.Time `json:"meetupTime"`
		StartTime       time.Time `json:"startTime"`
		EndTime         time.Time `json:"endTime"`
		MaxParticipants int       `json:"maxParticipants"`
		Fee             int       `json:"fee"`
		Description     string    `json:"description"`
		Status          int8      `json:"status"`
	}
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	if err := r.deps.Services.Event.Update(req.Context(), id, body.Title, body.RouteID, body.MeetupAddress, body.MeetupLat, body.MeetupLng, body.MeetupTime, body.StartTime, body.EndTime, body.MaxParticipants, body.Fee, body.Description, body.Status); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}
func (r *router) handleEventsDelete(w http.ResponseWriter, req *http.Request) {
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	if err := r.deps.Services.Event.Delete(req.Context(), id); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}
func (r *router) handleEventsJoin(w http.ResponseWriter, req *http.Request) {
	claims, ok := crypto.ClaimsFromContext(req.Context())
	if !ok {
		Fail(w, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	var body struct {
		Name   string `json:"name"`
		Remark string `json:"remark"`
	}
	if err := decodeJSON(req, &body); err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	partID, err := r.deps.Services.Event.Join(req.Context(), id, claims.UserID, body.Name, body.Remark)
	if errors.Is(err, repository.ErrConflict) {
		OK(w, map[string]any{"id": int64(0), "queued": true})
		return
	}
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"id": partID})
}
func (r *router) handleEventsCancelJoin(w http.ResponseWriter, req *http.Request) {
	claims, ok := crypto.ClaimsFromContext(req.Context())
	if !ok {
		Fail(w, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	changed, err := r.cancelEventJoin(req, id, claims.UserID)
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true, "changed": changed})
}

func (r *router) cancelEventJoin(req *http.Request, eventID, userID int64) (bool, error) {
	if err := r.deps.Services.Event.CancelJoin(req.Context(), eventID, userID); err != nil {
		return false, err
	}
	return true, nil
}
func (r *router) handleEventsCheckIn(w http.ResponseWriter, req *http.Request) {
	claims, ok := crypto.ClaimsFromContext(req.Context())
	if !ok {
		Fail(w, http.StatusUnauthorized, 1002, "unauthorized")
		return
	}
	id, err := pathInt64(req, "id")
	if err != nil {
		Fail(w, http.StatusBadRequest, 1001, "参数错误")
		return
	}
	if err := r.deps.Services.Event.CheckIn(req.Context(), id, claims.UserID); err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, map[string]any{"success": true})
}

func (r *router) handleRecommend(w http.ResponseWriter, req *http.Request) {
	var latPtr, lngPtr *float64
	if v := req.URL.Query().Get("lat"); v != "" {
		lat := parseFloat(v, 0)
		latPtr = &lat
	}
	if v := req.URL.Query().Get("lng"); v != "" {
		lng := parseFloat(v, 0)
		lngPtr = &lng
	}
	data, err := r.deps.Services.Recommend.Home(req.Context(), 0, latPtr, lngPtr)
	if err != nil {
		failFromErr(w, err)
		return
	}
	OK(w, data)
}

func (r *router) serveStatic() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if r.deps.Cfg == nil {
			Fail(w, http.StatusNotFound, 1004, "not found")
			return
		}
		http.StripPrefix("/static/", http.FileServer(http.Dir(r.deps.Cfg.DataDir))).ServeHTTP(w, req)
	}
}

func decodeJSON(req *http.Request, dest any) error {
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dest)
}

func clientIP(req *http.Request) string {
	if xff := req.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xr := req.Header.Get("X-Real-IP"); xr != "" {
		return xr
	}
	host, _, found := strings.Cut(req.RemoteAddr, ":")
	if !found {
		return req.RemoteAddr
	}
	if host == "" {
		return req.RemoteAddr
	}
	return host
}

func bearer(req *http.Request) string {
	value := req.Header.Get("Authorization")
	return strings.TrimPrefix(value, "Bearer ")
}

func pathInt64(req *http.Request, key string) (int64, error) {
	return strconv.ParseInt(req.PathValue(key), 10, 64)
}
func parseInt(value string, fallback int) int {
	if v, err := strconv.Atoi(value); err == nil {
		return v
	}
	return fallback
}
func parseInt64(value string, fallback int64) int64 {
	if v, err := strconv.ParseInt(value, 10, 64); err == nil {
		return v
	}
	return fallback
}
func parseFloat(value string, fallback float64) float64 {
	if v, err := strconv.ParseFloat(value, 64); err == nil {
		return v
	}
	return fallback
}
func parseNullableInt8(value string) *int8 {
	if value == "" {
		return nil
	}
	if v, err := strconv.ParseInt(value, 10, 8); err == nil {
		vv := int8(v)
		return &vv
	}
	return nil
}
