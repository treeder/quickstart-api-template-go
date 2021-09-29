package globals

import (
	"cloud.google.com/go/firestore"
	"github.com/dgraph-io/ristretto"
)

var (
	Cache *ristretto.Cache
	Fs    *firestore.Client
)
