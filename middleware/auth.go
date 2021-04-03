package middleware

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/mivinci/pi"
	"hash"
	"net/http"
)

var (
	ErrTokenFormat = errors.New("illegal token format")
	ErrTokenSign   = errors.New("incorrect token signature")
)

type authKey struct{}

func AuthFromContext(ctx context.Context) (MD, bool) {
	return FromContext(ctx, authKey{})
}

func AuthToken(b []byte, key string, h hash.Hash) string {
	raw := base64.URLEncoding.EncodeToString(b)
	sign := signature(b, key, h)
	return raw + "." + sign
}

func Auth(header, key string, h hash.Hash) pi.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get(header)
			_, md, err := verify(token, key, h)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := WithMetadata(r.Context(), authKey{}, md)
			next(w, r.WithContext(ctx))
		}
	}
}

func verify(token, key string, h hash.Hash) (raw string, md MD, err error) {
	var sign string
	var b []byte
	raw, sign, err = split(token)
	if err != nil {
		return
	}
	b, err = base64.URLEncoding.DecodeString(raw)
	if err != nil {
		return
	}
	raw = string(b)
	if sign != signature(b, key, h) {
		err = ErrTokenSign
		return
	}
	err = json.Unmarshal(b, &md)
	return
}

func signature(data []byte, key string, h hash.Hash) string {
	h.Reset()
	h.Write(data)
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}

func split(token string) (string, string, error) {
	i := 0
	for ; i < len(token) && token[i] != '.'; i++ {
	}
	if i >= len(token) {
		return "", "", ErrTokenFormat
	}
	return token[:i], token[i+1:], nil
}
