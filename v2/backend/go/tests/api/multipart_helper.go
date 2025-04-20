// backend/go/tests/api/multipart_helper.go
package api

import (
	"bytes"
	"io"
	"mime/multipart"
	"strings"
)

// TestMultipartWriter is a wrapper around multipart.Writer for easier testing
type TestMultipartWriter struct {
	writer *multipart.Writer
	buffer *bytes.Buffer
}

// NewTestMultipartWriter creates a new multipart writer for testing
func NewTestMultipartWriter(buffer *bytes.Buffer) *TestMultipartWriter {
	return &TestMultipartWriter{
		writer: multipart.NewWriter(buffer),
		buffer: buffer,
	}
}

// AddFile adds a file to the multipart form
func (t *TestMultipartWriter) AddFile(fieldName, fileName, content string) error {
	part, err := t.writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return err
	}
	
	_, err = io.Copy(part, strings.NewReader(content))
	return err
}

// AddField adds a text field to the multipart form
func (t *TestMultipartWriter) AddField(fieldName, value string) error {
	return t.writer.WriteField(fieldName, value)
}

// Close finalizes the form
func (t *TestMultipartWriter) Close() error {
	return t.writer.Close()
}

// GetContentType returns the content type including boundary
func (t *TestMultipartWriter) GetContentType() string {
	return t.writer.FormDataContentType()
}

// GetBuffer returns the underlying buffer
func (t *TestMultipartWriter) GetBuffer() *bytes.Buffer {
	return t.buffer
}