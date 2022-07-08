package main

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/treeder/firetils"
	"github.com/treeder/gotils/v2"
	"github.com/treeder/quickstart/globals"
)

type Msg struct {
	firetils.Firestored
	firetils.TimeStamped
	firetils.IDed
	firetils.OwnedBy
	Msg   string `firestore:"msg" json:"msg"`
	Image string `firestore:"image" json:"image"`
}

type MsgInput struct {
	Msg *Msg `json:"msg"`
}

func postMsg(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var err error

	id := chi.URLParam(r, "id")
	fmt.Println("ID:", id)
	randid := ""
	if id == "" {
		randid, _ = gonanoid.New()
	}

	input := &MsgInput{}

	var file multipart.File
	filename := ""
	gotils.L(ctx).Info().Printf("ctype: %v", r.Header.Get("Content-Type"))
	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		// this is the path it should always take now
		if err := r.ParseMultipartForm(40 << 20); err != nil { // think that's 40MB?
			return err
		}
		jsonPart := r.FormValue("json")
		fmt.Println("jsonPart:", jsonPart)
		err = gotils.ParseJSONBytes([]byte(jsonPart), input)
		if err != nil {
			return gotils.C(ctx).Error(err)
		}

		// and now media part
		f2, header, err := r.FormFile("image")
		if err != nil {
			return err
		}
		defer f2.Close()
		file = f2

		ext := filepath.Ext(header.Filename)

		filename = fmt.Sprintf("media-%v%v", randid, ext)

	} else {
		err := gotils.ParseJSONReader(r.Body, input)
		if err != nil {
			return gotils.C(ctx).Errorf("bad input: %w", err)
		}
	}
	msg := input.Msg
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
	} else {
		msg.ID = randid
	}

	// store file if one was uploaded
	if filename != "" {
		imagePath := mediaPath(msg.ID, filename)
		gotils.L(ctx).Info().Printf("storing image to %v", imagePath)
		err = storeImage(ctx, globals.App.StorageBucket, imagePath, file)
		if err != nil {
			return gotils.C(ctx).Errorf("error storing image: %w", err)
		}
		msg.Image = mediaURL(globals.App.StorageBucket, msg.ID, filename)
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
