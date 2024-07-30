/*
Copyright (c) Ronan LE MEILLAT 2024
Licensed under the AGPLv3 License
https://www.gnu.org/licenses/agpl-3.0.html

Package translator provides functions for uploading, retrieving, and deleting files from Azure Blob Storage,
generating JSON documents for translation, and translating documents using the Azure Translator service.

This file contains test functions for the various operations provided by the translator package.
*/

package translator

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

// BLOB_NAME is the name of the blob file used for testing.
const (
	BLOB_NAME = "TestREADME.md"
)

// TestUploadFileToBlobStorage is a test function that tests the uploadFileToBlobStorage function.
func TestUploadFileToBlobStorage(t *testing.T) {
	// Test configuration
	config := TranslatorConfig{
		TranslatorEndpoint: os.Getenv("TRANSLATOR_ENDPOINT"),
		TranslatorKey:      os.Getenv("TRANSLATOR_KEY"),
		TranslatorRegion:   os.Getenv("TRANSLATOR_REGION"),
		BlobAccountName:    os.Getenv("BLOB_STORAGE_ACCOUNT_NAME"),
		BlobAccountKey:     os.Getenv("BLOB_STORAGE_ACCOUNT_KEY"),
		BlobContainerName:  os.Getenv("BLOB_STORAGE_CONTAINER_NAME"),
		Timeout:            30,
		Verbose:            true,
	}

	// File paths
	filePath := "../../README.md"
	blobName := BLOB_NAME

	// Get absolute file path
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Upload file to blob storage
	err = uploadFileToBlobStorage(config, absFilePath, blobName)
	if err != nil {
		t.Errorf("UploadFileToBlobStorage failed: %v", err)
	}

	// Add additional assertions or verifications here
}

// TestGetBlobURLWithSASToken is a test function that tests the getBlobURLWithSASToken function.
func TestGetBlobURLWithSASToken(t *testing.T) {
	// Test configuration
	config := TranslatorConfig{
		TranslatorEndpoint: os.Getenv("TRANSLATOR_ENDPOINT"),
		TranslatorKey:      os.Getenv("TRANSLATOR_KEY"),
		TranslatorRegion:   os.Getenv("TRANSLATOR_REGION"),
		BlobAccountName:    os.Getenv("BLOB_STORAGE_ACCOUNT_NAME"),
		BlobAccountKey:     os.Getenv("BLOB_STORAGE_ACCOUNT_KEY"),
		BlobContainerName:  os.Getenv("BLOB_STORAGE_CONTAINER_NAME"),
		Timeout:            30,
		Verbose:            true,
	}

	// Blob name
	blobName := BLOB_NAME

	// Get blob URL with SAS token
	urlWithSASToken, _, err := getBlobURLWithSASToken(config, blobName)
	fmt.Println(urlWithSASToken)
	if err != nil {
		t.Errorf("getBlobURLWithSASToken failed: %v", err)
	}

	// Test if the URL returns a 200 status code
	_, err = http.Get(urlWithSASToken)
	if err != nil {
		t.Errorf("getBlobURLWithSASToken failed: %v", err)
	}

	// Add additional assertions or verifications here
}

// TestDeleteFileFromBlobStorage is a test function that tests the deleteFileFromBlobStorage function.
func TestDeleteFileFromBlobStorage(t *testing.T) {
	// Test configuration
	config := TranslatorConfig{
		TranslatorEndpoint: os.Getenv("TRANSLATOR_ENDPOINT"),
		TranslatorKey:      os.Getenv("TRANSLATOR_KEY"),
		TranslatorRegion:   os.Getenv("TRANSLATOR_REGION"),
		BlobAccountName:    os.Getenv("BLOB_STORAGE_ACCOUNT_NAME"),
		BlobAccountKey:     os.Getenv("BLOB_STORAGE_ACCOUNT_KEY"),
		BlobContainerName:  os.Getenv("BLOB_STORAGE_CONTAINER_NAME"),
		Timeout:            30,
		Verbose:            true,
	}

	// Blob name
	blobName := BLOB_NAME

	// Delete file from blob storage
	err := deleteFileFromBlobStorage(config, blobName)
	if err != nil {
		t.Errorf("DeleteFileFromBlobStorage failed: %v", err)
	}

	// Add additional assertions or verifications here
}

// TestGenerateJSONDocument is a test function that tests the generateJSONDocument function.
func TestGenerateJSONDocument(t *testing.T) {
	// Test data
	sourceSASUrl := "https://example.com/sourceSASUrl"
	targetSASUrl := "https://example.com/targetSASUrl"
	sourceLanguage := "en"
	targetLanguage := "fr"

	// Expected JSON document
	expectedJSON := `{
  "inputs": [
    {
      "storageType": "File",
      "source": {
        "sourceUrl": "https://example.com/sourceSASUrl",
        "language": "en"
      },
      "targets": [
        {
          "targetUrl": "https://example.com/targetSASUrl",
          "language": "fr"
        }
      ]
    }
  ]
}`

	// Generate JSON document
	jsonDoc, err := generateJSONDocument(sourceSASUrl, targetSASUrl, sourceLanguage, targetLanguage)
	if err != nil {
		t.Errorf("generateJSONDocument failed: %v", err)
	}

	// Compare generated JSON document with expected JSON document
	if jsonDoc != expectedJSON {
		t.Errorf("generateJSONDocument returned incorrect JSON.\nExpected:\n%s\n\nActual:\n%s", expectedJSON, jsonDoc)
	}

	// Add additional assertions or verifications here
}

// TestTranslateDocument is a test function that tests the TranslateDocument function.
func TestTranslateDocument(t *testing.T) {
	// Test data
	fileToTranslate := "../../README.md"
	fileTranslated := "../../README_fr.md"
	sourceLanguage := "en"
	targetLanguage := "fr"

	// Test configuration
	config := TranslatorConfig{
		TranslatorEndpoint: os.Getenv("TRANSLATOR_ENDPOINT"),
		TranslatorKey:      os.Getenv("TRANSLATOR_KEY"),
		TranslatorRegion:   os.Getenv("TRANSLATOR_REGION"),
		BlobAccountName:    os.Getenv("BLOB_STORAGE_ACCOUNT_NAME"),
		BlobAccountKey:     os.Getenv("BLOB_STORAGE_ACCOUNT_KEY"),
		BlobContainerName:  os.Getenv("BLOB_STORAGE_CONTAINER_NAME"),
		Timeout:            30,
		Verbose:            true,
	}

	// Translate document
	err := TranslateDocument(fileToTranslate, fileTranslated, sourceLanguage, targetLanguage, config)
	if err != nil {
		t.Errorf("TranslateDocument failed: %v", err)
	}

	// Add additional assertions or verifications here
}
