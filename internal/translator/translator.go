/* Copyright (c) Ronan LE MEILLAT 2024
 * Licensed under the AGPLv3 License
 * https://www.gnu.org/licenses/agpl-3.0.html
 */
package translator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ATranslateDocument struct {
	Inputs []Input `json:"inputs"`
}

type Input struct {
	StorageType string   `json:"storageType"`
	Source      Source   `json:"source"`
	Targets     []Target `json:"targets"`
}

type Source struct {
	SourceURL string `json:"sourceUrl"`
	Language  string `json:"language,omitempty"`
}

type Target struct {
	TargetURL string `json:"targetUrl"`
	Language  string `json:"language"`
}

const (
	APIVersion = "2024-05-01"
)

// TranslatorConfig represents the configuration for the Translator service.
type TranslatorConfig struct {
	BlobAccountName    string `json:"blobAccountName"`
	BlobAccountKey     string `json:"blobAccountKey"`
	BlobContainerName  string `json:"blobContainerName"`
	TranslatorEndpoint string `json:"translatorEndpoint"`
	TranslatorKey      string `json:"translatorKey"`
	TranslatorRegion   string `json:"translatorRegion"`
	Timeout            int    `json:"timeout"`
	Verbose            bool   `json:"verbose"`
	Logger             *logrus.Logger
}

// uploadFileToBlobStorage uploads a local file to Azure Blob Storage.
// It takes the following parameters:
// - config: The TranslatorConfig object.
// - filePath: The path to the local file to be uploaded.
// - blobName: The name of the blob in Azure Blob Storage.
// It returns an error if any.
func uploadFileToBlobStorage(config TranslatorConfig, filePath, blobName string) error {
	config.Logger.Debugf("Uploading file %s to blob storage as %s", filePath, blobName)
	accountName := config.BlobAccountName
	accountKey := config.BlobAccountKey
	containerName := config.BlobContainerName
	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return err
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// Create a URL to the blob storage container.
	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	// Create a container URL object.
	containerURL := azblob.NewContainerURL(*URL, p)

	// Open the local file.
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get the file size.
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()

	// Create a blob URL object.
	blobURL := containerURL.NewBlockBlobURL(blobName)

	// Create a context with cancellation.
	ctx := context.Background()

	// Create a buffer to read the file content.
	buffer := make([]byte, fileSize)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	// Upload the file to Azure Blob Storage.
	_, err = azblob.UploadBufferToBlockBlob(ctx, buffer, blobURL, azblob.UploadToBlockBlobOptions{})
	if err != nil {
		return err
	}

	config.Logger.Debugf("Successfully uploaded file to blob storage: %s", blobName)
	return nil
}

// deleteFileFromBlobStorage deletes a file from Azure Blob Storage.
// It takes the following parameters:
// - config: The TranslatorConfig object.
// - blobName: The name of the blob in Azure Blob Storage.
// It returns an error if any.
func deleteFileFromBlobStorage(config TranslatorConfig, blobName string) error {
	accountName := config.BlobAccountName
	accountKey := config.BlobAccountKey
	containerName := config.BlobContainerName
	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return err
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// Create a URL to the blob storage container.
	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	// Create a container URL object.
	containerURL := azblob.NewContainerURL(*URL, p)

	// Create a blob URL object.
	blobURL := containerURL.NewBlockBlobURL(blobName)

	// Create a context with cancellation.
	ctx := context.Background()

	// Delete the file from Azure Blob Storage.
	_, err = blobURL.Delete(ctx, azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{})
	if err != nil {
		return err
	}

	return nil
}

// dowloadFileFromBlobStorage downloads a file from Azure Blob Storage.
// It takes the following parameters:
// - config: The TranslatorConfig object.
// - destFilePath: The path to save the downloaded file.
// - blobName: The name of the blob in Azure Blob Storage.
// It returns an error if any.
func downloadFileFromBlobStorage(config TranslatorConfig, destFilePath, blobName string) error {
	accountName := config.BlobAccountName
	accountKey := config.BlobAccountKey
	containerName := config.BlobContainerName
	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return err
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// Create a URL to the blob storage container.
	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	// Create a container URL object.
	containerURL := azblob.NewContainerURL(*URL, p)

	// Create a blob URL object.
	blobURL := containerURL.NewBlobURL(blobName)

	// Create a context with cancellation.
	ctx := context.Background()

	err = os.MkdirAll(filepath.Dir(destFilePath), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Download the file from Azure Blob Storage.
	err = azblob.DownloadBlobToFile(ctx, blobURL, 0, 0, file, azblob.DownloadFromBlobOptions{})
	if err != nil {
		return err
	}

	return nil
}

// getBlobURLWithSASToken generates a Shared Access Signature (SAS) token for a blob in Azure Blob Storage.
// It takes the following parameters:
// - config: The TranslatorConfig object.
// - blobName: The name of the blob in Azure Blob Storage (use "" for container SAS).
// It returns the generated URL with the SAS query parameter as a string, the SAS token and an error if any.
func getBlobURLWithSASToken(config TranslatorConfig, blobName string) (string, string, error) {
	accountName := config.BlobAccountName
	accountKey := config.BlobAccountKey
	containerName := config.BlobContainerName
	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return "", "", err
	}
	path := ""
	permissions := ""
	if blobName != "" {
		path = fmt.Sprintf("/%s", blobName)
		permissions = azblob.BlobSASPermissions{Add: true, Read: true, Write: true, List: true}.String()

	} else {
		permissions = azblob.ContainerSASPermissions{Read: true, Write: true}.String()
		path = ""
	}
	// Get a SAS token for the blob.
	sasQueryParams, err := azblob.BlobSASSignatureValues{
		Protocol:      azblob.SASProtocolHTTPS,              // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: containerName,
		BlobName:      blobName,

		// To produce a container SAS (as opposed to a blob SAS), assign to Permissions using
		// ContainerSASPermissions and make sure the BlobName field is "" (the default).
		Permissions: permissions,
	}.NewSASQueryParameters(credential)
	if err != nil {
		return "", "", err
	}
	qp := sasQueryParams.Encode()
	urlToSendToSomeone := fmt.Sprintf("https://%s.blob.core.windows.net/%s%s?%s",
		accountName, containerName, path, qp)
	return urlToSendToSomeone, qp, nil
}

// generateJSONDocument generates a JSON document for translation.
// It takes the source SAS URL, target SAS URL, source language, and target language as input parameters.
// It returns the generated JSON document as a string and an error if any.
func generateJSONDocument(sourceSASUrl string, targetSASUrl string, sourceLanguage string, targetLanguage string) (string, error) {
	doc := ATranslateDocument{
		Inputs: []Input{
			{
				StorageType: "File",
				Source: Source{
					SourceURL: sourceSASUrl,
					Language:  sourceLanguage,
				},
				Targets: []Target{
					{
						TargetURL: targetSASUrl,
						Language:  targetLanguage,
					},
				},
			},
		},
	}

	jsonBytes, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}

	return string(jsonBytes), nil
}

// generateUUIDv4WithoutHyphens generates a UUID version 4 without hyphens.
// It uses the github.com/google/uuid package to generate a UUID with hyphens,
// then removes the hyphens from the generated UUID and returns the result.
func generateUUIDv4WithoutHyphens() string {
	uuidWithHyphens := uuid.New()
	uuidWithoutHyphens := uuidWithHyphens.String()
	uuidWithoutHyphens = strings.ReplaceAll(uuidWithoutHyphens, "-", "")
	return uuidWithoutHyphens
}

// TranslateDocument translates a document from one language to another using the Azure Translator service.
// It takes the following parameters:
// - fileToTranslate: The path to the local file to be translated.
// - destinationFile: The path to save the translated file.
// - sourceLanguage: The language of the source document if the provided string has zero length, the service will attempt to auto-detect the language.
// - targetLanguage: The language to translate the document to.
// - config: The TranslatorConfig object.
// It returns an error if any.
// The process of translation involves the following steps:
// 1. Generate a UUID for the translation job.
// 2. Upload the file to Azure Blob Storage.
// 3. Generate a JSON document for translation.
// 4. Translate the document.
// 5. Delete the file from Azure Blob Storage.
// 6. wait for the translated document to be ready in the target container.
func TranslateDocument(fileToTranslate, destinationFile, sourceLanguage, targetLanguage string, config TranslatorConfig) error {
	// populate endpoint, key, region, blobAccountName, blobContainerName, timeout and verbose with the values from the config object
	endpoint := config.TranslatorEndpoint
	key := config.TranslatorKey
	region := config.TranslatorRegion
	blobAccountName := config.BlobAccountName
	blobContainerName := config.BlobContainerName

	// Generate a UUID for the translation job.
	jobID := generateUUIDv4WithoutHyphens()
	config.Logger.Debugf("Starting translation job %s", jobID)
	filename := filepath.Base(fileToTranslate)
	srcJobID := fmt.Sprintf("%s-%s", jobID, filename)
	dstJobID := fmt.Sprintf("%s-translated-%s", jobID, filename)
	// Upload the file to Azure Blob Storage.
	err := uploadFileToBlobStorage(config, fileToTranslate, srcJobID)
	if err != nil {
		return fmt.Errorf("error uploading file to Azure Blob Storage: %v", err)
	}

	// Generate a JSON document for translation.
	containerSASurl, containerSASToken, err := getBlobURLWithSASToken(config, "")
	if err != nil {
		return fmt.Errorf("error generating container SAS token: %v", err)
	}
	config.Logger.Debugf("containerSASurl: %s", containerSASurl)

	sourceSASUrl := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s", blobAccountName, blobContainerName, srcJobID, containerSASToken)

	config.Logger.Debugf("sourceSASUrl: %s", sourceSASUrl)

	targetSASUrl := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s", blobAccountName, blobContainerName, dstJobID, containerSASToken)

	config.Logger.Debugf("targetSASUrl: %s", targetSASUrl)

	jsonDocument, err := generateJSONDocument(sourceSASUrl, targetSASUrl, sourceLanguage, targetLanguage)
	if err != nil {
		return fmt.Errorf("error generating JSON document: %v", err)
	}

	// Translate the document.
	basePath := fmt.Sprintf("%s/translator/document/batches", endpoint)
	uri := fmt.Sprintf("%s?api-version=%s", basePath, APIVersion)
	method := "POST"
	req, err := http.NewRequest(method, uri, bytes.NewBuffer([]byte(jsonDocument)))
	req.Header.Add("Ocp-Apim-Subscription-Key", key)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Ocp-Apim-Subscription-Region", region)
	client := &http.Client{}
	if err != nil {

		return fmt.Errorf("error creating HTTP request: %v", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %v", err)

	}
	defer res.Body.Close()

	config.Logger.Debugf("response status: %s", res.Status)
	config.Logger.Debugf("response headers: %v", res.Header)

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		// Wait for the translated document to be ready in the target container.
		_ = waitForTranslatedFile(config, destinationFile, dstJobID)
	}

	// Delete the file from Azure Blob Storage.
	err = deleteFileFromBlobStorage(config, srcJobID)
	if err != nil {
		return fmt.Errorf("error deleting source document: %v", err)
	}
	err = deleteFileFromBlobStorage(config, dstJobID)
	if err != nil {
		return fmt.Errorf("error deleting translated document: %v", err)
	}
	config.Logger.Debugf("Translation job %s completed successfully", jobID)
	return nil
}

// waitForTranslatedFile waits for the translated file to be available in the destination storage.
// It continuously checks for the availability of the file by downloading it from the blob storage.
// The function takes the TranslatorConfig, destinationFile path, and dstJobID as parameters.
// It returns an error if the file is not available within the specified timeout period.
func waitForTranslatedFile(config TranslatorConfig, destinationFile string, dstJobID string) error {
	maxTry := config.Timeout
	for {
		err := downloadFileFromBlobStorage(config, destinationFile, dstJobID)
		if err == nil {
			return nil
		}
		if maxTry <= 0 {
			return fmt.Errorf("timeout waiting for the translated document")
		}
		config.Logger.Infoln("File not yet ready, wait for 1s…")
		time.Sleep(1 * time.Second)
		maxTry--
	}
}
