package rest

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-pkgz/auth/logger"

	"github.com/go-pkgz/auth/token"

	"github.com/go-pkgz/auth/avatar"

	"github.com/go-pkgz/auth/provider"

	"github.com/Semior001/gotemplate/app/store/user"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-pkgz/auth"
	R "github.com/go-pkgz/rest"
)

// Rest defines a simple web server for routing to REST api methods
type Rest struct {
	Version    string
	AppName    string
	AppAuthor  string
	JWTSecret  string
	ServiceURL string

	// Data services
	UserService user.Service

	Auth struct {
		TTL struct {
			JWT    time.Duration
			Cookie time.Duration
		}
	}

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

func (s *Rest) makeAuth() *auth.Service {
	authenticator := auth.NewService(auth.Opts{
		TokenDuration:  s.Auth.TTL.JWT,
		CookieDuration: s.Auth.TTL.Cookie,
		Issuer:         s.AppName,
		SecureCookies:  strings.HasPrefix(s.ServiceURL, "https://"),
		URL:            strings.TrimSuffix(s.ServiceURL, "/"),
		AvatarStore:    avatar.NewNoOp(),
		JWTQuery:       "jwt",
		Logger:         logger.Std,
		SecretReader: token.SecretFunc(func(_ string) (string, error) { // secret key for JWT
			// todo is thread-safe?
			return s.JWTSecret, nil
		}),

		// c.User.Audience - address of front end,

		ClaimsUpd: token.ClaimsUpdFunc(func(c token.Claims) token.Claims {
			if c.User == nil {
				return c
			}
			adm, err := s.UserService.IsAdmin(c.User.ID)
			if err != nil {
				log.Printf("[WARN] failed to recognize is user admin, id: %s", c.User.ID)
				return c
			}
			c.User.SetAdmin(adm)
			return c
		}),
		Validator: token.ValidatorFunc(func(_ string, claims token.Claims) bool {
			// allow only dev_* names
			return claims.User != nil && strings.HasPrefix(claims.User.Name, "dev_")
		}),
	})
	// adding custom credentials checker
	authenticator.AddDirectProvider("local", provider.CredCheckerFunc(s.UserService.CheckUserCredentials))
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
		// protected routes

	})

	r.Group(func(r chi.Router) {
		// public routes
	})

	return r
}
