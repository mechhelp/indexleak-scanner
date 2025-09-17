package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type enterDirectoryArgs struct {
	URL      string `json:"url" description:"URL to enter directory"`
	Page     int    `json:"page,omitempty" description:"Page number (1-based, default: 1)"`
	PageSize int    `json:"page_size,omitempty" description:"Number of items per page (default: 30)"`
}

type directoryEntry struct {
	Name string `json:"name" description:"Name of the entry"`
	Type string `json:"type" description:"Type: DIRECTORY or FILE"`
	Link string `json:"link" description:"Full URL link"`
}

type scanResult struct {
	URL           string
	Type          string
	FileExtension string
}

// Global variable to store scan results
var scanResults []scanResult

func getDefaultPageParams(params *enterDirectoryArgs) (page, pageSize int) {
	page = params.Page
	if page < 1 {
		page = 1
	}

	pageSize = params.PageSize
	if pageSize < 1 {
		pageSize = 30
	}

	return page, pageSize
}

func enterDirectory(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[enterDirectoryArgs]) (*mcp.CallToolResult, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", params.Arguments.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	htmlContent := string(body)

	// Check for "Index of" string
	if !strings.Contains(htmlContent, "Index of") {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "This page is not a directory listing (Index of not found)"},
			},
		}, nil
	}

	// Parse links in HTML content
	linkRegex := regexp.MustCompile(`<a href="([^"]+)">([^<]+)</a>`)
	matches := linkRegex.FindAllStringSubmatch(htmlContent, -1)

	var entries []directoryEntry
	baseURL := strings.TrimSuffix(params.Arguments.URL, "/")

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		href := match[1]
		name := strings.TrimSpace(match[2])

		// Exclude ../ links
		if href == "../" {
			continue
		}

		// Create full URL
		fullURL := baseURL + "/" + strings.TrimPrefix(href, "/")

		// Check if it's a directory or file
		entryType := "FILE"
		if strings.HasSuffix(href, "/") {
			entryType = "DIRECTORY"
		}

		entries = append(entries, directoryEntry{
			Name: name,
			Type: entryType,
			Link: fullURL,
		})
	}

	// Apply pagination
	page, pageSize := getDefaultPageParams(&params.Arguments)
	totalEntries := len(entries)
	totalPages := (totalEntries + pageSize - 1) / pageSize

	if page > totalPages {
		page = totalPages
	}

	startIdx := (page - 1) * pageSize
	endIdx := startIdx + pageSize
	if endIdx > totalEntries {
		endIdx = totalEntries
	}

	// Get entries for current page
	pageEntries := entries
	if startIdx < totalEntries {
		pageEntries = entries[startIdx:endIdx]
	} else {
		pageEntries = []directoryEntry{}
	}

	// Format results
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Directory contents (%s):\n", params.Arguments.URL))
	result.WriteString(fmt.Sprintf("\nPage %d of %d (Total entries: %d, Entries per page: %d)\n\n",
		page, totalPages, totalEntries, pageSize))

	for _, entry := range pageEntries {
		result.WriteString(fmt.Sprintf("%s: %s (%s)\n", entry.Type, entry.Name, entry.Link))

		// Store scan result for later export
		extension := ""
		if entry.Type == "FILE" {
			if idx := strings.LastIndex(entry.Name, "."); idx != -1 {
				extension = entry.Name[idx:]
			}
		}

		scanResults = append(scanResults, scanResult{
			URL:           entry.Link,
			Type:          entry.Type,
			FileExtension: extension,
		})
	}

	if totalPages > 1 {
		result.WriteString(fmt.Sprintf("\nUse page parameter to navigate between pages (1-%d)", totalPages))
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result.String()},
		},
	}, nil
}

func main() {
	flag.Parse()

	server := mcp.NewServer(&mcp.Implementation{Name: "indexleak-scanner"}, nil)

	// Tools
	mcp.AddTool(server, &mcp.Tool{Name: "enter_directory", Description: "enter directory and list contents from Index of pages"}, enterDirectory)

	t := mcp.NewLoggingTransport(mcp.NewStdioTransport(), os.Stderr)
	if err := server.Run(context.Background(), t); err != nil {
		log.Printf("Server failed: %v", err)
	}
}
