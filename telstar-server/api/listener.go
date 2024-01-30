package api

// See https://github.com/go-chi/chi/blob/master/_examples/rest/main.go for details

//
// Client requests:
// ----------------
// $ curl http://localhost:3333/
// $ curl http://localhost:3333/articles
// $ curl http://localhost:3333/articles/1
// $ curl -X DELETE http://localhost:3333/articles/1
// $ curl -X POST -d '{"id":"will-be-omitted","title":"awesomeness"}' http://localhost:3333/articles
// $ curl http://localhost:3333/articles/97
// $ curl http://localhost:3333/articles

import (
	"bitbucket.org/johnnewcombe/telstar-library/logger"
	"bitbucket.org/johnnewcombe/telstar/config"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"net/http"
	"time"
)

func Start(apiPort int, settings config.Config) error {

	logger.LogInfo.Printf("Starting Restful API Server %s on port %d", settings.Server.DisplayName, apiPort)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", homeHandler)
	r.Get("/ping", pingHandler)
	r.Get("/status", HandlerWrapper(http.HandlerFunc(getStatus), settings))
	r.Put("/login", HandlerWrapper(http.HandlerFunc(putLogin), settings))

	// Mount the frame sub-router (frame and frames for compatibility with Telstar 0.x and 1.x
	r.Mount("/frame", frameRouter(settings))
	r.Mount("/frames", frameRouter(settings))

	// Mount the user sub-router
	r.Mount("/user", userRouter(settings))

	return http.ListenAndServe(fmt.Sprintf(":%d", apiPort), r)
}

func createJwtToken(userId string, expiryMinutes time.Duration) (string, error) {

	var (
		tokenString string
		err         error
	)

	claims := map[string]interface{}{"user-id": userId}

	jwtauth.SetExpiry(claims, time.Now().Add(time.Minute*expiryMinutes))

	// jwt token with claims = `user_id:nnn`
	if _, tokenString, err = TokenAuth.Encode(claims); err != nil {
		return "", err
	}

	return tokenString, nil
}
