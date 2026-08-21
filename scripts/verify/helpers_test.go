package verify

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"citywalk/internal/app"
	"citywalk/internal/config"
	appcrypto "citywalk/internal/platform/crypto"

	_ "github.com/go-sql-driver/mysql"
)

type verifyEnv struct {
	t      *testing.T
	db     *sql.DB
	server *httptest.Server
	jwt    *appcrypto.JWT
}

type apiResult struct {
	Status int
	Body   map[string]any
}

var verifyOnce sync.Once
var fixtureSequence atomic.Int64

func newVerifyEnv(t *testing.T) *verifyEnv {
	t.Helper()
	verifyOnce.Do(func() {
		_, file, _, _ := runtime.Caller(0)
		root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
		if err := os.Chdir(root); err != nil {
			panic(err)
		}
	})
	cfg := &config.Config{
		Environment: "test", ListenAddr: ":0", ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second,
		IdleTimeout: 60 * time.Second, ShutdownAfter: 5 * time.Second,
		MySQLHost: "127.0.0.1", MySQLPort: 13306, MySQLDatabase: "citywalk", MySQLUser: "citywalk", MySQLPassword: "citywalk2024",
		RedisHost: "127.0.0.1", RedisPort: 16379, JWTSecret: "trial-verification-secret", AdminUsername: "admin", AdminPassword: "Admin123!",
		DataDir: t.TempDir(), CORSAllowOrigin: "*",
	}
	application, err := app.New(cfg, nil)
	if err != nil {
		t.Fatalf("start application: %v", err)
	}
	database, err := sql.Open("mysql", cfg.MySQLDSN())
	if err != nil {
		_ = application.Close()
		t.Fatalf("open verification database: %v", err)
	}
	env := &verifyEnv{t: t, db: database, server: httptest.NewServer(application.Router()), jwt: appcrypto.NewJWT(cfg.JWTSecret)}
	t.Cleanup(func() {
		env.server.Close()
		_ = env.db.Close()
		_ = application.Close()
	})
	return env
}

func (e *verifyEnv) api(method, path, token string, body any) apiResult {
	e.t.Helper()
	var payload bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&payload).Encode(body); err != nil {
			e.t.Fatalf("encode request: %v", err)
		}
	}
	req, err := http.NewRequestWithContext(context.Background(), method, e.server.URL+path, &payload)
	if err != nil {
		e.t.Fatalf("create request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		e.t.Fatalf("request %s %s: %v", method, path, err)
	}
	defer resp.Body.Close()
	result := apiResult{Status: resp.StatusCode, Body: map[string]any{}}
	_ = json.NewDecoder(resp.Body).Decode(&result.Body)
	return result
}

func (e *verifyEnv) user(perms ...string) (int64, string) {
	e.t.Helper()
	name := fmt.Sprintf("verify_%d_%d", time.Now().UnixNano(), fixtureSequence.Add(1))
	hash, err := appcrypto.HashPassword("Password123!")
	if err != nil {
		e.t.Fatalf("hash password: %v", err)
	}
	res, err := e.db.Exec(`INSERT INTO users(username, password_hash, nickname, status) VALUES (?, ?, ?, 1)`, name, hash, name)
	if err != nil {
		e.t.Fatalf("create user: %v", err)
	}
	id, _ := res.LastInsertId()
	token, err := e.jwt.Issue(id, name, []string{"verification"}, perms, appcrypto.TokenTypeAccess, time.Hour)
	if err != nil {
		e.t.Fatalf("issue token: %v", err)
	}
	return id, token
}

func requireStatus(t *testing.T, result apiResult, expected int) {
	t.Helper()
	if result.Status != expected {
		t.Fatalf("expected HTTP %d, got %d body=%v", expected, result.Status, result.Body)
	}
}
