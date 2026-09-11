package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// DiskStorage implements FileStorage using the local filesystem.
type DiskStorage struct {
	baseDir string
}

func NewDiskStorage(baseDir string) *DiskStorage {
	return &DiskStorage{baseDir: baseDir}
}

func (d *DiskStorage) resolvePath(path string) (string, error) {
	cleanBase := filepath.Clean(d.baseDir)
	cleanPath := filepath.Clean(
		filepath.Join(cleanBase, filepath.FromSlash(path)),
	)

	rel, err := filepath.Rel(cleanBase, cleanPath)
	if err != nil || strings.HasPrefix(rel, "..") || rel == ".." {
		log.Printf(
			"[SECURITY_ALERT] Path traversal breach attempt detected! "+
				"BaseDir: %s, RequestedPath: %s",
			cleanBase,
			path,
		)
		return "", fmt.Errorf(
			"security: path traversal security breach attempt detected",
		)
	}

	return cleanPath, nil
}

func (d *DiskStorage) Upload(
	ctx context.Context,
	path string,
	reader io.ReadSeeker,
	contentType string,
) error {
	fullPath, err := d.resolvePath(path)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0o750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (d *DiskStorage) Download(
	ctx context.Context,
	path string,
	writer io.Writer,
) error {
	fullPath, err := d.resolvePath(path)
	if err != nil {
		return err
	}

	f, err := os.Open(fullPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(writer, f); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	return nil
}

func (d *DiskStorage) Delete(ctx context.Context, path string) error {
	fullPath, err := d.resolvePath(path)
	if err != nil {
		return err
	}

	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
