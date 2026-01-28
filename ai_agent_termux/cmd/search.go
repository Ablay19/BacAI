package cmd

import (
	"fmt"
	"os"

	"ai_agent_termux/config"
	"ai_agent_termux/pkg/serpapi"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search documents semantically or on the web",
	Long: `Search documents using semantic similarity with FAISS vector database or perform web search using SerpAPI.
Returns documents most similar to the query or web search results.

Examples:
  ai_agent search "machine learning algorithms"           # Semantic search
  ai_agent search "machine learning algorithms" --web    # Web search using SerpAPI`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]

		// Load config
		cfg := config.LoadConfig()

		slog.Info("Searching documents", "query", query, "output_dir", cfg.OutputDir)

		// Check if we want web search or semantic search
		webSearch, _ := cmd.Flags().GetBool("web")

		if webSearch {
			// Perform web search using SerpAPI
			performWebSearch(query, cfg)
		} else {
			// Perform semantic search using FAISS
			performSemanticSearch(query, cfg)
		}
	},
}

func init() {
	searchCmd.Flags().BoolP("web", "w", false, "Perform web search using SerpAPI instead of semantic search")
}

// performSemanticSearch searches documents using FAISS vector database
func performSemanticSearch(query string, cfg *config.Config) {
	fmt.Printf("Searching for: %s\n", query)
	fmt.Println("Semantic search functionality is not yet implemented.")
	fmt.Println("This feature will use FAISS vector database to find relevant documents.")

	// Placeholder for future implementation
	/*
		// Generate embeddings for query
		embeddings, err := embedding_generator.GenerateEmbeddingsForFile(query, cfg)
		if err != nil {
			slog.Error("Error generating embeddings for query", "error", err)
			fmt.Printf("Error generating embeddings: %v\n", err)
			return
		}

		// Search FAISS index
		faissInterface := vector_db.NewFaissInterface(cfg)
		results, err := faissInterface.Search(embeddings, 5) // Top 5 results
		if err != nil {
			slog.Error("Error searching FAISS index", "error", err)
			fmt.Printf("Error searching: %v\n", err)
			return
		}

		fmt.Printf("Search results for '%s':\n", query)
		// Display results
	*/
}

// performWebSearch searches the web using SerpAPI
func performWebSearch(query string, cfg *config.Config) {
	fmt.Printf("Performing web search for: %s\n", query)

	// Get API key from environment variable
	apiKey := os.Getenv("SERPAPI_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: SERPAPI_API_KEY environment variable not set")
		fmt.Println("Please set your SerpAPI API key as an environment variable:")
		fmt.Println("  export SERPAPI_API_KEY=your_api_key_here")
		return
	}

	// Create SerpAPI client
	client := serpapi.NewClient(apiKey)

	// Perform search
	results, err := client.Search(query)
	if err != nil {
		fmt.Printf("Error performing web search: %v\n", err)
		return
	}

	// Display results
	fmt.Printf("\nFound %d results:\n", len(results.OrganicResults))
	for i, result := range results.OrganicResults {
		fmt.Printf("\n%d. %s\n", i+1, result.Title)
		fmt.Printf("   URL: %s\n", result.Link)
		fmt.Printf("   Snippet: %s\n", result.Snippet)
	}
}
