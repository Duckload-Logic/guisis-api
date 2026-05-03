package files

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/hash"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/ocr"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/storage"
)

type Service struct {
	repo      *Repository
	storage   storage.FileStorage
	ocrClient *ocr.OCRClient
}

func NewService(
	repo *Repository,
	storage storage.FileStorage,
	ocrClient *ocr.OCRClient,
) *Service {
	return &Service{
		repo:      repo,
		storage:   storage,
		ocrClient: ocrClient,
	}
}

func (s *Service) GetFileByID(ctx context.Context, id string) (*File, error) {
	return s.repo.GetFileByID(ctx, id)
}

func (s *Service) GetOCRResult(
	ctx context.Context,
	fileID string,
) (*OCRResult, error) {
	return s.repo.GetOCRResult(ctx, fileID)
}

func (s *Service) UploadFile(
	ctx context.Context,
	fileHeader *multipart.FileHeader,
	prefix string,
) (File, error) {
	files, err := s.UploadFiles(
		ctx,
		[]*multipart.FileHeader{fileHeader},
		prefix,
	)
	if err != nil {
		return File{}, err
	}

	return files[0], nil
}

const (
	MaxFileSize = 5 * 1024 * 1024 // 5MB limit
)

var AllowedMimeTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"application/pdf": true,
}

func (s *Service) UploadFiles(
	ctx context.Context,
	filesHeaders []*multipart.FileHeader,
	prefix string,
) ([]File, error) {
	var filesToCreate []File
	var ocrResults []OCRResult
	var uploadedBlobPaths []string

	for _, fh := range filesHeaders {
		// Size Validation
		if fh.Size > MaxFileSize {
			return nil, fmt.Errorf(
				"file %s exceeds maximum size of 5MB",
				fh.Filename,
			)
		}

		ext := strings.ToLower(filepath.Ext(fh.Filename))
		fileHash := hash.GetSHA256Hash(
			fmt.Sprintf("%s%d", fh.Filename, time.Now().UnixNano()),
			16,
		)
		uniqueFileName := fileHash + ext
		folderHash := hash.GetSHA256Hash(
			fmt.Sprintf("%s%d", prefix, time.Now().UnixNano()),
			8,
		)
		blobPath := fmt.Sprintf("%s/%s/%s", prefix, folderHash, uniqueFileName)

		src, err := fh.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()

		data, err := io.ReadAll(src)
		if err != nil {
			return nil, err
		}

		// MIME Type Validation
		contentType := http.DetectContentType(data)
		if !AllowedMimeTypes[contentType] {
			return nil, fmt.Errorf(
				"file %s has unauthorized type: %s",
				fh.Filename,
				contentType,
			)
		}

		reader := bytes.NewReader(data)

		if err := s.storage.Upload(
			ctx, blobPath, reader, contentType,
		); err != nil {
			return nil, fmt.Errorf("failed to upload %s: %w", fh.Filename, err)
		}
		uploadedBlobPaths = append(uploadedBlobPaths, blobPath)

		fileID := uuid.New().String()
		fileRecord := File{
			ID:       fileID,
			FileName: fh.Filename,
			FileURL:  "/uploads/" + blobPath,
			FileType: contentType,
			FileSize: fh.Size,
			MimeType: contentType,
		}
		filesToCreate = append(filesToCreate, fileRecord)

		switch prefix {
		case "cors":
			corResp, err := s.ocrClient.ProcessCOR(
				ctx,
				fh.Filename,
				bytes.NewReader(data),
			)
			if err != nil {
				return nil, fmt.Errorf("this file does not appear to be a valid COR: %w", err)
			}
			if corResp == nil {
				return nil, fmt.Errorf("AI service returned empty response for COR")
			}

			marshaled, _ := json.Marshal(corResp)
			ocrResults = append(ocrResults, OCRResult{
				FileID:         fileID,
				StructuredData: string(marshaled),
				EngineV:        "paddleocr-v4-cor",
				CreatedAt:      time.Now(),
			})
		case "slips":
			// ocrResp, err := s.ocrClient.ProcessDocument(
			// 	ctx,
			// 	fh.Filename,
			// 	bytes.NewReader(data),
			// )
			// if err != nil {
			// 	fmt.Printf("[FileService] {OCR Generic Error}: %v\n", err)
			// } else if ocrResp != nil {
			// 	for i := range ocrResp.Pages {
			// 		ocrResp.Pages[i].Words = nil
			// 	}

			// 	marshaled, _ := json.Marshal(ocrResp)
			// 	ocrResults = append(ocrResults, OCRResult{
			// 		FileID:         fileID,
			// 		RawText:        ocrResp.FullText,
			// 		StructuredData: string(marshaled),
			// 		EngineV:        "paddleocr-v4-generic",
			// 		CreatedAt:      time.Now(),
			// 	})
			// }
		}
	}

	// Transaction with Cleanup (Reliability)
	err := s.repo.WithTransaction(ctx, func(tx datastore.DB) error {
		if _, err := s.repo.CreateBulk(ctx, tx, filesToCreate); err != nil {
			return err
		}

		for _, res := range ocrResults {
			if err := s.repo.SaveOCRResult(ctx, tx, res); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		// Cleanup uploaded files if DB transaction fails
		for _, path := range uploadedBlobPaths {
			_ = s.storage.Delete(ctx, path)
		}
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}

	return filesToCreate, nil
}

func (s *Service) DeleteFile(ctx context.Context, id string) error {
	file, err := s.repo.GetFileByID(ctx, id)
	if err != nil {
		return err
	}
	if file == nil {
		return nil
	}

	err = s.repo.WithTransaction(ctx, func(tx datastore.DB) error {
		blobPath := strings.TrimPrefix(file.FileURL, "/uploads/")

		// Security: Path Traversal Protection
		if strings.Contains(blobPath, "..") {
			return fmt.Errorf("security: invalid file path detected")
		}

		_ = s.storage.Delete(ctx, blobPath)

		return s.repo.Delete(ctx, tx, id)
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}
