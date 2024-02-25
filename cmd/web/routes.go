package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir(cfg.staticDir))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	publicRoutes := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/", publicRoutes.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", publicRoutes.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", publicRoutes.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", publicRoutes.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", publicRoutes.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", publicRoutes.ThenFunc(app.userLoginPost))

	protectedRoutes := publicRoutes.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/snippet/create", protectedRoutes.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protectedRoutes.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protectedRoutes.ThenFunc(app.userLogoutPost))

	middlewares := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return middlewares.Then(router)
}
