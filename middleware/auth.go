package middleware

import (
	"hash"
	"net/http"

	"github.com/mivinci/pi"
)

func Auth(key string, h hash.Hash) pi.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

		}
	}
}
