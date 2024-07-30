package main

import (
	"flag"
	"fmt"
	"os"
	"translator/internal/translator"
)

const (
	envTranslatorEndpoint = "TRANSLATOR_ENDPOINT"
	envTranslatorKey      = "TRANSLATOR_KEY"
	envTranslatorRegion   = "TRANSLATOR_REGION"
	envBlobAccount        = "BLOB_STORAGE_ACCOUNT_NAME"
	envBlobAccountKey     = "BLOB_STORAGE_ACCOUNT_KEY"
	envBlobContainer      = "BLOB_STORAGE_CONTAINER_NAME"
)

var verbose bool

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	endpoint := flag.String("endpoint", os.Getenv(envTranslatorEndpoint), "Azure Translator API endpoint")
	key := flag.String("key", os.Getenv(envTranslatorKey), "Azure Translator API key")
	region := flag.String("region", os.Getenv(envTranslatorRegion), "Azure region")
	in := flag.String("in", "", "Input file path")
	from := flag.String("from", "", "Source language")
	to := flag.String("to", "", "Target language")
	out := flag.String("out", "", "Destination file path")
	timeout := flag.Int("timeout", 30, "Timeout in seconds")
	blobAccount := flag.String("blobAccount", os.Getenv(envBlobAccount), "Azure Blob Storage account name")
	blobAccountKey := flag.String("blobAccountKey", os.Getenv(envBlobAccountKey), "Azure Blob Storage account key")
	blobContainer := flag.String("blobContainer", os.Getenv(envBlobContainer), "Azure Blob Storage container name")
	flag.BoolVar(&verbose, "v", false, "enable verbose logging")
	flag.Parse()

	if err := validateInputs(*endpoint, *key, *region, *in, *out, *to, *blobAccount, *blobAccountKey, *blobContainer); err != nil {
		return err
	}

	return translator.TranslateDocument(*in, *out, *from, *to, *endpoint, *key, *region, *blobAccount, *blobAccountKey, *blobContainer, *timeout, verbose)
}

func validateInputs(endpoint, key, region, in, out, to, blobAccount, blobAccountKey, blobContainer string) error {
	if endpoint == "" || key == "" || region == "" || in == "" || out == "" || to == "" || blobAccount == "" || blobAccountKey == "" || blobContainer == "" {
		return fmt.Errorf("all parameters are required")
	}
	// Add more specific validations here if needed
	return nil
}
