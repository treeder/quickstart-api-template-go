package globals

import (
	"cloud.google.com/go/firestore"
	"github.com/dgraph-io/ristretto"
)

func init() {
	App = &MyApp{}
}

var (
	App *MyApp
)

type MyApp struct {
	Cache *ristretto.Cache
	Db    *firestore.Client
}
