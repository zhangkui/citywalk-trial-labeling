package verify

import (
	"net/http"
	"strings"
	"testing"
)

func TestBug004_BusinessRegression(t *testing.T) {
	env := newVerifyEnv(t)
	userID, token := env.user("story:edit")
	res, _ := env.db.Exec(`INSERT INTO stories(title,content,created_by,status) VALUES ('old story','old content',?,1)`, userID)
	storyID, _ := res.LastInsertId()
	_, _ = env.db.Exec(`INSERT INTO story_media(story_id,type,url,`+"`order`"+`) VALUES (?,1,'old.jpg',0)`, storyID)
	result := env.api(http.MethodPut, "/api/v1/stories/"+fmtID(storyID), token, map[string]any{
		"title": "new story", "content": "new content", "status": 1,
		"media": []map[string]any{{"type": 1, "url": "new.jpg"}, {"type": 1, "url": strings.Repeat("x", 300)}},
	})
	if result.Status < 400 {
		t.Fatalf("invalid media replacement unexpectedly succeeded: %v", result.Body)
	}
	var title string
	var oldMedia int
	_ = env.db.QueryRow(`SELECT title FROM stories WHERE id=?`, storyID).Scan(&title)
	_ = env.db.QueryRow(`SELECT COUNT(1) FROM story_media WHERE story_id=? AND url='old.jpg'`, storyID).Scan(&oldMedia)
	if title != "old story" || oldMedia != 1 {
		t.Errorf("story update partially committed: title=%q old_media=%d", title, oldMedia)
	}
	res2, _ := env.db.Exec(`INSERT INTO stories(title,content,created_by,status) VALUES ('ordered story','body',?,1)`, userID)
	orderedID, _ := res2.LastInsertId()
	requireStatus(t, env.api(http.MethodPut, "/api/v1/stories/"+fmtID(orderedID), token, map[string]any{"title": "ordered story", "content": "body", "status": 1, "media": []map[string]any{{"type": 1, "url": "late.jpg", "order": 8}, {"type": 1, "url": "early.jpg", "order": 2}}}), http.StatusOK)
	rows, _ := env.db.Query(`SELECT url,`+"`order`"+` FROM story_media WHERE story_id=?`, orderedID)
	defer rows.Close()
	orders := map[string]int{}
	for rows.Next() {
		var url string
		var order int
		_ = rows.Scan(&url, &order)
		orders[url] = order
	}
	if orders["late.jpg"] != 8 || orders["early.jpg"] != 2 {
		t.Errorf("media order was discarded: %v", orders)
	}
}
