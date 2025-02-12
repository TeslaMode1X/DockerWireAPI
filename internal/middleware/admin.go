package middleware

import "net/http"

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value("role").(float64)
		if !ok || role != 1 {
			permissionDenied(w, r, "admin access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}
