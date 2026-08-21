package middleware

import "net/http"

func RequireRole(next http.HandlerFunc, have func(*http.Request) []string, need []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		has := have(r)
		if !containsAny(has, need) { http.Error(w, "forbidden", http.StatusForbidden); return }
		next(w, r)
	}
}

func containsAny(have, need []string) bool {
	set := map[string]struct{}{}
	for _, v := range have { set[v] = struct{}{} }
	for _, v := range need { if _, ok := set[v]; ok { return true } }
	return false
}

