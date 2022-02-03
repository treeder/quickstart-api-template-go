package main

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi/v5"
	"github.com/treeder/firetils"
	"github.com/treeder/gotils/v2"
	"github.com/treeder/temp/globals"
)

func setupRoutes(ctx context.Context, r chi.Router) {
	r.Get("/", gotils.ErrorHandler(hi))
	r.Route("/v1", func(r chi.Router) {

		r.Post("/msg", gotils.ErrorHandler(postMsg))
		r.Get("/msgs", gotils.ErrorHandler(getMsgs))

		r.Route("/chains", func(r chi.Router) {
			// r.Get("/{id}", gotils.ErrorHandler(getChain))
		})
		r.Route("/tokens", func(r chi.Router) {
			// r.Get("/{id}", gotils.ErrorHandler(getToken))
		})
	})
}

func hi(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	gotils.L(ctx).Info().Println("hi!")

	// TODO: store this in our own db as we build it up
	gotils.WriteObject(w, http.StatusOK, map[string]interface{}{"hello": "world"})

	return nil
}

type Msg struct {
	firetils.Firestored
	firetils.TimeStamped
	firetils.IDed
	Msg string `firestore:"msg" json:"msg"`
}

func postMsg(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	mi := &Msg{}
	err := gotils.ParseJSONReader(r.Body, mi)
	if err != nil {
		return gotils.C(ctx).Errorf("bad input: %w", err)
	}

	v, err := firetils.Save(ctx, globals.App.Db, "msgs", mi)
	if err != nil {
		return gotils.C(ctx).Errorf("fs error: %w", err)
	}

	// TODO: store this in our own db as we build it up
	gotils.WriteObject(w, http.StatusOK, map[string]interface{}{"message": v})

	return nil
}

func getMsgs(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	mi := &Msg{}

	vs, err := firetils.GetAllByQuery2(ctx, globals.App.Db.Collection("msgs").OrderBy("createdAt", firestore.Desc), mi)
	if err != nil {
		return gotils.C(ctx).Errorf("fs error: %w", err)
	}

	// TODO: store this in our own db as we build it up
	gotils.WriteObject(w, http.StatusOK, map[string]interface{}{"messages": vs})

	return nil
}
