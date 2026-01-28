# Google Cloud Vision API Package

This package provides a simplified interface to the Google Cloud Vision API for text detection in images.

## Features

- Text detection from image files
- Automatic handling of Google Cloud authentication
- Simple API for extracting text from documents and images

## Prerequisites

To use this package, you need:

1. A Google Cloud Platform account
2. A project with the Vision API enabled
3. Authentication credentials (service account key or default credentials)

## Setup

1. Enable the Vision API in your Google Cloud Console
2. Create a service account and download the JSON key file
3. Set the `GOOGLE_APPLICATION_CREDENTIALS` environment variable to point to your key file:

```bash
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/your/service-account-key.json
```

Alternatively, you can use default credentials if running on Google Cloud Platform.

## Usage

```go
import "ai_agent_termux/pkg/gcpvision"

// Create a new client
client, err := gcpvision.NewClient()
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// Detect text in an image
results, err := client.DetectTextFromFile("/path/to/image.jpg")
if err != nil {
    log.Fatal(err)
}

// Print results
for _, result := range results {
    fmt.Printf("Text: %s\n", result.Text)
}
```

## Limitations

This package currently only supports document text detection. Future versions may include support for:
- OCR with layout preservation
- Language detection
- Handwriting recognition
- Barcode/QR code detection