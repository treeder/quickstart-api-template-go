package main

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"cloud.google.com/go/storage"
	"github.com/treeder/gotils/v2"
	"github.com/treeder/quickstart/globals"
)

func mediaPath(thingID, filename string) string {
	path := fmt.Sprintf("things/%v/%v", thingID, filename)
	return path
}

func mediaURL(bucketName, thingID, filename string) string {
	// if strings.HasPrefix(path, "/") {
	// 	// need to strip the slash since it'll get encoded below
	// 	path = path[1:]
	// }
	return fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%v/o/%v?alt=media", bucketName, url.PathEscape(mediaPath(thingID, filename)))
}
func storeImage(ctx context.Context, bucketName string, path string, file io.Reader) error {
	bucket, err := globals.App.Storage.Bucket(bucketName)
	if err != nil {
		return gotils.C(ctx).Errorf("error getting storage bucket handle: $v", err)
	}
	if err := upload(ctx, bucket, path, file); err != nil {
		return err
	}
	return nil

}
func upload(ctx context.Context, bucket *storage.BucketHandle, path string, r io.Reader) error {
	wc := bucket.Object(path).NewWriter(ctx)
	defer wc.Close()
	_, err := io.Copy(wc, r)
	return err
}
