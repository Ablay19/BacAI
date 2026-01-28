package gcpvision

import (
	"context"
	"fmt"
	"io/ioutil"

	vision "cloud.google.com/go/vision/v2/apiv1"
	visionpb "cloud.google.com/go/vision/v2/apiv1/visionpb"
)

// Client represents a Google Cloud Vision API client
type Client struct {
	client *vision.ImageAnnotatorClient
	ctx    context.Context
}

// NewClient creates a new Google Cloud Vision API client
func NewClient() (*Client, error) {
	ctx := context.Background()
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create vision client: %w", err)
	}

	return &Client{
		client: client,
		ctx:    ctx,
	}, nil
}

// TextDetectionResult represents the result of text detection
type TextDetectionResult struct {
	Text         string
	LanguageCode string
}

// DetectTextFromFile detects text in an image file
func (c *Client) DetectTextFromFile(imagePath string) ([]TextDetectionResult, error) {
	// Read image file
	fileBytes, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file: %w", err)
	}

	// Create image from bytes
	image := &visionpb.Image{
		Content: fileBytes,
	}

	// Create request
	req := &visionpb.AnnotateImageRequest{
		Image: image,
		Features: []*visionpb.Feature{
			{
				Type: visionpb.Feature_DOCUMENT_TEXT_DETECTION,
			},
		},
	}

	// Create batch request
	batchReq := &visionpb.BatchAnnotateImagesRequest{
		Requests: []*visionpb.AnnotateImageRequest{req},
	}

	// Make the request
	resp, err := c.client.BatchAnnotateImages(c.ctx, batchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to detect text: %w", err)
	}

	// Process results
	var results []TextDetectionResult
	if len(resp.Responses) > 0 && resp.Responses[0].TextAnnotations != nil {
		// The first annotation contains the full text
		if len(resp.Responses[0].TextAnnotations) > 0 {
			fullText := resp.Responses[0].TextAnnotations[0].Description

			// For detailed word-level analysis, we can parse the full text
			// or use the more detailed Page/Block/Paragraph/Word structure
			results = append(results, TextDetectionResult{
				Text:         fullText,
				LanguageCode: "", // Language detection would require additional processing
			})
		}
	}

	return results, nil
}

// Close closes the Google Cloud Vision client
func (c *Client) Close() error {
	return c.client.Close()
}
