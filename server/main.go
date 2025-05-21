package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var baseDir string

// FileInfo represents file metadata
type FileInfo struct {
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Extension  interface{} `json:"extension"` // Can be string or null
	Size       interface{} `json:"size"`      // Can be int64 or null
	ModifiedAt string      `json:"modifiedAt"`
	CreatedAt  string      `json:"createdAt"`
}

// DirectoryResponse represents API response for directory listing
type DirectoryResponse struct {
	Files []FileInfo `json:"files"`
}

// ErrorResponse represents API error response
type ErrorResponse struct {
	Error string `json:"error"`
}

func init() {
	_ = os.MkdirAll("./files", 0777)

	var err error
	baseDir, err = filepath.Abs("./files")
	if err != nil {
		log.Fatal("Failed to resolve base directory:", err)
	}
}

func main() {
	// Create router and setup API routes
	http.HandleFunc("/api/files/list", listFilesHandler)
	http.HandleFunc("/api/files", serveFileHandler)

	// Serve static files from the public directory
	http.Handle("/", http.FileServer(http.Dir("public")))

	// Start server
	port := "8088" //TODO: add config

	fmt.Printf("Server running on port %s, serving files from %s\n", port, baseDir)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// resolvePath validates and resolves a requested path to prevent directory traversal
func resolvePath(requestPath string) (string, error) {
	// Replace empty path with "/"
	if requestPath == "" {
		requestPath = "/"
	}

	// Clean the path to remove ".." segments
	cleaned := filepath.Clean(requestPath)
	cleanedPath := strings.Replace(cleaned, "\\", "/", -1)

	// Make sure it doesn't start with ".."
	if strings.HasPrefix(cleanedPath, "../") || cleanedPath == ".." {
		return "", fmt.Errorf("access denied: path is outside the allowed directory")
	}

	// Remove leading slash for joining with base directory
	if strings.HasPrefix(cleanedPath, "/") {
		cleanedPath = cleanedPath[1:]
	}

	// Join with base directory
	fullPath := filepath.Join(baseDir, cleanedPath)

	// Make sure the path is within the base directory
	relPath, err := filepath.Rel(baseDir, fullPath)
	if err != nil || strings.HasPrefix(relPath, "..") {
		return "", fmt.Errorf("access denied: path is outside the allowed directory")
	}

	// Check if path exists
	_, err = os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("path does not exist")
		}
		return "", err
	}

	return fullPath, nil
}

// listFilesHandler handles requests to list directory contents
func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get path from query parameters
	requestPath := r.URL.Query().Get("path")
	fullPath, err := resolvePath(requestPath)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if it's a directory
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to access path")
		return
	}

	if !fileInfo.IsDir() {
		respondWithError(w, http.StatusBadRequest, "Not a directory")
		return
	}

	// Read directory contents
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to read directory")
		return
	}

	// Process directory entries
	files := make([]FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			// Skip entries that can't be accessed
			continue
		}

		fileEntry := FileInfo{
			Name:       entry.Name(),
			ModifiedAt: info.ModTime().Format(time.RFC3339),
		}

		// Get file birth time (created at)
		fileEntry.CreatedAt = info.ModTime().Format(time.RFC3339)

		if entry.IsDir() {
			fileEntry.Type = "directory"
			fileEntry.Extension = nil
			fileEntry.Size = nil
		} else {
			fileEntry.Type = "file"
			ext := strings.TrimPrefix(filepath.Ext(entry.Name()), ".")
			if ext == "" {
				fileEntry.Extension = nil
			} else {
				fileEntry.Extension = strings.ToLower(ext)
			}
			fileEntry.Size = info.Size()
		}

		files = append(files, fileEntry)
	}

	// Sort files (directories first, then files alphabetically)
	sort.Slice(files, func(i, j int) bool {
		if files[i].Type == "directory" && files[j].Type != "directory" {
			return true
		}
		if files[i].Type != "directory" && files[j].Type == "directory" {
			return false
		}
		return files[i].Name < files[j].Name
	})

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(DirectoryResponse{Files: files}); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
	}
}

// serveFileHandler handles requests to download files
func serveFileHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get path from query parameters
	requestPath := r.URL.Query().Get("path")
	if requestPath == "" {
		respondWithError(w, http.StatusBadRequest, "Path parameter is required")
		return
	}

	fullPath, err := resolvePath(requestPath)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check if it's a file
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to access file")
		return
	}

	if fileInfo.IsDir() {
		respondWithError(w, http.StatusBadRequest, "Not a file")
		return
	}

	// Determine content type
	ext := filepath.Ext(fullPath)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Set headers for download
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(fullPath)))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// Serve the file
	http.ServeFile(w, r, fullPath)
}

// respondWithError sends an error response in JSON format
func respondWithError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(ErrorResponse{Error: message}); err != nil {
		log.Printf("Error encoding error response: %v", err)
	}
}
