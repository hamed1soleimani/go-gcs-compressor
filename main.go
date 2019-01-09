package main

import (
	"archive/zip"
	"context"
	"log"
	"os"
	"io"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func main() {
	folder := os.Args[1]
	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
            	log.Fatal(err)
	}

	rw := storageClient.Bucket(os.Getenv("ZIP_BUCKET")).Object(folder+".zip").NewWriter(ctx)

    	myZip := zip.NewWriter(rw)

	it := storageClient.Bucket(os.Getenv("BUCKET")).Objects(ctx, &storage.Query{Prefix: folder + "/"})
	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
            		log.Fatal(err)
		}

		rc, err := storageClient.Bucket(os.Getenv("BUCKET")).Object(objAttrs.Name).NewReader(ctx)
		if err != nil {
            		log.Fatal(err)
		}

		fileName := strings.TrimPrefix(objAttrs.Name, folder)
		fileName = strings.TrimPrefix(fileName, "/")
        	zipFile, err := myZip.Create(fileName)
        	if err != nil {
            		log.Fatal(err)
        	}
      
     		_, err = io.Copy(zipFile, rc)
        	if err != nil {
            		log.Fatal(err)
        	}
    	}

 	err = myZip.Close()
    	if err != nil {
            	log.Fatal(err)
	}

	err = rw.Close()
 	if err != nil {
            	log.Fatal(err)
	}
}
