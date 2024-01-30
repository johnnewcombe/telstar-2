package api

import (
	"bitbucket.org/johnnewcombe/telstar/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"net/http"
)

/*
Process is as follows

1 - request to /articles/{id} invokes the 'GetArticle' handler with id is in ctx
2 - ArticleContext middleware extracts the ID and calls the appropriate dal function to get the data and place it in ctx
3 - The handler gets the article from ctx and renders the output of NewArticleResponse(article) function
4 - The NewNewArticleResponse function calls the db functions to get user details bassed on the user idf in the article
    then adds the result NewUserPayloadResponse(user) to the previously obtained result of NewArticleResponse(article)

*/

func frameRouter(settings config.Config) chi.Router {

	r := chi.NewRouter()

	// Seek, verify and validate JWT tokens
	r.Use(jwtauth.Verifier(TokenAuth))

	// Handle valid / invalid tokens. This is a custom authenticator
	// based on the jwtauth.Authenticator method.
	r.Use(Authenticator)

	r.Get("/", HandlerWrapper(http.HandlerFunc(getFrames), settings))
	r.Put("/", HandlerWrapper(http.HandlerFunc(updateFrame), settings))

	r.Route("/{pageId:^[0-9]+[a-z]$}", func(r chi.Router) {
		// HandlerWrapper is a custom wrapper that returns a Handler,
		// it allows settings to be passed to the handler
		r.Get("/", HandlerWrapper(http.HandlerFunc(getFrame), settings))
		r.Delete("/", HandlerWrapper(http.HandlerFunc(deleteFrame), settings))
	})

	return r

}

func userRouter(settings config.Config) chi.Router {

	r := chi.NewRouter()

	// Seek, verify and validate JWT tokens
	r.Use(jwtauth.Verifier(TokenAuth))

	// Handle valid / invalid tokens. This is a custom authenticator
	// based on the jwtauth.Authenticator method.
	r.Use(Authenticator)

	//r.Get("/", HandlerWrapper(http.HandlerFunc(getUsers), settings))
	r.Put("/", HandlerWrapper(http.HandlerFunc(updateUser), settings))

	r.Route("/{userId:^[0-9]+$}", func(r chi.Router) {
		// HandlerWrapper is a custom wrapper that returns a Handler,
		// it allows settings to be passed to the handler
		r.Get("/", HandlerWrapper(http.HandlerFunc(getUser), settings))
		r.Delete("/", HandlerWrapper(http.HandlerFunc(deleteUser), settings))
	})

	return r
}
