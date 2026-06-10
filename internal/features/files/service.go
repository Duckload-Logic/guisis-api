package files

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/hash"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/ocr"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/storage"
)

type Service struct {
	repo      *Repository
	storage   storage.FileStorage
	ocrClient *ocr.OCRClient
	logger    audit.Logger
	cfg       *config.Config
}

func NewService(
	repo *Repository,
	storage storage.FileStorage,
	ocrClient *ocr.OCRClient,
	cfg *config.Config,
) *Service {
	return &Service{
		repo:      repo,
		storage:   storage,
		ocrClient: ocrClient,
		cfg:       cfg,
	}
}

func (s *Service) SetLogger(logger audit.Logger) {
	s.logger = logger
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

		var envPrefix string
		if s.cfg.IsStaging {
			envPrefix = "staging/"
		} else if s.cfg.IsProduction {
			envPrefix = "production/"
		} else {
			envPrefix = "development/"
		}

		blobPath := fmt.Sprintf(
			"%s%s/%s/%s",
			envPrefix,
			prefix,
			folderHash,
			uniqueFileName,
		)

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
				if s.logger != nil {
					id, ip, ua, email, _, trace := audit.ExtractMeta(ctx)
					logLevel := audit.LevelError
					logAction := audit.ActionOCRProcessingFailed
					logMsg := fmt.Sprintf(
						"OCR COR processing failed for %s: %v",
						fh.Filename,
						err,
					)

					var httpErr *ocr.HTTPError
					if errors.As(err, &httpErr) &&
						httpErr.StatusCode >= 400 &&
						httpErr.StatusCode < 500 {
						logLevel = audit.LevelWarning
						logAction = audit.ActionOCRValidationFailed
						logMsg = fmt.Sprintf(
							"COR validation failed for %s: %v",
							fh.Filename,
							err,
						)
					} else if strings.Contains(err.Error(), "status: 400") ||
						strings.Contains(err.Error(), "status: 422") {
						logLevel = audit.LevelWarning
						logAction = audit.ActionOCRValidationFailed
						logMsg = fmt.Sprintf(
							"COR validation failed for %s: %v",
							fh.Filename,
							err,
						)
					}

					s.logger.Record(ctx, nil, audit.LogEntry{
						Level:     logLevel,
						Category:  audit.CategorySystem,
						Action:    logAction,
						Message:   logMsg,
						UserID:    structs.StringToNullableString(id),
						UserEmail: structs.StringToNullableString(email),
						IPAddress: structs.StringToNullableString(ip),
						UserAgent: structs.StringToNullableString(ua),
						TraceID:   structs.StringToNullableString(trace),
					})
				}
				return nil, fmt.Errorf(
					"this file does not appear to be a valid COR",
				)
			}
			if corResp == nil {
				if s.logger != nil {
					id, ip, ua, email, _, trace := audit.ExtractMeta(ctx)
					s.logger.Record(ctx, nil, audit.LogEntry{
						Level:    audit.LevelError,
						Category: audit.CategorySystem,
						Action:   audit.ActionOCRProcessingFailed,
						Message: fmt.Sprintf(
							"AI service returned empty response for COR: %s",
							fh.Filename,
						),
						UserID:    structs.StringToNullableString(id),
						UserEmail: structs.StringToNullableString(email),
						IPAddress: structs.StringToNullableString(ip),
						UserAgent: structs.StringToNullableString(ua),
						TraceID:   structs.StringToNullableString(trace),
					})
				}
				return nil, fmt.Errorf(
					"AI service returned empty response for COR",
				)
			}

			if s.logger != nil {
				id, ip, ua, email, _, trace := audit.ExtractMeta(ctx)
				s.logger.Record(ctx, nil, audit.LogEntry{
					Level:    audit.LevelInfo,
					Category: audit.CategorySystem,
					Action:   audit.ActionOCRProcessingSuccess,
					Message: fmt.Sprintf(
						"Successfully processed COR document using OCR: %s",
						fh.Filename,
					),
					UserID:    structs.StringToNullableString(id),
					UserEmail: structs.StringToNullableString(email),
					IPAddress: structs.StringToNullableString(ip),
					UserAgent: structs.StringToNullableString(ua),
					TraceID:   structs.StringToNullableString(trace),
				})
			}

			marshaled, _ := json.Marshal(corResp)
			ocrResults = append(ocrResults, OCRResult{
				FileID:         fileID,
				StructuredData: string(marshaled),
				EngineV:        "paddleocr-v4-cor",
				CreatedAt:      time.Now(),
			})
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

func (s *Service) DownloadFile(
	ctx context.Context,
	fileURL string,
	writer io.Writer,
) (string, error) {
	file, err := s.repo.GetFileByURL(ctx, fileURL)
	if err != nil {
		return "", fmt.Errorf("[FileService] {DownloadFile Meta}: %w", err)
	}

	blobPath := strings.TrimPrefix(file.FileURL, "/uploads/")
	if err := s.storage.Download(ctx, blobPath, writer); err != nil {
		return "", fmt.Errorf("[FileService] {DownloadFile Storage}: %w", err)
	}

	return file.MimeType, nil
}
