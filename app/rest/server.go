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
	authenticator auth.Service
	http          *http.Server
	lock          sync.Mutex
}

// Run starts the web-server for listening
func (s *Rest) Run(port int) {
	s.authenticator = *s.makeAuth()

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
		URL:            strings.TrimSuffix(s.ServiceURL, "/"),
		Issuer:         s.AppName,
		TokenDuration:  s.Auth.TTL.JWT,
		CookieDuration: s.Auth.TTL.Cookie,
		SecureCookies:  strings.HasPrefix(s.ServiceURL, "https://"),
		AvatarStore:    avatar.NewNoOp(),
		JWTQuery:       "jwt",
		Logger:         logger.Std,
		SecretReader: token.SecretFunc(func(_ string) (string, error) {
			// todo is thread-safe?
			return s.JWTSecret, nil
		}),
		// c.User.Audience - address of front end,
		ClaimsUpd: token.ClaimsUpdFunc(func(c token.Claims) token.Claims {
			if c.User == nil {
				return c
			}
			uInfo, err := s.UserService.GetBasicUserInfo(c.User.ID)
			if err != nil {
				log.Printf("[WARN] failed to recognize is user admin, id: %s, error: %s", c.User.ID, err.Error())
				return c
			}
			privs := []string{}
			for k, v := range uInfo.Privileges {
				if v {
					privs = append(privs, k)
				}
			}
			c.User.SetSliceAttr("privs", privs)
			return c
		}),
		Validator: token.ValidatorFunc(func(token string, claims token.Claims) bool {
			// todo do we need validator?
			return claims.User != nil
		}),
	})
	authenticator.AddDirectProvider("local", provider.CredCheckerFunc(s.UserService.CheckUserCredentials))
	return authenticator
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

	authHandler, _ := s.authenticator.Handlers()

	r.Group(func(r chi.Router) {
		r.Use(middleware.Timeout(5 * time.Second))
		r.Mount("/auth", authHandler)
	})

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("[DEBUG] registered route: %s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		log.Printf("[WARN] error occurred while printing routes: %s", err.Error())
	}

	r.Group(func(r chi.Router) {
		// protected routes

	})

	r.Group(func(r chi.Router) {
		// public routes
	})

	return r
}
