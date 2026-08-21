package verify

import (
	"net/http"
	"testing"
)

func TestBug002_BusinessRegression(t *testing.T) {
	env := newVerifyEnv(t)
	userID, token := env.user("route:favorite")
	ownerID, _ := env.user()
	res, _ := env.db.Exec(`INSERT INTO routes(title,created_by,status,favorite_count) VALUES ('favorite target',?,2,0)`, ownerID)
	routeID, _ := res.LastInsertId()
	_, _ = env.db.Exec(`INSERT INTO favorites(user_id,target_type,target_id) VALUES (?,'route',?)`, userID, routeID)
	t.Run("reconcile stale aggregate", func(t *testing.T) {
		requireStatus(t, env.api(http.MethodPost, "/api/v1/routes/"+fmtID(routeID)+"/favorite", token, map[string]any{}), http.StatusOK)
		var count int
		_ = env.db.QueryRow(`SELECT favorite_count FROM routes WHERE id=?`, routeID).Scan(&count)
		if count != 1 {
			t.Errorf("favorite relation exists but display state was not reconciled: %d", count)
		}
	})
	t.Run("duplicate request stays idempotent", func(t *testing.T) {
		_, _ = env.db.Exec(`UPDATE routes SET favorite_count=1 WHERE id=?`, routeID)
		requireStatus(t, env.api(http.MethodPost, "/api/v1/routes/"+fmtID(routeID)+"/favorite", token, map[string]any{}), http.StatusOK)
		var count, relations int
		_ = env.db.QueryRow(`SELECT favorite_count FROM routes WHERE id=?`, routeID).Scan(&count)
		_ = env.db.QueryRow(`SELECT COUNT(1) FROM favorites WHERE user_id=? AND target_type='route' AND target_id=?`, userID, routeID).Scan(&relations)
		if count != 1 || relations != 1 {
			t.Errorf("duplicate favorite changed persisted state: count=%d relations=%d", count, relations)
		}
	})
}
