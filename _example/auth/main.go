package auth

import (
	"crypto/sha256"
	"fmt"
	"github.com/mivinci/pi"
	"github.com/mivinci/pi/middleware"
	"net/http"
)

const (
	reqHeader = "Token"
	secretKey = "secretKey"
)

type User struct{}

func (User) Get(w http.ResponseWriter, r *http.Request) {
	auth, _ := middleware.AuthFromContext(r.Context())
	_, _ = fmt.Fprint(w, auth["user_id"])
}

func main() {
	auth := middleware.Auth(reqHeader, secretKey, sha256.New())
	pi.RegisterAndRun(":8080", new(User), auth)
}
