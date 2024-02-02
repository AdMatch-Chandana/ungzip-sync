package main

import (
	"context"
	"os"

	ungzip "github.com/admatch/adform-sync-ungzip"
	"github.com/peterbourgon/ff/v4"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()

	fs := ff.NewFlagSet("adform-sync")
	sourceBucket := fs.String('s', "source", "", "bucket name to retrieve the zip file")
	objectPath := fs.String('o', "object", "", "object path")
	_ = fs.String('d', "destination", "", "bucket name to unzip files")
	_ = fs.String('p', "project-id", "", "project id for logs")
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVars())
	if err != nil {
		sugar.Fatalf("Failed to parse flags: %v", err)
	}

	ctx := context.TODO()
	e := ungzip.GCSEvent{
		Bucket: *sourceBucket, //"bkt-prj-eng-d-adform-sync-svc-ewh5-adform-sync-out",
		Name:   *objectPath,   //"",
	}

	sugar.Infof("Source: %s Object: %s", *sourceBucket, *objectPath)

	err = ungzip.UncompressFile(ctx, e)
	if err != nil {
		sugar.Fatalf("Fail to unzip and copy: %v", err)
	}

	sugar.Info("Success!")
}
