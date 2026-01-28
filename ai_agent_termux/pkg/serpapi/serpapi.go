package serpapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// SearchResult represents a single search result from SerpAPI
type SearchResult struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
}

// SearchResponse represents the response from SerpAPI
type SearchResponse struct {
	SearchMetadata struct {
		ID          string    `json:"id"`
		Status      string    `json:"status"`
		TotalTime   float64   `json:"total_time"`
		CreatedAt   time.Time `json:"created_at"`
		ProcessedAt time.Time `json:"processed_at"`
	} `json:"search_metadata"`
	SearchParameters struct {
		Engine string `json:"engine"`
		Q      string `json:"q"`
	} `json:"search_parameters"`
	OrganicResults []SearchResult `json:"organic_results"`
}

// Client represents a SerpAPI client
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new SerpAPI client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "https://serpapi.com/search",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Search performs a search using SerpAPI
func (c *Client) Search(query string) (*SearchResponse, error) {
	// Prepare the request URL
	params := url.Values{}
	params.Add("engine", "google")
	params.Add("q", query)
	params.Add("api_key", c.apiKey)

	url := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())

	// Make the HTTP request
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var searchResp SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &searchResp, nil
}

// SearchWithParams performs a search with custom parameters
func (c *Client) SearchWithParams(params map[string]string) (*SearchResponse, error) {
	// Prepare the request URL
	queryParams := url.Values{}
	queryParams.Add("api_key", c.apiKey)

	for key, value := range params {
		queryParams.Add(key, value)
	}

	url := fmt.Sprintf("%s?%s", c.baseURL, queryParams.Encode())

	// Make the HTTP request
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var searchResp SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &searchResp, nil
}
