package verify

import (
	"fmt"
	"net/http"
	"testing"
)

func TestBug003_BusinessRegression(t *testing.T) {
	env := newVerifyEnv(t)
	ownerID, _ := env.user()
	city := fmt.Sprintf("VerifyCity-%d", ownerID)
	var themeID int64
	if err := env.db.QueryRow(`SELECT id FROM themes ORDER BY id LIMIT 1`).Scan(&themeID); err != nil {
		t.Fatal(err)
	}
	_, _ = env.db.Exec(`INSERT INTO routes(title,city,theme_id,rating,created_by,status) VALUES ('high',?,?,4.8,?,2),('low',?,?,2.0,?,2)`, city, themeID, ownerID, city, themeID, ownerID)
	result := env.api(http.MethodGet, "/api/v1/routes?city="+city+"&themeId="+fmtID(themeID)+"&minRate=4&page=1&pageSize=20", "", nil)
	requireStatus(t, result, http.StatusOK)
	data := result.Body["data"].(map[string]any)
	list := data["list"].([]any)
	total := int(data["total"].(float64))
	if len(list) != 1 || total != 1 {
		t.Errorf("filtered list and total diverged: list=%d total=%d", len(list), total)
	}
	invalid := env.api(http.MethodGet, "/api/v1/routes?minRate=not-a-number", "", nil)
	if invalid.Status != http.StatusBadRequest {
		t.Errorf("invalid minRate should be rejected, got HTTP %d", invalid.Status)
	}
}
