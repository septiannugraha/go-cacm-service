package uploader

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Uploader handles file uploads to the service
type Uploader struct {
	serviceURL string
	timeout    time.Duration
	client     *http.Client
}

// NewUploader creates a new uploader instance
func NewUploader(serviceURL string, timeout time.Duration) *Uploader {
	return &Uploader{
		serviceURL: serviceURL,
		timeout:    timeout,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// UploadFile uploads a file to the service
func (u *Uploader) UploadFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a buffer to store our request body
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create form file field
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	// Copy file content to form
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	// Close the writer to finalize the form
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", u.serviceURL, &requestBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set the content type with boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	resp, err := u.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// UploadFiles uploads multiple files
func (u *Uploader) UploadFiles(filePaths []string) error {
	for _, filePath := range filePaths {
		fmt.Printf("Uploading %s...\n", filePath)
		if err := u.UploadFile(filePath); err != nil {
			return fmt.Errorf("failed to upload %s: %w", filePath, err)
		}
		fmt.Printf("Successfully uploaded %s\n", filePath)
	}
	return nil
}