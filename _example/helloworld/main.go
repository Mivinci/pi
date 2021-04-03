package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/mivinci/pi"
	"github.com/mivinci/pi/middleware"
)

type Foo struct{}

func (Foo) Get(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "This is a GET request")
}

func main() {
	pi.RegisterAndRun(":8080", new(Foo), middleware.Log(os.Stdout))
}
