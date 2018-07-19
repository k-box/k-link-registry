package klinkregistry

import (
	"context"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/pkg/errors"
)

// Errors that are returned by the session manager
var (
	ErrUnsupportedSignature = errors.New("Unsupported signature")
)

// SessionsProvider is implemented by all providers that handle
// sessions
type SessionsProvider interface {
	Middleware(http.Handler) http.Handler
	RequireAuthorized(http.Handler) http.Handler
	GetUser(req *http.Request) *User
	GenerateToken(u User, expires time.Time) (string, error)
}

// claims are the JWT claims.
type claims struct {
	User

	jwt.StandardClaims
}

// User is the user information inside a session
type User struct {
	ID          int64  `json:"id"`
	Role        string `json:"role"`
	DisplayName string `json:"name"`
}

// ContextKey is a custom type for providing keys to Context.Value. It is
// a custom datatype to prevent collisions, since other packages may also
// provide context values with the same key names.
type ContextKey string

// UserKey is the Key used for user objects stored inside Context.Value
var UserKey ContextKey = "user"

// Middleware provides a middleware that generates a session from the user
// request, and stores it in the request context for later use. This way
// subsequent middlewares or handlers can re-use the session object, without
// a performance penalty.
func (s JWTSession) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// keyFunc returns the correct key for the token
		keyFunc := func(token *jwt.Token) (interface{}, error) {
			switch token.Method {
			case jwt.SigningMethodHS256:
				return s.Key, nil
			default:
				return nil, ErrUnsupportedSignature
			}
		}

		// get token from auth header
		tokenExtractor := request.AuthorizationHeaderExtractor

		// serialize the claims in our custom claims format
		var claims claims
		withClaims := request.WithClaims(&claims)

		token, err := request.ParseFromRequest(req, tokenExtractor, keyFunc, withClaims)
		if err != nil || !token.Valid {
			// token is invalid, just continue like we didn't find a token.
			next.ServeHTTP(w, req)
			return
		}

		u := claims.User

		// extract context from old request, augment context with custom values
		ctx := context.WithValue(req.Context(), UserKey, u)
		req = req.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}

// JWTSession is a session Provider that uses JWT tokens
type JWTSession struct {
	Key []byte
}

// RequireAuthorized requires the request to have an authorized user
func (s JWTSession) RequireAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		user := s.GetUser(req)
		if user == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, req)
	})
}

// GetUser returns a user from a request, if no authorized user was detected,
// the returned user will be nil
func (s JWTSession) GetUser(req *http.Request) *User {
	ctx := req.Context()
	value := ctx.Value(UserKey)
	if value == nil {
		// return nil user on nil value, since we cannot cast nil
		return nil
	}
	u := value.(User)
	return &u
}

// GenerateToken returns a valid token that proves the user has access until
// the set expiry time
func (s JWTSession) GenerateToken(u User, expires time.Time) (string, error) {
	// populate claims
	claims := claims{
		u,
		jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
			Issuer:    "Registry",
		},
	}

	// create unsigned token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign token
	return token.SignedString(s.Key)
}
