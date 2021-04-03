package pi

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	typeResponseWriter = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	typeRequest        = reflect.TypeOf((*http.Request)(nil))
)

var ErrMethodNotFound = errors.New("method not found")

type prototype struct {
	method reflect.Method
	count  uint64
}

func (p prototype) Count() uint64 {
	return atomic.LoadUint64(&p.count)
}

type service struct {
	name    string
	value   reflect.Value
	methods map[string]*prototype
}

func newService() *service {
	return &service{methods: make(map[string]*prototype)}
}

func (s *service) register(o interface{}, name string) {
	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)
	if name == "" {
		name = reflect.Indirect(v).Type().Name()
	}
	s.value = v
	s.name = strings.ToLower(name)
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		methodType := method.Type
		// skip unexported methods
		if method.PkgPath != "" {
			continue
		}
		if methodType.NumIn() != 3 || methodType.NumOut() != 0 {
			continue
		}
		if rw := methodType.In(1); rw != typeResponseWriter {
			log.Printf("pi: the first argument of method %s.%s must be an http.ResponseWriter\n", name, method.Name)
		}
		if req := methodType.In(2); req != typeRequest {
			log.Printf("pi: the second argument of method %s.%s must be an http.Request\n", name, method.Name)
		}
		s.methods[strings.ToUpper(method.Name)] = &prototype{method: method}
		log.Printf("pi: registered %s.%s\n", name, method.Name)
	}
}

func (s *service) Call(name string, rw, req reflect.Value) error {
	proto, ok := s.methods[name]
	if !ok {
		return ErrMethodNotFound
	}
	proto.method.Func.Call([]reflect.Value{s.value, rw, req})
	return nil
}

type Server struct {
	options  Options
	services sync.Map
}

func New(opts ...Option) *Server {
	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}
	return &Server{options: o}
}

func (s *Server) RegisterName(name string, o interface{}) error {
	svc := newService()
	svc.register(o, name)
	if _, dup := s.services.LoadOrStore(path.Join(s.options.group, svc.name), svc); dup {
		return fmt.Errorf("pi: service %s is already registered", svc.name)
	}
	return nil
}

func (s *Server) Register(o interface{}) error {
	return s.RegisterName("", o)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.EscapedPath()[1:]
	v, ok := s.services.Load(name)
	if !ok {
		http.NotFound(w, r)
		return
	}
	svc, ok := v.(*service)
	if !ok {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	endpoint := func(w http.ResponseWriter, r *http.Request) {
		if err := svc.Call(r.Method, reflect.ValueOf(w), reflect.ValueOf(r)); err != nil {
			http.Error(w, err.Error(), errCode(err))
		}
	}
	chain(endpoint, s.options.chain...)(w, r)
}

type Middleware func(next http.HandlerFunc) http.HandlerFunc

// chain is the magic.
func chain(endpoint http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		endpoint = middleware[i](endpoint)
	}
	return endpoint
}

func errCode(err error) int {
	if errors.Is(err, ErrMethodNotFound) {
		return http.StatusMethodNotAllowed
	}
	return http.StatusInternalServerError
}

func RegisterAndRun(addr string, o interface{}, chain ...Middleware) {
	h := New(Chain(chain...))
	if err := h.Register(o); err != nil {
		log.Fatal("pi:", err)
	}
	log.Fatal("pi:", http.ListenAndServe(addr, h))
}
