package verify

import (
	"fmt"
	"net/http"
	"testing"
)

func TestBug001_BusinessRegression(t *testing.T) {
	env := newVerifyEnv(t)
	userID, token := env.user("route:edit")
	t.Run("rollback aggregate", func(t *testing.T) {
		res, err := env.db.Exec(`INSERT INTO routes(title, description, created_by, status) VALUES ('old title','old body',?,1)`, userID)
		if err != nil {
			t.Fatal(err)
		}
		routeID, _ := res.LastInsertId()
		_, _ = env.db.Exec(`INSERT INTO waypoints(route_id,name,lat,lng,stay_duration,`+"`order`"+`) VALUES (?,'old stop',31.2,121.4,10,0)`, routeID)
		result := env.api(http.MethodPut, "/api/v1/routes/"+fmtID(routeID), token, map[string]any{"title": "new title", "description": "new body", "status": 1, "waypoints": []map[string]any{{"name": "new stop", "lat": nil, "lng": 121.5, "stayDuration": 5, "order": 4}}})
		if result.Status < 400 {
			t.Fatalf("invalid waypoint update unexpectedly succeeded: %v", result.Body)
		}
		var title string
		var waypointCount int
		_ = env.db.QueryRow(`SELECT title FROM routes WHERE id=?`, routeID).Scan(&title)
		_ = env.db.QueryRow(`SELECT COUNT(1) FROM waypoints WHERE route_id=? AND name='old stop'`, routeID).Scan(&waypointCount)
		if title != "old title" || waypointCount != 1 {
			t.Errorf("failed update partially committed: title=%q old_waypoints=%d", title, waypointCount)
		}
	})
	t.Run("preserve waypoint order", func(t *testing.T) {
		res, _ := env.db.Exec(`INSERT INTO routes(title,created_by,status) VALUES ('ordered',?,1)`, userID)
		routeID, _ := res.LastInsertId()
		requireStatus(t, env.api(http.MethodPut, "/api/v1/routes/"+fmtID(routeID), token, map[string]any{"title": "ordered", "status": 1, "waypoints": []map[string]any{{"name": "late", "lat": 31.1, "lng": 121.1, "order": 9}, {"name": "early", "lat": 31.2, "lng": 121.2, "order": 2}}}), http.StatusOK)
		rows, _ := env.db.Query(`SELECT name,`+"`order`"+` FROM waypoints WHERE route_id=? ORDER BY `+"`order`", routeID)
		defer rows.Close()
		orders := map[string]int{}
		for rows.Next() {
			var name string
			var order int
			_ = rows.Scan(&name, &order)
			orders[name] = order
		}
		if orders["late"] != 9 || orders["early"] != 2 {
			t.Errorf("client waypoint order was discarded: %v", orders)
		}
	})
}

func fmtID(id int64) string { return fmt.Sprintf("%d", id) }
