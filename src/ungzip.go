package adformSyncUnzip

import (
	"compress/gzip"
	"context"
	"io"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/dhduvall/gcloudzap"
	"github.com/peterbourgon/ff/v4"
	"go.uber.org/zap"
)

var (
	fs                = ff.NewFlagSet("adform-sync")
	destinationBucket = fs.String('d', "destination", "", "bucket name to ungzip files")
	projectID         = fs.String('p', "project-id", "", "project id for logs")
	_                 = fs.String('s', "source", "", "bucket name to retrieve the gzip file")
	_                 = fs.String('o', "object", "", "object path")
	storageClient     *storage.Client
	logger            *zap.Logger
)

func initialise() {
	var err error
	err = ff.Parse(fs, os.Args[1:], ff.WithEnvVars())
	if err != nil {
		log.Fatalf("Failed to parse flags: %v", err)
	}
	logger, err = gcloudzap.NewProduction(*projectID, "adform-sync-svc-ungzip")
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	storageClient, err = storage.NewClient(context.Background())
	if err != nil {
		logger.Fatal("Failed to create storage client", zap.Error(err))
	}
}

// GCSEvent is the payload of a GCS event.
type GCSEvent struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

func UncompressFile(ctx context.Context, e GCSEvent) error {
	initialise()

	fnLogger := logger.With(
		zap.String("object", e.Name),
		zap.String("source", e.Bucket),
	)
	fnLogger.Info("Decompressing file")

	// Get the source and destination buckets
	sourceBucketObj := storageClient.Bucket(e.Bucket)
	destinationBucketObj := storageClient.Bucket(*destinationBucket)

	// Read file from the bucket
	sourceReader, err := sourceBucketObj.Object(e.Name).NewReader(ctx)
	if err != nil {
		fnLogger.Fatal("Error reading file", zap.Error(err))
	}
	defer sourceReader.Close()

	// Create a gzip reader
	gzipReader, err := gzip.NewReader(sourceReader)
	if err != nil {
		fnLogger.Fatal("Failed to create gzip reader", zap.Error(err))
	}
	defer gzipReader.Close()

	// Create a writer to write the unzipped file to the destination bucket
	destinationFile := strings.TrimSuffix(e.Name, ".gz")
	destinationWriter := destinationBucketObj.Object(destinationFile).NewWriter(ctx)
	defer destinationWriter.Close()

	fnLogger = fnLogger.With(
		zap.String("destination", *destinationBucket),
		zap.String("file", destinationFile),
	)

	// Copy the unzipped data from the gzip reader to the destination writer
	_, err = io.Copy(destinationWriter, gzipReader)
	if err != nil {
		fnLogger.Fatal("Failed to copy unzipped data to destination bucket", zap.Error(err))
	}

	fnLogger.Info("File unzipped and uploaded to destination bucket")

	return nil
}
