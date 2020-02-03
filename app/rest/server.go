package rest

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	R "github.com/go-pkgz/rest"
)

// Rest defines a simple web server for routing to REST api methods
type Rest struct {
	Version   string
	AppName   string
	AppAuthor string
	JWTSecret string
	// Data services

	// URLs, and other configs for running web-server
	jwtTokenAuth *jwtauth.JWTAuth

	// Private fields (http object, etc.)
	http *http.Server
	lock sync.Mutex
}

// Run starts the web-server for listening
func (s *Rest) Run(port int) {
	s.lock.Lock()
	s.http = s.makeHTTPServer(port, s.routes())
	s.http.ErrorLog = log.New(os.Stdout, "", log.Flags())
	s.lock.Unlock()
	log.Printf("[INFO] started web server at port %d", port)
	err := s.http.ListenAndServe()
	log.Printf("[WARN] web server terminated reason: %s", err)
}

func (s *Rest) makeHTTPServer(port int, routes chi.Router) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           routes,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
}

func (s *Rest) routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(R.AppInfo(s.AppName, s.AppAuthor, s.Version), R.Ping)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(s.jwtTokenAuth))
		r.Use(jwtauth.Authenticator)

		// protected routes

	})

	r.Group(func(r chi.Router) {
		// public routes
	})

	return r
}
