package middleware

import (
	"net/http"
	"time"
)

type Bucket struct {
	window time.Duration
	limit  int
}

func NewBucket(window time.Duration, limit int) Bucket { return Bucket{window: window, limit: limit} }

func (b Bucket) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// placeholder for shared limiter integration
		next(w, r)
	}
}

