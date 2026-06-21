package slips

import (
	"bytes"
	"mime/multipart"
	"testing"
)

func createMockFileHeader(
	t *testing.T,
	filename string,
	content []byte,
) *multipart.FileHeader {
	t.Helper()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}

	_, err = part.Write(content)
	if err != nil {
		t.Fatalf("failed to write mock content: %v", err)
	}
	writer.Close()

	reader := multipart.NewReader(&buf, writer.Boundary())
	form, err := reader.ReadForm(10 << 20)
	if err != nil {
		t.Fatalf("failed to read form: %v", err)
	}

	return form.File["file"][0]
}

func TestService_ValidateFiles(t *testing.T) {
	s := &Service{}

	// Case 1: Valid PDF
	pdfContent := append([]byte("%PDF-1.4\n"), make([]byte, 100)...)
	pdfHeader := createMockFileHeader(t, "test.pdf", pdfContent)
	err := s.validateFiles([]*multipart.FileHeader{pdfHeader})
	if err != nil {
		t.Errorf("expected valid PDF to pass, got error: %v", err)
	}

	// Case 2: Valid PNG
	pngContent := append(
		[]byte("\x89PNG\r\n\x1a\n"),
		make([]byte, 100)...,
	)
	pngHeader := createMockFileHeader(t, "test.png", pngContent)
	err = s.validateFiles([]*multipart.FileHeader{pngHeader})
	if err != nil {
		t.Errorf("expected valid PNG to pass, got error: %v", err)
	}

	// Case 3: Invalid mime type (fake extension)
	fakePdfContent := []byte("plain text content pretending to be pdf")
	fakePdfHeader := createMockFileHeader(t, "fake.pdf", fakePdfContent)
	err = s.validateFiles([]*multipart.FileHeader{fakePdfHeader})
	if err == nil {
		t.Errorf("expected error for fake PDF, got nil")
	}

	// Case 4: Invalid extension
	txtContent := []byte("plain text content")
	txtHeader := createMockFileHeader(t, "test.txt", txtContent)
	err = s.validateFiles([]*multipart.FileHeader{txtHeader})
	if err == nil {
		t.Errorf("expected error for text extension, got nil")
	}

	// Case 5: File size too large
	largeContent := make([]byte, MaxFileSize+1)
	copy(largeContent, []byte("\x89PNG\r\n\x1a\n"))
	largeHeader := createMockFileHeader(t, "large.png", largeContent)
	err = s.validateFiles([]*multipart.FileHeader{largeHeader})
	if err == nil {
		t.Errorf("expected error for too large file, got nil")
	}
}
