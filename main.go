package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
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
	configFile := flag.String("config", "", "Configuration file path")
	flag.BoolVar(&verbose, "v", false, "enable verbose logging")
	flag.Parse()

	if *configFile != "" {
		err := loadConfigFromFile(*configFile, endpoint, key, region, blobAccount, blobAccountKey, blobContainer)
		if err != nil {
			return err
		}
	}

	if err := validateInputs(*endpoint, *key, *region, *in, *out, *to, *blobAccount, *blobAccountKey, *blobContainer); err != nil {
		return err
	}

	config := translator.TranslatorConfig{
		TranslatorEndpoint: *endpoint,
		TranslatorKey:      *key,
		TranslatorRegion:   *region,
		BlobAccountName:    *blobAccount,
		BlobAccountKey:     *blobAccountKey,
		BlobContainerName:  *blobContainer,
		Timeout:            *timeout,
		Verbose:            true,
	}
	return translator.TranslateDocument(*in, *out, *from, *to, config)
}

func loadConfigFromFile(configFile string, endpoint, key, region, blobAccount, blobAccountKey, blobContainer *string) error {
	file, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	var config translator.TranslatorConfig
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return err
	}

	*endpoint = config.TranslatorEndpoint
	*key = config.TranslatorKey
	*region = config.TranslatorRegion
	*blobAccount = config.BlobAccountName
	*blobAccountKey = config.BlobAccountKey
	*blobContainer = config.BlobContainerName

	return nil
}

func validateInputs(endpoint, key, region, in, out, to, blobAccount, blobAccountKey, blobContainer string) error {
	missingArgs := []string{}
	if endpoint == "" {
		missingArgs = append(missingArgs, "endpoint")
	}
	if key == "" {
		missingArgs = append(missingArgs, "key")
	}
	if region == "" {
		missingArgs = append(missingArgs, "region")
	}
	if in == "" {
		missingArgs = append(missingArgs, "in")
	}
	if out == "" {
		missingArgs = append(missingArgs, "out")
	}
	if to == "" {
		missingArgs = append(missingArgs, "to")
	}
	if blobAccount == "" {
		missingArgs = append(missingArgs, "blobAccount")
	}
	if blobAccountKey == "" {
		missingArgs = append(missingArgs, "blobAccountKey")
	}
	if blobContainer == "" {
		missingArgs = append(missingArgs, "blobContainer")
	}

	if len(missingArgs) > 0 {
		return fmt.Errorf("missing required arguments: %s", strings.Join(missingArgs, ", "))
	}

	return nil
}
