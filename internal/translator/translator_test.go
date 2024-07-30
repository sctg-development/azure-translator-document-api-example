package translator

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

const (
	BLOB_NAME = "TestREADME.md"
)

func TestUploadFileToBlobStorage(t *testing.T) {
	accountName := os.Getenv("ACCOUNT_NAME")
	accountKey := os.Getenv("ACCOUNT_KEY")
	containerName := os.Getenv("CONTAINER_NAME")
	filePath := "../../README.md"
	blobName := BLOB_NAME

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}
	err = uploadFileToBlobStorage(accountName, accountKey, containerName, absFilePath, blobName)

	if err != nil {
		t.Errorf("UploadFileToBlobStorage failed: %v", err)
	}

	// Add additional assertions or verifications here
}
func TestGetBlobURLWithSASToken(t *testing.T) {
	accountName := os.Getenv("ACCOUNT_NAME")
	accountKey := os.Getenv("ACCOUNT_KEY")
	containerName := os.Getenv("CONTAINER_NAME")
	blobName := BLOB_NAME

	urlWithSASToken, _, err := getBlobURLWithSASToken(accountName, accountKey, containerName, blobName)
	fmt.Println(urlWithSASToken)
	if err != nil {
		t.Errorf("getBlobURLWithSASToken failed: %v", err)
	}

	// test if the URL return 200 status code
	_, err = http.Get(urlWithSASToken)
	if err != nil {
		t.Errorf("getBlobURLWithSASToken failed: %v", err)
	}

	// Add additional assertions or verifications here
}

func TestDeleteFileFromBlobStorage(t *testing.T) {
	accountName := os.Getenv("ACCOUNT_NAME")
	accountKey := os.Getenv("ACCOUNT_KEY")
	containerName := os.Getenv("CONTAINER_NAME")
	blobName := BLOB_NAME

	err := deleteFileFromBlobStorage(accountName, accountKey, containerName, blobName)
	if err != nil {
		t.Errorf("DeleteFileFromBlobStorage failed: %v", err)
	}

	// Add additional assertions or verifications here
}
func TestGenerateJSONDocument(t *testing.T) {
	sourceSASUrl := "https://example.com/sourceSASUrl"
	targetSASUrl := "https://example.com/targetSASUrl"
	sourceLanguage := "en"
	targetLanguage := "fr"

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

	jsonDoc, err := generateJSONDocument(sourceSASUrl, targetSASUrl, sourceLanguage, targetLanguage)
	if err != nil {
		t.Errorf("generateJSONDocument failed: %v", err)
	}

	if jsonDoc != expectedJSON {
		t.Errorf("generateJSONDocument returned incorrect JSON.\nExpected:\n%s\n\nActual:\n%s", expectedJSON, jsonDoc)
	}

	// Add additional assertions or verifications here
}
func TestTranslateDocument(t *testing.T) {
	fileToTranslate := "../../README.md"
	fileTranslated := "../../README_fr.md"
	sourceLanguage := "en"
	targetLanguage := "fr"
	endpoint := os.Getenv("TRANSLATOR_ENDPOINT")
	key := os.Getenv("TRANSLATOR_KEY")
	region := os.Getenv("TRANSLATOR_REGION")
	blobAccountName := os.Getenv("ACCOUNT_NAME")
	blobAccountKey := os.Getenv("ACCOUNT_KEY")
	blobContainerName := os.Getenv("CONTAINER_NAME")

	err := TranslateDocument(fileToTranslate, fileTranslated, sourceLanguage, targetLanguage, endpoint, key, region, blobAccountName, blobAccountKey, blobContainerName, 30, true)
	if err != nil {
		t.Errorf("TranslateDocument failed: %v", err)
	}

	// Add additional assertions or verifications here
}
