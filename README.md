# pi
Making a web server is easy as π.


## Example
```go
type Foo struct{}

func (Foo) Get(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "This is a GET request")
}

func main() {
    pi.RegisterAndRun(":8080", new(Foo), middleware.Log(os.Stdout))
}
```
