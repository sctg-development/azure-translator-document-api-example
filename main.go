/* Copyright (c) Ronan LE MEILLAT 2024
 * Licensed under the AGPLv3 License
 * https://www.gnu.org/licenses/agpl-3.0.html
 */
package main

import (
	"flag"
	"log"
	"os"
	"translator/internal/translator"
)

var verbose bool

func main() {
	endpoint := flag.String("endpoint", os.Getenv("TRANSLATOR_ENDPOINT"), "Azure Translator API endpoint (default to variable TRANSLATOR_ENDPOINT)")
	key := flag.String("key", os.Getenv("TRANSLATOR_KEY"), "Azure Translator API key default to variable TRANSLATOR_KEY")
	region := flag.String("region", os.Getenv("TRANSLATOR_REGION"), "Azure region default to variable TRANSLATOR_REGION")
	in := flag.String("in", "", "Input file path")
	from := flag.String("from", "", "Source language (if not defined, the service will auto-detect the language)")
	to := flag.String("to", "", "Target language")
	out := flag.String("out", "", "Destination file path")
	timeout := flag.Int("timeout", 30, "Timeout in seconds")
	blobAccount := flag.String("blobAccount", os.Getenv("BLOB_STORAGE_ACCOUNT_NAME"), "Azure Blob Storage account name (default to variable BLOB_STORAGE_ACCOUNT_NAME)")
	blobAccountKey := flag.String("blobAccountKey", os.Getenv("BLOB_STORAGE_ACCOUNT_KEY"), "Azure Blob Storage account key (default to variable BLOB_STORAGE_ACCOUNT_KEY)")
	blobContainer := flag.String("blobContainer", os.Getenv("BLOB_STORAGE_CONTAINER_NAME"), "Azure Blob Storage container name (default to variable BLOB_STORAGE_CONTAINER_NAME)")
	flag.BoolVar(&verbose, "v", false, "enable verbose logging")
	flag.Parse()

	if *endpoint == "" || *key == "" || *region == "" || *in == "" || *out == "" || *to == "" || *blobAccount == "" || *blobAccountKey == "" || *blobContainer == "" {
		flag.Usage()
		log.Fatal("All parameters are required")
	}

	translator.TranslateDocument(*in, *out, *from, *to, *endpoint, *key, *region, *blobAccount, *blobAccountKey, *blobContainer, *timeout, verbose)
}
