package storage

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Local struct {
	BaseDir string
	BaseURL string
}

func NewLocal(baseDir, baseURL string) *Local {
	_ = os.MkdirAll(baseDir, 0o755)
	return &Local{BaseDir: baseDir, BaseURL: baseURL}
}

func (l *Local) Save(name string, reader io.Reader) (string, error) {
	fileName := sanitizeName(name)
	full := filepath.Join(l.BaseDir, time.Now().Format("20060102"), fileName)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return "", err
	}
	file, err := os.Create(full)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if _, err := io.Copy(file, reader); err != nil {
		return "", err
	}
	return l.PublicURL(full), nil
}

func (l *Local) PublicURL(path string) string {
	clean := strings.TrimPrefix(filepath.ToSlash(path), filepath.ToSlash(l.BaseDir)+"/")
	return strings.TrimRight(l.BaseURL, "/") + "/" + clean
}

func sanitizeName(name string) string {
	name = filepath.Base(name)
	name = strings.ReplaceAll(name, " ", "_")
	return name
}

