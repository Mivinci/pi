package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/mivinci/pi"
)

func Timeout(d time.Duration) pi.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer func() {
				cancel()
				if ctx.Err() == context.DeadlineExceeded {
					w.WriteHeader(http.StatusGatewayTimeout)
				}
			}()
			next(w, r.WithContext(ctx))
		}
	}
}
