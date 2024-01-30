package api

import (
	"bitbucket.org/johnnewcombe/telstar/config"
	"context"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/lestrrat-go/jwx/jwt"
	"net/http"
)

var (
	TokenAuth *jwtauth.JWTAuth
)

type ContextData struct {
	Settings config.Config
}

// this should get called before main (telstar.go)
func init() {
	/*
		// this is the env for Telstar v2.0 which uses HS256 jwt tokens
		apiSecret := os.Getenv(ENV_API_SECRET)

		// this is the env for early versions of Telstar which uses a encrypted abstract token
		cookieSecret := os.Getenv(ENV_COOKIE_SECRET)

		if len(apiSecret) == 0 {
			// support for env vars used in Telstar 0.x and 1.x
			if len(cookieSecret) == 0 {
				log.Fatal("The Telstar API is unable to start as the TELSTAR_AUTH_SECRET environment variable has not been set. \r\nIdeally, his should be set using at least a 32 character string. The string should be kept secret!")
			} else {
				logger.LogWarn.Println("The TELSTAR_AUTH_SECRET has not been set. Using TELSTAR_COOKIE_SECRET instead.")
				apiSecret = cookieSecret
			}
		}
		TokenAuth = jwtauth.New("HS256", []byte(apiSecret), nil)

	*/
}

// Authenticator is authentication middleware to enforce access from the
// Verifier middleware request context values. The Authenticator sends a 401 Unauthorized
// response for any unverified tokens and passes the good ones through.
// This is a custom authenticator based on the jwtauth.Authenticator method.
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// claims can be retrieved from context, see handlers
		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil {
			render.Render(w, r, ErrUnauthorizedRequest(err))
			return
		}
		if token == nil || jwt.Validate(token) != nil {
			render.Render(w, r, ErrUnauthorizedRequest(err))
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

// AdminOnly middleware restricts access to just administrators.
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value("acl.admin").(bool)
		if !ok || !isAdmin {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AdminOnly middleware restricts access to just administrators.
func LoggedInOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value("acl.admin").(bool)
		if !ok || !isAdmin {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// HandlerWrapper is a wrapper designed to wrap an http.Handler, the wrapper allows the settings to be passed
// These are placed in context before calling the actualk handler. This is used as follows...
// Instead of using a basic handler like this;
//
//          r.Get("/", getFrame)
//
// We can leverage the wrapper as follows;
//
//          r.Get("/", HandlerWrapper(http.HandlerFunc(getFrame), settings))
//
// This second approach places settings into context so that the actual handler
// can get access to them.
func HandlerWrapper(h http.Handler, settings config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// add settings to the context and add it back to the request
		ctxData := ContextData{settings}
		ctx := context.WithValue(r.Context(), "ctx-data", &ctxData)

		// call the handler we are wrapping
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
