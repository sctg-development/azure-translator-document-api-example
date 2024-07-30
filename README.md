# Azure Translator Document API Example

This project demonstrates the usage of Azure Translator Document API to translate supported documents from one language to another using the command line. It leverages Azure Blob Storage for temporary file storage during the translation process.

## Features

- Translate documents using Azure Translator Document API
- Auto-detect source language if not specified
- Upload and download files from Azure Blob Storage
- Generate and use SAS tokens for secure blob access
- Verbose logging option for debugging

## Prerequisites

- Go (version 1.21 or higher)
- An Azure account with access to:
  - Azure Translator API (Document Translation feature)
  - Azure Blob Storage

## Installation

Clone this repository and navigate to the project directory:

```sh
git clone <REPOSITORY_URL>
cd <DIRECTORY_NAME>
```

## Configuration

Set up the following environment variables or provide them as command-line arguments:

```bash
export BLOB_STORAGE_ACCOUNT_NAME="your_blob_storage_account_name"
export BLOB_STORAGE_ACCOUNT_KEY="your_blob_storage_account_key"
export BLOB_STORAGE_CONTAINER_NAME="your_container_name"
export TRANSLATOR_ENDPOINT="https://your_translator.cognitiveservices.azure.com/"
export TRANSLATOR_KEY="your_translator_api_key"
export TRANSLATOR_REGION="your_azure_region"
```

## Usage

### Build the program

```sh
go build -o translator
```

### Run the program

```sh
./translator -in ./input_file.docx -out ./output_file.docx -to fr -v
```

### Command-line Arguments

- `-endpoint`: Azure Translator API endpoint (default: TRANSLATOR_ENDPOINT env var)
- `-key`: Azure Translator API key (default: TRANSLATOR_KEY env var)
- `-region`: Azure region (default: TRANSLATOR_REGION env var)
- `-in`: Input file path (required)
- `-out`: Output file path (required)
- `-from`: Source language (optional, auto-detected if not provided)
- `-to`: Target language (required)
- `-blobAccount`: Azure Blob Storage account name (default: BLOB_STORAGE_ACCOUNT_NAME env var)
- `-blobAccountKey`: Azure Blob Storage account key (default: BLOB_STORAGE_ACCOUNT_KEY env var)
- `-blobContainer`: Azure Blob Storage container name (default: BLOB_STORAGE_CONTAINER_NAME env var)
- `-timeout`: Timeout in seconds (default: 30)
- `-v`: Enable verbose logging

## How It Works

1. The program generates a unique job ID for the translation task.
2. It uploads the input file to Azure Blob Storage.
3. A JSON document is generated with translation parameters and SAS URLs.
4. The document is submitted for translation using the Azure Translator Document API.
5. The program waits for the translation to complete, polling the output blob.
6. Once ready, the translated document is downloaded to the specified output path.
7. Temporary blobs are deleted from Azure Blob Storage.

## Testing

Run the tests with:

```bash
export BLOB_STORAGE_ACCOUNT_NAME="my_blob_store"
export BLOB_STORAGE_ACCOUNT_KEY="MWIxYmYwMjJhMGU1MTdhMWRkZDE1YjM1OGJiNmIzOTIyYjc5MWRhNzViZTBmNTQzMTYxNWM4NWMwM2JiY2M1Ngo="
export BLOB_STORAGE_CONTAINER_NAME="translator"
export TRANSLATOR_ENDPOINT="https://mytranslator.cognitiveservices.azure.com/"
export TRANSLATOR_KEY="ebe909c559af69d7b285bfb246b214b8"
export TRANSLATOR_REGION="westeurope"
go test -v ./internal/translator
```

## License

This project is licensed under the AGPLv3 License. See the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
