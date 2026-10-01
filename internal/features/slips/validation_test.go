package slips

import (
	"bytes"
	"mime/multipart"
	"strings"
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

func TestValidateDocumentFileCounts(t *testing.T) {
	threeFiles := []*multipart.FileHeader{{}, {}, {}}
	fourFiles := []*multipart.FileHeader{{}, {}, {}, {}}

	if err := validateDocumentFileCounts(threeFiles, threeFiles, threeFiles); err != nil {
		t.Fatalf("expected three files per document type to pass, got: %v", err)
	}

	tests := []struct {
		name       string
		excuse     []*multipart.FileHeader
		parentID   []*multipart.FileHeader
		medical    []*multipart.FileHeader
		wantErrMsg string
	}{
		{
			name:       "excuse letter exceeds limit",
			excuse:     fourFiles,
			wantErrMsg: "Excuse Letter file limit exceeded",
		},
		{
			name:       "parent ID exceeds limit",
			parentID:   fourFiles,
			wantErrMsg: "Parent's ID file limit exceeded",
		},
		{
			name:       "medical certificate exceeds limit",
			medical:    fourFiles,
			wantErrMsg: "Medical Certificate file limit exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDocumentFileCounts(tt.excuse, tt.parentID, tt.medical)
			if err == nil {
				t.Fatal("expected file-count validation error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErrMsg) {
				t.Fatalf("expected error containing %q, got %q", tt.wantErrMsg, err.Error())
			}
		})
	}
}

func TestValidateUpdatedDocumentFileCounts_IncludesKeptFiles(t *testing.T) {
	oldAttachments := []SlipAttachment{
		{FileID: "parent-1", FileName: "parentId-page-1.pdf", AttachmentType: "OTHER"},
		{FileID: "parent-2", FileName: "parentId-page-2.pdf", AttachmentType: "OTHER"},
	}
	newParentFiles := []*multipart.FileHeader{{}, {}}

	err := validateUpdatedDocumentFileCounts(
		oldAttachments,
		[]string{"parent-1", "parent-2"},
		nil,
		newParentFiles,
		nil,
	)
	if err == nil {
		t.Fatal("expected kept + new parent ID files to exceed the limit")
	}
	if !strings.Contains(err.Error(), "Parent's ID file limit exceeded") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateDocumentFileCounts_PreservesFiveMBTotalLimit(t *testing.T) {
	files := []*multipart.FileHeader{
		{Size: 2 * 1024 * 1024},
		{Size: 2 * 1024 * 1024},
		{Size: 2 * 1024 * 1024},
	}

	err := validateDocumentFileCounts(files, nil, nil)
	if err == nil {
		t.Fatal("expected aggregate file size to exceed the existing 5MB limit")
	}
	if !strings.Contains(err.Error(), "total file size limit exceeded") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateUpdatedDocumentFileCounts_IncludesKeptFileSize(t *testing.T) {
	oldAttachments := []SlipAttachment{
		{
			FileID:         "medical-1",
			FileName:       "medicalCert-page-1.pdf",
			AttachmentType: "OTHER",
			FileSize:       4 * 1024 * 1024,
		},
	}
	newMedicalFiles := []*multipart.FileHeader{{Size: 2 * 1024 * 1024}}

	err := validateUpdatedDocumentFileCounts(
		oldAttachments,
		[]string{"medical-1"},
		nil,
		nil,
		newMedicalFiles,
	)
	if err == nil {
		t.Fatal("expected kept + new medical files to exceed the 5MB limit")
	}
	if !strings.Contains(err.Error(), "total file size limit exceeded") {
		t.Fatalf("unexpected error: %v", err)
	}
}
