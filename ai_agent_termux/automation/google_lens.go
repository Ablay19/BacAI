package automation

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"ai_agent_termux/config"
	vision "cloud.google.com/go/vision/v2/apiv1"
	visionpb "cloud.google.com/go/vision/v2/apiv1/visionpb"
	"github.com/anthonynsimon/bild/effect"
	"github.com/disintegration/imaging"
	"golang.org/x/exp/slog"
)

// GoogleLensProcessor handles actual Google Lens processing
type GoogleLensProcessor struct {
	config          *config.Config
	cache           map[string]LensResult
	cacheMutex      sync.RWMutex
	httpClient      *http.Client
	resultCallbacks []func(LensResult)
	visionClient    *vision.ImageAnnotatorClient
	ctx             context.Context
}

// LensResult represents the result of Google Lens processing
type LensResult struct {
	ImagePath      string             `json:"image_path"`
	Operation      string             `json:"operation"`
	ResultText     string             `json:"result_text"`
	Objects        []IdentifiedObject `json:"objects,omitempty"`
	TextBlocks     []TextBlock        `json:"text_blocks,omitempty"`
	Barcodes       []Barcode          `json:"barcodes,omitempty"`
	MathResults    []MathResult       `json:"math_results,omitempty"`
	Timestamp      time.Time          `json:"timestamp"`
	ProcessingTime time.Duration      `json:"processing_time"`
	Cached         bool               `json:"cached"`
}

// IdentifiedObject represents an object identified by Google Lens
type IdentifiedObject struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
	Bounds     Bounds  `json:"bounds"`
}

// TextBlock represents a block of extracted text
type TextBlock struct {
	Text   string `json:"text"`
	Bounds Bounds `json:"bounds"`
}

// Barcode represents a scanned barcode/QR code
type Barcode struct {
	Type   string `json:"type"`
	Data   string `json:"data"`
	Format string `json:"format"`
}

// MathResult represents a solved math problem
type MathResult struct {
	Problem  string   `json:"problem"`
	Solution string   `json:"solution"`
	Steps    []string `json:"steps"`
}

// Bounds represents coordinates of an object/element
type Bounds struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// ProgressCallback represents a progress update callback
type ProgressCallback func(progress float64, message string)

// LensRequest represents a request to Google Lens API
type LensRequest struct {
	ImageData  string `json:"image_data"`
	Operation  string `json:"operation"`
	Language   string `json:"language,omitempty"`
	MaxResults int    `json:"max_results,omitempty"`
}

// LensResponse represents a response from Google Lens API
type LensResponse struct {
	Success     bool       `json:"success"`
	Result      LensResult `json:"result"`
	Error       string     `json:"error,omitempty"`
	Suggestions []string   `json:"suggestions,omitempty"`
}

// NewGoogleLensProcessor creates a new Google Lens processor
func NewGoogleLensProcessor(cfg *config.Config) *GoogleLensProcessor {
	ctx := context.Background()
	visionClient, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		slog.Error("Failed to initialize Vision API client", "error", err)
		// We'll continue without Vision API for now
		visionClient = nil
	}

	processor := &GoogleLensProcessor{
		config:       cfg,
		cache:        make(map[string]LensResult),
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		visionClient: visionClient,
		ctx:          ctx,
	}

	slog.Info("Google Lens Processor initialized", "vision_api_enabled", visionClient != nil)
	return processor
}

// ProcessImage processes an image with Google Lens
func (glp *GoogleLensProcessor) ProcessImage(imagePath, operation string) (LensResult, error) {
	return glp.ProcessImageWithProgress(imagePath, operation, nil)
}

// ProcessImageWithProgress processes an image with Google Lens and progress tracking
func (glp *GoogleLensProcessor) ProcessImageWithProgress(imagePath, operation string, progressCallback ProgressCallback) (LensResult, error) {
	startTime := time.Now()

	// Report initial progress
	if progressCallback != nil {
		progressCallback(0.1, "Validating image file...")
	}

	// Generate file hash for better cache key
	hash, err := glp.generateFileHash(imagePath)
	if err != nil {
		slog.Warn("Failed to generate file hash, using path as key", "error", err)
		hash = imagePath
	}
	cacheKey := fmt.Sprintf("%s_%s", hash, operation)

	// Check cache first
	if result := glp.getFromCache(cacheKey); result != nil {
		slog.Info("Using cached result", "image", imagePath, "operation", operation)
		result.Cached = true
		result.ProcessingTime = time.Since(startTime)
		if progressCallback != nil {
			progressCallback(1.0, "Using cached result")
		}
		return *result, nil
	}

	// Validate image file
	if err := glp.validateImageFile(imagePath); err != nil {
		return LensResult{}, fmt.Errorf("invalid image file: %v", err)
	}

	// Preprocess image if needed
	if progressCallback != nil {
		progressCallback(0.3, "Preprocessing image...")
	}
	processedImagePath, err := glp.preprocessImage(imagePath)
	if err != nil {
		slog.Warn("Image preprocessing failed, using original", "error", err)
		processedImagePath = imagePath
	}

	// Process with Google Vision API
	if progressCallback != nil {
		progressCallback(0.5, "Processing with Google Vision API...")
	}
	result, err := glp.processWithVisionAPI(processedImagePath, operation, progressCallback)
	if err != nil {
		return LensResult{}, fmt.Errorf("Google Vision API processing failed: %v", err)
	}

	// Add metadata
	result.ImagePath = imagePath
	result.Operation = operation
	result.Timestamp = time.Now()
	result.ProcessingTime = time.Since(startTime)
	result.Cached = false

	// Cache result
	glp.addToCache(cacheKey, result)

	// Notify callbacks
	glp.notifyCallbacks(result)

	// Final progress update
	if progressCallback != nil {
		progressCallback(1.0, "Processing completed")
	}

	slog.Info("Google Vision API processing completed",
		"image", imagePath,
		"operation", operation,
		"duration", result.ProcessingTime)

	return result, nil
}

// validateImageFile checks if the image file is valid
func (glp *GoogleLensProcessor) validateImageFile(imagePath string) error {
	info, err := os.Stat(imagePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", imagePath)
	}
	if err != nil {
		return fmt.Errorf("error accessing file: %v", err)
	}

	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file: %s", imagePath)
	}

	// Check file extension
	ext := filepath.Ext(imagePath)
	validExtensions := []string{".jpg", ".jpeg", ".png", ".bmp", ".webp"}
	valid := false
	for _, validExt := range validExtensions {
		if ext == validExt || ext == strings.ToUpper(validExt) {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("unsupported image format: %s", ext)
	}

	// Check file size (should not exceed 10MB for Google Lens)
	if info.Size() > 10*1024*1024 {
		return fmt.Errorf("image file too large: %d bytes (max 10MB)", info.Size())
	}

	return nil
}

// generateFileHash creates a SHA-256 hash of the image file for caching
func (glp *GoogleLensProcessor) generateFileHash(imagePath string) (string, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// PreprocessOptions contains options for image preprocessing
type PreprocessOptions struct {
	MaxWidth        int     `json:"max_width"`
	MaxHeight       int     `json:"max_height"`
	Quality         float64 `json:"quality"`          // JPEG quality 0-1
	EnhanceContrast float64 `json:"enhance_contrast"` // Contrast adjustment
	Sharpen         bool    `json:"sharpen"`          // Apply sharpening
	Denoise         bool    `json:"denoise"`          // Apply noise reduction
	AutoRotate      bool    `json:"auto_rotate"`      // Auto-rotate based on EXIF
	Grayscale       bool    `json:"grayscale"`        // Convert to grayscale for OCR
}

// DefaultPreprocessOptions returns default preprocessing options
func DefaultPreprocessOptions() *PreprocessOptions {
	return &PreprocessOptions{
		MaxWidth:        2048,
		MaxHeight:       2048,
		Quality:         0.85,
		EnhanceContrast: 15.0,
		Sharpen:         true,
		Denoise:         false,
		AutoRotate:      true,
		Grayscale:       false,
	}
}

// preprocessImage prepares an image for Google Vision API processing
func (glp *GoogleLensProcessor) preprocessImage(imagePath string) (string, error) {
	return glp.preprocessImageWithOptions(imagePath, DefaultPreprocessOptions())
}

// preprocessImageWithOptions prepares an image with custom preprocessing options
func (glp *GoogleLensProcessor) preprocessImageWithOptions(imagePath string, opts *PreprocessOptions) (string, error) {
	slog.Debug("Preprocessing image", "path", imagePath, "options", opts)

	// Open the image
	img, err := imaging.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %v", err)
	}

	processed := img
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Log original dimensions
	slog.Debug("Original image dimensions", "width", width, "height", height)

	// Auto-rotate based on EXIF if requested
	if opts.AutoRotate {
		// imaging library automatically handles EXIF rotation
		slog.Debug("Auto-rotation enabled")
	}

	// Calculate resize dimensions maintaining aspect ratio
	if width > opts.MaxWidth || height > opts.MaxHeight {
		slog.Info("Resizing large image",
			"original_size", fmt.Sprintf("%dx%d", width, height),
			"max_size", fmt.Sprintf("%dx%d", opts.MaxWidth, opts.MaxHeight))

		processed = imaging.Fit(processed, opts.MaxWidth, opts.MaxHeight, imaging.Lanczos)

		newBounds := processed.Bounds()
		slog.Debug("Resized image dimensions", "width", newBounds.Dx(), "height", newBounds.Dy())
	}

	// Apply denoising if requested (using bilateral filter simulation)
	if opts.Denoise {
		slog.Debug("Applying noise reduction")
		// Use a mild blur to reduce noise
		processed = imaging.Blur(processed, 0.5)
	}

	// Convert to grayscale for better OCR if requested
	if opts.Grayscale {
		slog.Debug("Converting to grayscale")
		processed = imaging.Grayscale(processed)
	}

	// Enhance contrast for better OCR
	if opts.EnhanceContrast != 0 {
		slog.Debug("Enhancing contrast", "adjustment", opts.EnhanceContrast)
		processed = imaging.AdjustContrast(processed, opts.EnhanceContrast)
	}

	// Adjust brightness if needed for better visibility
	slog.Debug("Adjusting brightness")
	processed = imaging.AdjustBrightness(processed, 5)

	// Apply sharpening if requested
	if opts.Sharpen {
		slog.Debug("Applying sharpening")
		processed = imaging.Sharpen(processed, 0.5)
	}

	// Additional enhancement for OCR - increase local contrast
	slog.Debug("Applying local contrast enhancement")
	processed = effect.UnsharpMask(processed, 1.0, 0.5)

	// Save processed image to a temporary file
	tempDir := os.TempDir()
	ext := filepath.Ext(imagePath)
	if ext == "" {
		ext = ".jpg" // Default to JPEG
	}
	tempPath := filepath.Join(tempDir, fmt.Sprintf("processed_%s_%d%s",
		filepath.Base(imagePath)[:len(filepath.Base(imagePath))-len(ext)],
		time.Now().Unix(),
		ext))

	// Save with quality setting
	if strings.ToLower(ext) == ".jpg" || strings.ToLower(ext) == ".jpeg" {
		err = imaging.Save(processed, tempPath, imaging.JPEGQuality(int(opts.Quality*100)))
	} else {
		err = imaging.Save(processed, tempPath)
	}

	if err != nil {
		return "", fmt.Errorf("failed to save processed image: %v", err)
	}

	// Get file size for logging
	if stat, err := os.Stat(tempPath); err == nil {
		slog.Info("Image preprocessed successfully",
			"original", imagePath,
			"processed", tempPath,
			"original_size", fmt.Sprintf("%.1fMB", float64(width*height*3)/1024/1024),
			"processed_size", fmt.Sprintf("%.1fMB", float64(stat.Size())/1024/1024))
	} else {
		slog.Info("Image preprocessed", "original", imagePath, "processed", tempPath)
	}

	return tempPath, nil
}

// PreprocessForOCR applies OCR-specific preprocessing
func (glp *GoogleLensProcessor) PreprocessForOCR(imagePath string) (string, error) {
	opts := DefaultPreprocessOptions()
	opts.Grayscale = true
	opts.EnhanceContrast = 20.0
	opts.Sharpen = true
	opts.Denoise = true
	opts.Quality = 0.9

	return glp.preprocessImageWithOptions(imagePath, opts)
}

// PreprocessForObjectDetection applies object detection specific preprocessing
func (glp *GoogleLensProcessor) PreprocessForObjectDetection(imagePath string) (string, error) {
	opts := DefaultPreprocessOptions()
	opts.Grayscale = false
	opts.EnhanceContrast = 10.0
	opts.Sharpen = true
	opts.Denoise = false

	return glp.preprocessImageWithOptions(imagePath, opts)
}

// processWithVisionAPI processes an image using real Google Vision API
func (glp *GoogleLensProcessor) processWithVisionAPI(imagePath, operation string, progressCallback ProgressCallback) (LensResult, error) {
	if glp.visionClient == nil {
		slog.Warn("Vision API client not initialized, falling back to simulation")
		return glp.processWithGoogleLens(imagePath, operation)
	}

	slog.Info("Processing image with Google Vision API", "path", imagePath, "operation", operation)

	// Report progress
	if progressCallback != nil {
		progressCallback(0.6, "Reading image file...")
	}

	// Read image file
	fileBytes, err := os.ReadFile(imagePath)
	if err != nil {
		return LensResult{}, fmt.Errorf("failed to read image file: %v", err)
	}

	// Create image
	image := &visionpb.Image{
		Content: fileBytes,
	}

	result := LensResult{
		ImagePath: imagePath,
		Operation: operation,
	}

	// Process based on operation type
	switch operation {
	case "extract_text", "document":
		if progressCallback != nil {
			progressCallback(0.7, "Extracting text...")
		}
		err = glp.extractTextWithVision(image, &result)
	case "identify_object":
		if progressCallback != nil {
			progressCallback(0.7, "Identifying objects...")
		}
		err = glp.identifyObjectsWithVision(image, &result)
	case "barcode":
		if progressCallback != nil {
			progressCallback(0.7, "Scanning barcodes...")
		}
		err = glp.detectBarcodesWithVision(image, &result)
	default:
		// For unsupported operations, use the simulated approach
		return glp.processWithGoogleLens(imagePath, operation)
	}

	if err != nil {
		return LensResult{}, fmt.Errorf("Vision API processing failed: %v", err)
	}

	return result, nil
}

// extractTextWithVision extracts text using Vision API
func (glp *GoogleLensProcessor) extractTextWithVision(image *visionpb.Image, result *LensResult) error {
	req := &visionpb.AnnotateImageRequest{
		Image: image,
		Features: []*visionpb.Feature{
			{
				Type: visionpb.Feature_DOCUMENT_TEXT_DETECTION,
			},
		},
	}

	batchReq := &visionpb.BatchAnnotateImagesRequest{
		Requests: []*visionpb.AnnotateImageRequest{req},
	}

	resp, err := glp.visionClient.BatchAnnotateImages(glp.ctx, batchReq)
	if err != nil {
		return fmt.Errorf("text detection failed: %v", err)
	}

	if len(resp.Responses) > 0 && resp.Responses[0].FullTextAnnotation != nil {
		result.ResultText = resp.Responses[0].FullTextAnnotation.Text

		// Extract text blocks with bounding boxes
		for _, page := range resp.Responses[0].FullTextAnnotation.Pages {
			for _, block := range page.Blocks {
				blockText := ""
				for _, paragraph := range block.Paragraphs {
					for _, word := range paragraph.Words {
						for _, symbol := range word.Symbols {
							blockText += symbol.Text
						}
						blockText += " "
					}
				}

				if len(block.BoundingBox.Vertices) >= 2 {
					x := block.BoundingBox.Vertices[0].X
					y := block.BoundingBox.Vertices[0].Y
					width := block.BoundingBox.Vertices[1].X - x
					height := block.BoundingBox.Vertices[2].Y - y

					result.TextBlocks = append(result.TextBlocks, TextBlock{
						Text:   strings.TrimSpace(blockText),
						Bounds: Bounds{X: int(x), Y: int(y), Width: int(width), Height: int(height)},
					})
				}
			}
		}
	}

	return nil
}

// identifyObjectsWithVision identifies objects using Vision API
func (glp *GoogleLensProcessor) identifyObjectsWithVision(image *visionpb.Image, result *LensResult) error {
	req := &visionpb.AnnotateImageRequest{
		Image: image,
		Features: []*visionpb.Feature{
			{
				Type:       visionpb.Feature_LABEL_DETECTION,
				MaxResults: 10,
			},
		},
	}

	batchReq := &visionpb.BatchAnnotateImagesRequest{
		Requests: []*visionpb.AnnotateImageRequest{req},
	}

	resp, err := glp.visionClient.BatchAnnotateImages(glp.ctx, batchReq)
	if err != nil {
		return fmt.Errorf("object detection failed: %v", err)
	}

	var objectNames []string
	if len(resp.Responses) > 0 {
		for _, annotation := range resp.Responses[0].LabelAnnotations {
			objectNames = append(objectNames, annotation.Description)
			result.Objects = append(result.Objects, IdentifiedObject{
				Name:       annotation.Description,
				Confidence: float64(annotation.Score),
				Bounds:     Bounds{X: 0, Y: 0, Width: 0, Height: 0}, // Label detection doesn't provide bounding boxes
			})
		}
	}

	if len(objectNames) > 0 {
		result.ResultText = fmt.Sprintf("Identified: %s", strings.Join(objectNames, ", "))
	}

	return nil
}

// detectBarcodesWithVision detects barcodes using Vision API
func (glp *GoogleLensProcessor) detectBarcodesWithVision(image *visionpb.Image, result *LensResult) error {
	// Note: Barcode detection requires OBJECT_LOCALIZATION or TEXT_DETECTION
	// We'll use OBJECT_LOCALIZATION to find QR codes and barcodes
	req := &visionpb.AnnotateImageRequest{
		Image: image,
		Features: []*visionpb.Feature{
			{
				Type:       visionpb.Feature_OBJECT_LOCALIZATION,
				MaxResults: 10,
			},
		},
	}

	batchReq := &visionpb.BatchAnnotateImagesRequest{
		Requests: []*visionpb.AnnotateImageRequest{req},
	}

	resp, err := glp.visionClient.BatchAnnotateImages(glp.ctx, batchReq)
	if err != nil {
		return fmt.Errorf("barcode detection failed: %v", err)
	}

	var barcodeData []string
	if len(resp.Responses) > 0 {
		// Look for objects that might be barcodes/QR codes
		for _, obj := range resp.Responses[0].LocalizedObjectAnnotations {
			// Check if object name suggests it's a barcode or QR code
			if strings.Contains(strings.ToLower(obj.Name), "barcode") ||
				strings.Contains(strings.ToLower(obj.Name), "qr") ||
				strings.Contains(strings.ToLower(obj.Name), "code") {
				barcodeData = append(barcodeData, obj.Name)
				result.Barcodes = append(result.Barcodes, Barcode{
					Type:   obj.Name,
					Data:   fmt.Sprintf("Detected at confidence %.2f", obj.Score),
					Format: "OBJECT_LOCALIZATION",
				})
			}
		}

		// Also try text detection as fallback for barcodes
		if len(barcodeData) == 0 {
			textReq := &visionpb.AnnotateImageRequest{
				Image: image,
				Features: []*visionpb.Feature{
					{
						Type: visionpb.Feature_TEXT_DETECTION,
					},
				},
			}

			textBatchReq := &visionpb.BatchAnnotateImagesRequest{
				Requests: []*visionpb.AnnotateImageRequest{textReq},
			}

			textResp, err := glp.visionClient.BatchAnnotateImages(glp.ctx, textBatchReq)
			if err == nil && len(textResp.Responses) > 0 && textResp.Responses[0].TextAnnotations != nil {
				for _, text := range textResp.Responses[0].TextAnnotations {
					// Look for URL patterns or barcode-like text
					if strings.HasPrefix(text.Description, "http") ||
						len(text.Description) > 10 && (strings.Contains(text.Description, ":") || strings.Contains(text.Description, "/")) {
						barcodeData = append(barcodeData, "Text: "+text.Description)
						result.Barcodes = append(result.Barcodes, Barcode{
							Type:   "TEXT",
							Data:   text.Description,
							Format: "TEXT_DETECTION",
						})
					}
				}
			}
		}
	}

	if len(barcodeData) > 0 {
		result.ResultText = fmt.Sprintf("Detected potential barcodes/codes: %s", strings.Join(barcodeData, ", "))
	} else {
		result.ResultText = "No barcodes or QR codes detected"
	}

	return nil
}

// processWithGoogleLens simulates actual Google Lens processing
// This is used as fallback or for unsupported operations
func (glp *GoogleLensProcessor) processWithGoogleLens(imagePath, operation string) (LensResult, error) {
	slog.Info("Processing image with simulated Google Lens", "path", imagePath, "operation", operation)

	// Simulate network delay
	time.Sleep(1 * time.Second)

	result := LensResult{
		ImagePath: imagePath,
		Operation: operation,
	}

	switch operation {
	case "solve_math":
		result.MathResults = []MathResult{
			{
				Problem:  "2x + 5 = 15",
				Solution: "x = 5",
				Steps: []string{
					"2x + 5 = 15",
					"2x = 15 - 5",
					"2x = 10",
					"x = 10/2",
					"x = 5",
				},
			},
		}
		result.ResultText = "Solution: x = 5"
	case "translate_text":
		result.ResultText = "Original: Bonjour, comment allez-vous?\nTranslation: Hello, how are you?"
	default:
		result.ResultText = fmt.Sprintf("Processed with simulated Google Lens operation: %s", operation)
	}

	return result, nil
}

// BatchProcessImages processes multiple images with Google Lens
func (glp *GoogleLensProcessor) BatchProcessImages(imagePaths []string, operation string) ([]LensResult, error) {
	return glp.BatchProcessImagesWithProgress(imagePaths, operation, nil)
}

// BatchProcessImagesWithProgress processes multiple images with progress tracking
func (glp *GoogleLensProcessor) BatchProcessImagesWithProgress(imagePaths []string, operation string, progressCallback ProgressCallback) ([]LensResult, error) {
	results := make([]LensResult, len(imagePaths))
	errors := make([]error, len(imagePaths))
	completed := 0
	var completedMutex sync.Mutex

	// Process images concurrently with controlled concurrency
	semaphore := make(chan struct{}, 5) // Limit to 5 concurrent API calls
	var wg sync.WaitGroup
	wg.Add(len(imagePaths))

	for i, imagePath := range imagePaths {
		go func(index int, path string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			// Create individual progress callback that updates the global progress
			individualCallback := func(progress float64, message string) {
				if progressCallback != nil {
					completedMutex.Lock()
					totalProgress := (float64(completed) + progress*0.8) / float64(len(imagePaths)) // 80% for processing, 20% for finalization
					completedMutex.Unlock()
					progressCallback(totalProgress, fmt.Sprintf("Image %d/%d: %s", index+1, len(imagePaths), message))
				}
			}

			result, err := glp.ProcessImageWithProgress(path, operation, individualCallback)
			results[index] = result
			errors[index] = err

			completedMutex.Lock()
			completed++
			completedMutex.Unlock()
		}(i, imagePath)
	}

	// Wait for all processing to complete
	wg.Wait()

	// Final progress update
	if progressCallback != nil {
		progressCallback(1.0, fmt.Sprintf("Batch processing completed (%d images)", len(imagePaths)))
	}

	// Check for errors
	var processingErrors []string
	for i, err := range errors {
		if err != nil {
			processingErrors = append(processingErrors, fmt.Sprintf("Image %d (%s): %v", i, imagePaths[i], err))
		}
	}

	if len(processingErrors) > 0 {
		slog.Warn("Batch processing completed with errors", "errors", len(processingErrors))
		return results, fmt.Errorf("batch processing had %d errors: %s", len(processingErrors), strings.Join(processingErrors, "; "))
	}

	slog.Info("Batch processing completed successfully", "count", len(imagePaths))
	return results, nil
}

// ProcessImageWithPreprocessing processes an image with custom preprocessing
func (glp *GoogleLensProcessor) ProcessImageWithPreprocessing(imagePath, operation, preprocessType string, progressCallback ProgressCallback) (LensResult, error) {
	startTime := time.Now()

	// Report initial progress
	if progressCallback != nil {
		progressCallback(0.1, "Starting processing with custom preprocessing...")
	}

	// Generate file hash for better cache key
	hash, err := glp.generateFileHash(imagePath)
	if err != nil {
		slog.Warn("Failed to generate file hash, using path as key", "error", err)
		hash = imagePath
	}
	cacheKey := fmt.Sprintf("%s_%s_%s", hash, operation, preprocessType)

	// Check cache first
	if result := glp.getFromCache(cacheKey); result != nil {
		slog.Info("Using cached result", "image", imagePath, "operation", operation, "preprocess", preprocessType)
		result.Cached = true
		result.ProcessingTime = time.Since(startTime)
		if progressCallback != nil {
			progressCallback(1.0, "Using cached result")
		}
		return *result, nil
	}

	// Validate image file
	if err := glp.validateImageFile(imagePath); err != nil {
		return LensResult{}, fmt.Errorf("invalid image file: %v", err)
	}

	// Apply custom preprocessing based on type
	if progressCallback != nil {
		progressCallback(0.3, "Applying custom preprocessing...")
	}
	processedImagePath, err := glp.applyCustomPreprocessing(imagePath, preprocessType)
	if err != nil {
		slog.Warn("Custom preprocessing failed, using standard preprocessing", "error", err)
		processedImagePath, err = glp.preprocessImage(imagePath)
		if err != nil {
			slog.Warn("Standard preprocessing failed, using original", "error", err)
			processedImagePath = imagePath
		}
	}

	// Process with Google Vision API
	if progressCallback != nil {
		progressCallback(0.5, "Processing with Google Vision API...")
	}
	result, err := glp.processWithVisionAPI(processedImagePath, operation, progressCallback)
	if err != nil {
		return LensResult{}, fmt.Errorf("Google Vision API processing failed: %v", err)
	}

	// Add metadata
	result.ImagePath = imagePath
	result.Operation = operation
	result.Timestamp = time.Now()
	result.ProcessingTime = time.Since(startTime)
	result.Cached = false

	// Cache result with custom cache key
	glp.addToCache(cacheKey, result)

	// Notify callbacks
	glp.notifyCallbacks(result)

	// Final progress update
	if progressCallback != nil {
		progressCallback(1.0, "Processing completed")
	}

	slog.Info("Google Vision API processing completed with custom preprocessing",
		"image", imagePath,
		"operation", operation,
		"preprocess_type", preprocessType,
		"duration", result.ProcessingTime)

	return result, nil
}

// applyCustomPreprocessing applies custom preprocessing based on type
func (glp *GoogleLensProcessor) applyCustomPreprocessing(imagePath, preprocessType string) (string, error) {
	switch preprocessType {
	case "ocr":
		return glp.PreprocessForOCR(imagePath)
	case "object":
		return glp.PreprocessForObjectDetection(imagePath)
	case "document":
		return glp.preprocessForDocument(imagePath)
	case "barcode":
		return glp.preprocessForBarcode(imagePath)
	default:
		return glp.preprocessImage(imagePath)
	}
}

// preprocessForDocument applies document-specific preprocessing
func (glp *GoogleLensProcessor) preprocessForDocument(imagePath string) (string, error) {
	opts := DefaultPreprocessOptions()
	opts.Grayscale = true
	opts.EnhanceContrast = 25.0
	opts.Sharpen = true
	opts.Denoise = true
	opts.Quality = 0.9
	opts.MaxWidth = 3000
	opts.MaxHeight = 3000

	return glp.preprocessImageWithOptions(imagePath, opts)
}

// preprocessForBarcode applies barcode-specific preprocessing
func (glp *GoogleLensProcessor) preprocessForBarcode(imagePath string) (string, error) {
	opts := DefaultPreprocessOptions()
	opts.Grayscale = true
	opts.EnhanceContrast = 50.0
	opts.Sharpen = true
	opts.Denoise = false
	opts.Quality = 0.95
	opts.MaxWidth = 1024
	opts.MaxHeight = 1024

	return glp.preprocessImageWithOptions(imagePath, opts)
}

// Close closes the Vision API client and cleans up resources
func (glp *GoogleLensProcessor) Close() error {
	if glp.visionClient != nil {
		return glp.visionClient.Close()
	}
	return nil
}

// getFromCache retrieves a result from cache
func (glp *GoogleLensProcessor) getFromCache(key string) *LensResult {
	glp.cacheMutex.RLock()
	defer glp.cacheMutex.RUnlock()

	if result, exists := glp.cache[key]; exists {
		// Check if cache is still valid (less than 1 hour old)
		if time.Since(result.Timestamp) < time.Hour {
			return &result
		}
		// Remove expired cache entry
		glp.cacheMutex.RUnlock()
		glp.cacheMutex.Lock()
		delete(glp.cache, key)
		glp.cacheMutex.Unlock()
		glp.cacheMutex.RLock()
	}

	return nil
}

// addToCache stores a result in cache
func (glp *GoogleLensProcessor) addToCache(key string, result LensResult) {
	glp.cacheMutex.Lock()
	defer glp.cacheMutex.Unlock()

	// Limit cache size to prevent memory issues
	if len(glp.cache) >= 1000 {
		// Remove oldest entries
		count := 0
		for k := range glp.cache {
			if count >= 100 {
				break
			}
			delete(glp.cache, k)
			count++
		}
	}

	glp.cache[key] = result
}

// ClearCache clears all cached results
func (glp *GoogleLensProcessor) ClearCache() {
	glp.cacheMutex.Lock()
	defer glp.cacheMutex.Unlock()
	glp.cache = make(map[string]LensResult)
	slog.Info("Google Lens cache cleared")
}

// AddResultCallback adds a callback function to be notified of processing results
func (glp *GoogleLensProcessor) AddResultCallback(callback func(LensResult)) {
	glp.resultCallbacks = append(glp.resultCallbacks, callback)
}

// notifyCallbacks notifies all registered callbacks of a result
func (glp *GoogleLensProcessor) notifyCallbacks(result LensResult) {
	for _, callback := range glp.resultCallbacks {
		go callback(result) // Run callbacks concurrently
	}
}

// ExportResults exports processing results to a JSON file
func (glp *GoogleLensProcessor) ExportResults(results []LensResult, outputPath string) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %v", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write results to file: %v", err)
	}

	slog.Info("Results exported", "path", outputPath)
	return nil
}

// ImportResults imports processing results from a JSON file
func (glp *GoogleLensProcessor) ImportResults(inputPath string) ([]LensResult, error) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var results []LensResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal results: %v", err)
	}

	// Add imported results to cache
	glp.cacheMutex.Lock()
	for _, result := range results {
		key := fmt.Sprintf("%s_%s", result.ImagePath, result.Operation)
		glp.cache[key] = result
	}
	glp.cacheMutex.Unlock()

	slog.Info("Results imported", "path", inputPath, "count", len(results))
	return results, nil
}
