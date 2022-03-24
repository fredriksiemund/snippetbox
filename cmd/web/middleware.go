package main

import (
	"context"
	"errors"
	"fmt"
	"fredriksiemund/snippetbox/pkg/models"
	"net/http"

	"github.com/justinas/nosurf"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if id exists
		exists := app.session.Exists(r, "authenticatedUserId")
		if !exists {
			// Proceed to next handler without setting context
			next.ServeHTTP(w, r)
			return
		}

		// Check if a user with the provided id exists
		id := app.session.GetInt(r, "authenticatedUserId")
		user, err := app.users.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				// User does not exist, remove from session data and proceed without setting context
				app.session.Remove(r, "authenticatedUserId")
				next.ServeHTTP(w, r)
				return
			} else {
				// Something else went wrong
				app.serverError(w, err)
				return
			}
		}

		// Verify user is active
		if !user.Active {
			app.session.Remove(r, "authenticatedUserId")
			next.ServeHTTP(w, r)
			return
		}

		// Include information in request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, contextKeyIsAuthenticated, true)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
