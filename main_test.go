package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Mock HTML content for different test scenarios
const (
	validDirectoryHTML = `<!DOCTYPE html>
<html>
<head><title>Index of /test</title></head>
<body>
<h1>Index of /test</h1>
<hr>
<pre>
<a href="../">../</a>
<a href="file1.txt">file1.txt</a>                          01-Jan-2024 10:00               1024
<a href="file2.pdf">file2.pdf</a>                          01-Jan-2024 11:00               2048
<a href="subdirectory/">subdirectory/</a>                   01-Jan-2024 12:00                  -
<a href="image.jpg">image.jpg</a>                          01-Jan-2024 13:00               5120
</pre>
</body>
</html>`

	nonDirectoryHTML = `<!DOCTYPE html>
<html>
<head><title>Welcome</title></head>
<body>
<h1>Welcome to our website</h1>
<p>This is not a directory listing page.</p>
</body>
</html>`

	largeDirectoryHTML = `<!DOCTYPE html>
<html>
<head><title>Index of /large</title></head>
<body>
<h1>Index of /large</h1>
<hr>
<pre>
<a href="../">../</a>
%s
</pre>
</body>
</html>`
)

// generateLargeDirectoryHTML creates HTML with many entries for pagination testing
func generateLargeDirectoryHTML(count int) string {
	var entries []string
	for i := 1; i <= count; i++ {
		entries = append(entries, fmt.Sprintf(`<a href="file%d.txt">file%d.txt</a>                          01-Jan-2024 10:00               1024`, i, i))
	}
	return fmt.Sprintf(largeDirectoryHTML, strings.Join(entries, "\n"))
}

func TestGetDefaultPageParams(t *testing.T) {
	tests := []struct {
		name         string
		params       enterDirectoryArgs
		expectedPage int
		expectedSize int
	}{
		{
			name:         "default values",
			params:       enterDirectoryArgs{},
			expectedPage: 1,
			expectedSize: 30,
		},
		{
			name:         "valid page and size",
			params:       enterDirectoryArgs{Page: 2, PageSize: 10},
			expectedPage: 2,
			expectedSize: 10,
		},
		{
			name:         "zero page defaults to 1",
			params:       enterDirectoryArgs{Page: 0, PageSize: 15},
			expectedPage: 1,
			expectedSize: 15,
		},
		{
			name:         "negative page defaults to 1",
			params:       enterDirectoryArgs{Page: -5, PageSize: 20},
			expectedPage: 1,
			expectedSize: 20,
		},
		{
			name:         "zero page size defaults to 30",
			params:       enterDirectoryArgs{Page: 3, PageSize: 0},
			expectedPage: 3,
			expectedSize: 30,
		},
		{
			name:         "negative page size defaults to 30",
			params:       enterDirectoryArgs{Page: 2, PageSize: -10},
			expectedPage: 2,
			expectedSize: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page, pageSize := getDefaultPageParams(&tt.params)
			if page != tt.expectedPage {
				t.Errorf("getDefaultPageParams() page = %v, expected %v", page, tt.expectedPage)
			}
			if pageSize != tt.expectedSize {
				t.Errorf("getDefaultPageParams() pageSize = %v, expected %v", pageSize, tt.expectedSize)
			}
		})
	}
}

func TestEnterDirectory_ValidDirectory(t *testing.T) {
	// Clear scan results before test
	scanResults = []scanResult{}

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, validDirectoryHTML)
	}))
	defer server.Close()

	// Create test parameters
	params := &mcp.CallToolParamsFor[enterDirectoryArgs]{
		Arguments: enterDirectoryArgs{
			URL:      server.URL,
			Page:     1,
			PageSize: 30,
		},
	}

	// Call function
	result, err := enterDirectory(context.Background(), nil, params)

	// Assertions
	if err != nil {
		t.Fatalf("enterDirectory() error = %v, expected nil", err)
	}

	if result == nil {
		t.Fatal("enterDirectory() returned nil result")
	}

	if len(result.Content) == 0 {
		t.Fatal("enterDirectory() returned empty content")
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatal("enterDirectory() did not return TextContent")
	}

	// Check if the result contains expected directory listing
	if !strings.Contains(textContent.Text, "Directory contents") {
		t.Error("Result should contain 'Directory contents'")
	}

	if !strings.Contains(textContent.Text, "FILE: file1.txt") {
		t.Error("Result should contain 'FILE: file1.txt'")
	}

	if !strings.Contains(textContent.Text, "DIRECTORY: subdirectory") {
		t.Error("Result should contain 'DIRECTORY: subdirectory'")
	}

	// Check scan results were populated
	if len(scanResults) == 0 {
		t.Error("scanResults should be populated")
	}

	// Verify scan results contain expected entries
	foundFile := false
	foundDir := false
	for _, sr := range scanResults {
		if sr.Type == "FILE" && strings.Contains(sr.URL, "file1.txt") {
			foundFile = true
			if sr.FileExtension != ".txt" {
				t.Errorf("Expected file extension .txt, got %s", sr.FileExtension)
			}
		}
		if sr.Type == "DIRECTORY" && strings.Contains(sr.URL, "subdirectory") {
			foundDir = true
		}
	}

	if !foundFile {
		t.Error("Expected to find file entry in scan results")
	}
	if !foundDir {
		t.Error("Expected to find directory entry in scan results")
	}
}

func TestEnterDirectory_NonDirectory(t *testing.T) {
	// Create mock server that returns non-directory content
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, nonDirectoryHTML)
	}))
	defer server.Close()

	// Create test parameters
	params := &mcp.CallToolParamsFor[enterDirectoryArgs]{
		Arguments: enterDirectoryArgs{
			URL: server.URL,
		},
	}

	// Call function
	result, err := enterDirectory(context.Background(), nil, params)

	// Assertions
	if err != nil {
		t.Fatalf("enterDirectory() error = %v, expected nil", err)
	}

	if result == nil {
		t.Fatal("enterDirectory() returned nil result")
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatal("enterDirectory() did not return TextContent")
	}

	expectedMessage := "This page is not a directory listing (Index of not found)"
	if textContent.Text != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, textContent.Text)
	}
}

func TestEnterDirectory_InvalidURL(t *testing.T) {
	// Create test parameters with invalid URL
	params := &mcp.CallToolParamsFor[enterDirectoryArgs]{
		Arguments: enterDirectoryArgs{
			URL: "invalid-url",
		},
	}

	// Call function
	result, err := enterDirectory(context.Background(), nil, params)

	// Assertions
	if err == nil {
		t.Fatal("enterDirectory() should return error for invalid URL")
	}

	if result != nil {
		t.Error("enterDirectory() should return nil result for invalid URL")
	}

	if !strings.Contains(err.Error(), "request failed") {
		t.Errorf("Expected error to contain 'request failed', got: %v", err)
	}
}

func TestEnterDirectory_HTTPError(t *testing.T) {
	// Create mock server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Create test parameters
	params := &mcp.CallToolParamsFor[enterDirectoryArgs]{
		Arguments: enterDirectoryArgs{
			URL: server.URL,
		},
	}

	// Call function
	result, err := enterDirectory(context.Background(), nil, params)

	// For this test, the function should still succeed even with 404
	// because it reads the response body regardless of status code
	if err != nil {
		t.Fatalf("enterDirectory() error = %v, expected nil", err)
	}

	if result == nil {
		t.Fatal("enterDirectory() returned nil result")
	}
}

func TestEnterDirectory_Pagination(t *testing.T) {
	// Clear scan results before test
	scanResults = []scanResult{}

	// Create mock server with large directory (50 files)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, generateLargeDirectoryHTML(50))
	}))
	defer server.Close()

	// Test first page with page size 10
	params := &mcp.CallToolParamsFor[enterDirectoryArgs]{
		Arguments: enterDirectoryArgs{
			URL:      server.URL,
			Page:     1,
			PageSize: 10,
		},
	}

	result, err := enterDirectory(context.Background(), nil, params)

	if err != nil {
		t.Fatalf("enterDirectory() error = %v, expected nil", err)
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatal("enterDirectory() did not return TextContent")
	}

	// Check pagination info
	if !strings.Contains(textContent.Text, "Page 1 of 5") {
		t.Error("Expected 'Page 1 of 5' in pagination info")
	}

	if !strings.Contains(textContent.Text, "Total entries: 50") {
		t.Error("Expected 'Total entries: 50' in pagination info")
	}

	if !strings.Contains(textContent.Text, "Entries per page: 10") {
		t.Error("Expected 'Entries per page: 10' in pagination info")
	}

	// Check that only first 10 files are shown
	if !strings.Contains(textContent.Text, "file1.txt") {
		t.Error("Expected file1.txt in first page")
	}

	if !strings.Contains(textContent.Text, "file10.txt") {
		t.Error("Expected file10.txt in first page")
	}

	if strings.Contains(textContent.Text, "file11.txt") {
		t.Error("Did not expect file11.txt in first page")
	}

	// Test page out of bounds
	params.Arguments.Page = 10 // Only 5 pages exist
	result, err = enterDirectory(context.Background(), nil, params)

	if err != nil {
		t.Fatalf("enterDirectory() error = %v, expected nil", err)
	}

	textContent, ok = result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatal("enterDirectory() did not return TextContent")
	}

	// Should show last page (page 5)
	if !strings.Contains(textContent.Text, "Page 5 of 5") {
		t.Error("Expected 'Page 5 of 5' when requesting page out of bounds")
	}
}

func TestEnterDirectory_EmptyPage(t *testing.T) {
	// Clear scan results before test
	scanResults = []scanResult{}

	// Create mock server with small directory (only 5 files)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, generateLargeDirectoryHTML(5))
	}))
	defer server.Close()

	// Request page 3 when only 1 page exists (with page size 10)
	params := &mcp.CallToolParamsFor[enterDirectoryArgs]{
		Arguments: enterDirectoryArgs{
			URL:      server.URL,
			Page:     3,
			PageSize: 10,
		},
	}

	result, err := enterDirectory(context.Background(), nil, params)

	if err != nil {
		t.Fatalf("enterDirectory() error = %v, expected nil", err)
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatal("enterDirectory() did not return TextContent")
	}

	// Should show page 1 (last available page)
	if !strings.Contains(textContent.Text, "Page 1 of 1") {
		t.Error("Expected 'Page 1 of 1' when requesting page beyond available pages")
	}
}

func TestEnterDirectory_DefaultPageParams(t *testing.T) {
	// Clear scan results before test
	scanResults = []scanResult{}

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, validDirectoryHTML)
	}))
	defer server.Close()

	// Test with no page parameters (should use defaults)
	params := &mcp.CallToolParamsFor[enterDirectoryArgs]{
		Arguments: enterDirectoryArgs{
			URL: server.URL,
		},
	}

	result, err := enterDirectory(context.Background(), nil, params)

	if err != nil {
		t.Fatalf("enterDirectory() error = %v, expected nil", err)
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatal("enterDirectory() did not return TextContent")
	}

	// Should use default page size of 30
	if !strings.Contains(textContent.Text, "Entries per page: 30") {
		t.Error("Expected default page size of 30")
	}

	if !strings.Contains(textContent.Text, "Page 1 of 1") {
		t.Error("Expected 'Page 1 of 1' with default parameters")
	}
}
