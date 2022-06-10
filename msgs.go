package main

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi/v5"
	"github.com/treeder/firetils"
	"github.com/treeder/gotils/v2"
	"github.com/treeder/quickstart/globals"
)

type Msg struct {
	firetils.Firestored
	firetils.TimeStamped
	firetils.IDed
	firetils.OwnedBy
	Msg string `firestore:"msg" json:"msg"`
}

type MsgInput struct {
	Msg *Msg `json:"msg"`
}

func postMsg(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	fmt.Println("ID:", id)

	mi := &MsgInput{}
	err := gotils.ParseJSONReader(r.Body, mi)
	if err != nil {
		return gotils.C(ctx).Errorf("bad input: %w", err)
	}
	msg := mi.Msg

	if id != "" {
		current := &Msg{}
		err = firetils.GetByID(ctx, globals.App.Db, "msgs", id, current)
		if err != nil {
			return gotils.C(ctx).Errorf("firestore error: %w", err)
		}
		if current.UserID != firetils.UserID(ctx) {
			return gotils.C(ctx).Errorf("")
		}

		current.Msg = msg.Msg
		msg = current
	}

	v, err := firetils.Save(ctx, globals.App.Db, "msgs", msg)
	if err != nil {
		return gotils.C(ctx).Errorf("fs error: %w", err)
	}

	// TODO: store this in our own db as we build it up
	gotils.WriteObject(w, http.StatusOK, map[string]interface{}{"message": v})

	return nil
}

func getMsgs(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	userID := firetils.UserID(ctx)
	fmt.Println("USERID:", userID)
	mi := &Msg{}

	vs, err := firetils.GetAllByQuery2(ctx, globals.App.Db.Collection("msgs").OrderBy("createdAt", firestore.Desc), mi)
	if err != nil {
		return gotils.C(ctx).Errorf("fs error: %w", err)
	}

	// TODO: store this in our own db as we build it up
	gotils.WriteObject(w, http.StatusOK, map[string]interface{}{"messages": vs})

	return nil
}

func getMsg(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	fmt.Println("ID:", id)

	mi := &Msg{}
	err := firetils.GetByID(ctx, globals.App.Db, "msgs", id, mi)
	if err != nil {
		return gotils.C(ctx).Errorf("firestore error: %w", err)
	}

	// TODO: store this in our own db as we build it up
	gotils.WriteObject(w, http.StatusOK, map[string]interface{}{"message": mi})

	return nil
}

func deleteMsg(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	fmt.Println("ID:", id)

	err := firetils.Delete(ctx, globals.App.Db, "msgs", id)
	if err != nil {
		return gotils.C(ctx).Error(err)
	}

	// TODO: store this in our own db as we build it up
	gotils.WriteObject(w, http.StatusOK, map[string]interface{}{"all": "good"})

	return nil

}
