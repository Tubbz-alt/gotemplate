package rest

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-pkgz/rest"
)

// JwtAuthRouter provides router for all requests for authentication
type JwtAuthRouter struct {
	checkUserCredentials func(rest.JSON) (authenticated bool, claims jwt.MapClaims)
}

// POST /auth
//func (j *JwtAuthRouter) authenticateCtrl(w http.ResponseWriter, r *http.Request) {
//	decoder := json.NewDecoder(r.Body)
//	var t rest.JSON
//	err := decoder.Decode(&t)
//	if err != nil {
//		log.Printf("[ERROR] failed to decode json in request")
//		SendErrorJSON(w, r, http.StatusInternalServerError, merr, "can't decode JSON body", rest.ErrInternal)
//	}
//	authStatus, claims := j.checkUserCredentials(t)
//}
