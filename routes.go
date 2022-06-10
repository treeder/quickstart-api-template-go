package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/treeder/firetils"
	"github.com/treeder/gotils/v2"
	"github.com/treeder/quickstart/globals"
)

func setupRoutes(ctx context.Context, r chi.Router) {
	r.Get("/", gotils.ErrorHandler(hi))
	r.Route("/v1", func(r chi.Router) {

		r.With(firetils.FireAuth).Post("/session", gotils.ErrorHandler(createSession))

		r.Route("/msgs", func(r chi.Router) {
			r.With(firetils.FireAuth).Post("/", gotils.ErrorHandler(postMsg))
			r.With(firetils.OptionalAuth).Get("/", gotils.ErrorHandler(getMsgs))
			r.With(firetils.OptionalAuth).Get("/{id}", gotils.ErrorHandler(getMsg))
			r.With(firetils.FireAuth).Post("/{id}", gotils.ErrorHandler(postMsg))
			r.With(firetils.FireAuth).Delete("/{id}", gotils.ErrorHandler(deleteMsg))
		})

	})
}

func createSession(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	expiresIn := time.Hour * 24 * 14
	// Create the session cookie. This will also verify the ID token in the process.
	// The session cookie will have the same claims as the ID token.
	// To only allow session cookie setting on recent sign-in, auth_time in ID token
	// can be checked to ensure user was recently signed in before creating a session cookie.
	idToken := r.Header.Get("Authorization")
	splitToken := strings.Split(idToken, " ")
	cookie, err := globals.App.Auth.SessionCookie(ctx, splitToken[1], expiresIn)
	if err != nil {
		return gotils.NewHTTPError("Failed to create a cookie", http.StatusInternalServerError)
	}
	gotils.WriteObject(w, http.StatusOK, map[string]interface{}{"cookie": cookie, "expires": int(expiresIn.Seconds())})
	return nil
}

func hi(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	gotils.L(ctx).Info().Println("hi!")

	// TODO: store this in our own db as we build it up
	gotils.WriteObject(w, http.StatusOK, map[string]interface{}{"hello": "world"})

	return nil
}
