package middleware

import (
	"context"
	"github.com/mivinci/pi"
	"net/http"
)

type MD map[string]interface{}

func FromContext(ctx context.Context, key interface{}) (md MD, ok bool) {
	md, ok = ctx.Value(key).(MD)
	return
}

func WithMetadata(ctx context.Context, key interface{}, md MD) context.Context {
	return context.WithValue(ctx, key, md)
}

type DefaultKey struct{}

func Metadata(md MD) pi.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := WithMetadata(r.Context(), DefaultKey{}, md)
			next(w, r.WithContext(ctx))
		}
	}
}
