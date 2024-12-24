package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	mux := mux.NewRouter()

	mux.NotFoundHandler = http.HandlerFunc(app.notFound)
	mux.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowed)

	mux.Use(app.logAccess)
	mux.Use(app.recoverPanic)
	mux.Use(app.authenticate)
	mux = mux.PathPrefix("/api/" + app.config.apiVersion).Subrouter()
	// Public routes
	mux.HandleFunc("/status", app.status).Methods("GET")

	authRoutes := mux.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/register", app.createUser).Methods("POST")
	authRoutes.HandleFunc("/login", app.createAuthenticationToken).Methods("POST")

	//Protected routes
	authenticatedRoutes := mux.NewRoute().Subrouter()
	authenticatedRoutes.Use(app.requireAuthenticatedUser)
	authenticatedRoutes.HandleFunc("/protected", app.protected).Methods("GET")

	return mux
}
