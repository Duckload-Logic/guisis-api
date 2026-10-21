package storage_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockStorage implements storage.FileStorage in memory for testing
type mockStorage struct {
	files map[string][]byte
}

func newMockStorage() *mockStorage {
	return &mockStorage{files: make(map[string][]byte)}
}

func (m *mockStorage) Upload(ctx context.Context, path string, reader io.ReadSeeker, contentType string) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	m.files[path] = data
	return nil
}

func (m *mockStorage) Download(ctx context.Context, path string, writer io.Writer) error {
	data, ok := m.files[path]
	if !ok {
		return os.ErrNotExist
	}
	_, err := writer.Write(data)
	return err
}

func (m *mockStorage) Delete(ctx context.Context, path string) error {
	delete(m.files, path)
	return nil
}

func generateTestKey(t *testing.T) []byte {
	t.Helper()
	key := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key)
	require.NoError(t, err)
	return key
}

func TestParseKey(t *testing.T) {
	rawKey := generateTestKey(t)

	// Hex format (64 chars)
	hexStr := hex.EncodeToString(rawKey)
	parsedHex, err := storage.ParseKey(hexStr)
	require.NoError(t, err)
	assert.Equal(t, rawKey, parsedHex)

	// Base64 format
	b64Str := base64.StdEncoding.EncodeToString(rawKey)
	parsedB64, err := storage.ParseKey(b64Str)
	require.NoError(t, err)
	assert.Equal(t, rawKey, parsedB64)

	// Raw 32 bytes string
	rawStr := string(rawKey)
	parsedRaw, err := storage.ParseKey(rawStr)
	require.NoError(t, err)
	assert.Equal(t, rawKey, parsedRaw)

	// Invalid length
	_, err = storage.ParseKey("short-key")
	require.Error(t, err)

	// Empty string
	_, err = storage.ParseKey("")
	require.Error(t, err)
}

func TestEncryptDecryptData(t *testing.T) {
	key := generateTestKey(t)
	original := []byte("Hello, this is sensitive student COR data with grades and personal information.")

	encrypted, err := storage.EncryptData(key, original)
	require.NoError(t, err)

	// Must start with magic header
	assert.True(t, storage.IsEncrypted(encrypted))
	assert.NotEqual(t, original, encrypted)

	decrypted, err := storage.DecryptData(key, encrypted)
	require.NoError(t, err)
	assert.Equal(t, original, decrypted)
}

func TestDecryptTamperedData(t *testing.T) {
	key := generateTestKey(t)
	original := []byte("Important document to protect")

	encrypted, err := storage.EncryptData(key, original)
	require.NoError(t, err)

	// Tamper with ciphertext byte (after header and nonce)
	tampered := make([]byte, len(encrypted))
	copy(tampered, encrypted)
	tampered[len(tampered)-1] ^= 0xFF

	_, err = storage.DecryptData(key, tampered)
	require.Error(t, err)
	assert.ErrorIs(t, err, storage.ErrCorruptedData)
}

func TestEncryptedStorage_UploadAndDownload(t *testing.T) {
	key := generateTestKey(t)
	mock := newMockStorage()

	encStorage, err := storage.NewEncryptedStorage(mock, key)
	require.NoError(t, err)

	ctx := context.Background()
	testPath := "development/cors/test_file.pdf"
	content := []byte("%PDF-1.4 simulated binary PDF content with sensitive student information")

	// 1. Upload
	err = encStorage.Upload(ctx, testPath, bytes.NewReader(content), "application/pdf")
	require.NoError(t, err)

	// Verify underlying storage received encrypted bytes
	storedRaw, ok := mock.files[testPath]
	require.True(t, ok)
	assert.True(t, storage.IsEncrypted(storedRaw))
	assert.NotEqual(t, content, storedRaw)

	// 2. Download through EncryptedStorage
	var out bytes.Buffer
	err = encStorage.Download(ctx, testPath, &out)
	require.NoError(t, err)
	assert.Equal(t, content, out.Bytes())

	// 3. Delete
	err = encStorage.Delete(ctx, testPath)
	require.NoError(t, err)
	assert.Empty(t, mock.files)
}

func TestEncryptedStorage_LegacyFallback(t *testing.T) {
	key := generateTestKey(t)
	mock := newMockStorage()

	encStorage, err := storage.NewEncryptedStorage(mock, key)
	require.NoError(t, err)

	ctx := context.Background()
	testPath := "development/cors/legacy_unencrypted.pdf"
	legacyContent := []byte("%PDF-1.4 legacy unencrypted student profile picture")

	// Store directly into underlying storage WITHOUT encryption (simulating pre-existing file)
	mock.files[testPath] = legacyContent
	assert.False(t, storage.IsEncrypted(legacyContent))

	// Download through EncryptedStorage should fallback cleanly and return original plaintext
	var out bytes.Buffer
	err = encStorage.Download(ctx, testPath, &out)
	require.NoError(t, err)
	assert.Equal(t, legacyContent, out.Bytes())
}

func TestEncryptedStorage_DiskStorageIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	diskStorage := storage.NewDiskStorage(tmpDir)

	key := generateTestKey(t)
	encStorage, err := storage.NewEncryptedStorage(diskStorage, key)
	require.NoError(t, err)

	ctx := context.Background()
	relPath := "slips/subfolder/attachment.png"
	payload := []byte("image PNG binary data simulation")

	// Upload to disk
	err = encStorage.Upload(ctx, relPath, bytes.NewReader(payload), "image/png")
	require.NoError(t, err)

	// Check file on disk is encrypted
	diskFile := filepath.Join(tmpDir, filepath.FromSlash(relPath))
	diskBytes, err := os.ReadFile(diskFile)
	require.NoError(t, err)
	assert.True(t, storage.IsEncrypted(diskBytes))
	assert.NotEqual(t, payload, diskBytes)

	// Download from disk via encStorage
	var out bytes.Buffer
	err = encStorage.Download(ctx, relPath, &out)
	require.NoError(t, err)
	assert.Equal(t, payload, out.Bytes())
}

