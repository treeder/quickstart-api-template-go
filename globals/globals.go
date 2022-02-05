package globals

import (
	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/dgraph-io/ristretto"
)

func init() {
	App = &MyApp{}
}

var (
	App *MyApp
)

type MyApp struct {
	FireApp *firebase.App
	Auth    *auth.Client
	Cache   *ristretto.Cache
	Db      *firestore.Client
}
