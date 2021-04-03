package middleware

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mivinci/pi"
)

func Log(w io.Writer) pi.Middleware {
	// log.SetOutput(w)
	log.SetOutput(w)
	log.SetPrefix("[INFO] ")
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			t := time.Now()
			next(w, r)
			log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(t))
		}
	}
}
