package storage

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	// MagicHeader identifies an encrypted file stored at rest (4 bytes magic + 1 byte version).
	MagicHeader = []byte{'G', 'E', 'N', 'C', 0x01}

	ErrInvalidKeySize = errors.New("storage encryption key must be 32 bytes (AES-256)")
	ErrCorruptedData  = errors.New("encrypted file data is corrupted or authentication failed")
)

// ParseKey parses a 32-byte key from either 64-character hex, standard base64, or raw 32-byte string.
func ParseKey(keyStr string) ([]byte, error) {
	keyStr = strings.TrimSpace(keyStr)
	if len(keyStr) == 0 {
		return nil, errors.New("empty encryption key")
	}

	// 1. Try Hex (64 chars -> 32 bytes)
	if len(keyStr) == 64 {
		if keyBytes, err := hex.DecodeString(keyStr); err == nil && len(keyBytes) == 32 {
			return keyBytes, nil
		}
	}

	// 2. Try Base64
	if keyBytes, err := base64.StdEncoding.DecodeString(keyStr); err == nil && len(keyBytes) == 32 {
		return keyBytes, nil
	}

	// 3. Try Raw 32 bytes
	if len(keyStr) == 32 {
		return []byte(keyStr), nil
	}

	return nil, fmt.Errorf("%w (provided key length: %d)", ErrInvalidKeySize, len(keyStr))
}

// IsEncrypted checks whether the given byte slice starts with MagicHeader.
func IsEncrypted(data []byte) bool {
	return len(data) >= len(MagicHeader) && bytes.Equal(data[:len(MagicHeader)], MagicHeader)
}

// EncryptData encrypts plaintext using AES-256-GCM and prepends the magic header and nonce.
func EncryptData(key []byte, plaintext []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, ErrInvalidKeySize
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	result := make([]byte, 0, len(MagicHeader)+len(nonce)+len(ciphertext))
	result = append(result, MagicHeader...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result, nil
}

// DecryptData decrypts data produced by EncryptData.
// If the data does not start with MagicHeader, ErrCorruptedData is returned.
func DecryptData(key []byte, data []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, ErrInvalidKeySize
	}

	if !IsEncrypted(data) {
		return nil, ErrCorruptedData
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	headerLen := len(MagicHeader)

	if len(data) < headerLen+nonceSize+gcm.Overhead() {
		return nil, ErrCorruptedData
	}

	nonce := data[headerLen : headerLen+nonceSize]
	ciphertext := data[headerLen+nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCorruptedData, err)
	}

	return plaintext, nil
}

// EncryptedStorage decorates any FileStorage implementation with transparent AES-256-GCM encryption.
type EncryptedStorage struct {
	underlying FileStorage
	key        []byte
}

// NewEncryptedStorage constructs a new EncryptedStorage wrapper.
func NewEncryptedStorage(underlying FileStorage, key []byte) (*EncryptedStorage, error) {
	if len(key) != 32 {
		return nil, ErrInvalidKeySize
	}
	return &EncryptedStorage{
		underlying: underlying,
		key:        key,
	}, nil
}

// Upload encrypts data from reader in memory with AES-256-GCM and passes ciphertext to underlying storage.
func (e *EncryptedStorage) Upload(
	ctx context.Context,
	path string,
	reader io.ReadSeeker,
	contentType string,
) error {
	plaintext, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("[EncryptedStorage] {Upload Read}: %w", err)
	}

	encryptedData, err := EncryptData(e.key, plaintext)
	if err != nil {
		return fmt.Errorf("[EncryptedStorage] {Upload Encrypt}: %w", err)
	}

	cipherReader := bytes.NewReader(encryptedData)
	return e.underlying.Upload(ctx, path, cipherReader, contentType)
}

// Download retrieves data from underlying storage. If encrypted, it decrypts with AES-256-GCM.
// If not encrypted (legacy file prior to migration), it falls back to streaming plaintext.
func (e *EncryptedStorage) Download(ctx context.Context, path string, writer io.Writer) error {
	var buf bytes.Buffer
	if err := e.underlying.Download(ctx, path, &buf); err != nil {
		return err
	}

	raw := buf.Bytes()
	if IsEncrypted(raw) {
		plaintext, err := DecryptData(e.key, raw)
		if err != nil {
			return fmt.Errorf("[EncryptedStorage] {Download Decrypt}: %w", err)
		}
		if _, err := writer.Write(plaintext); err != nil {
			return fmt.Errorf("[EncryptedStorage] {Download Write}: %w", err)
		}
		return nil
	}

	// Legacy unencrypted file fallback during rollout/migration
	if _, err := writer.Write(raw); err != nil {
		return fmt.Errorf("[EncryptedStorage] {Download Fallback Write}: %w", err)
	}
	return nil
}

// Delete delegates directly to underlying storage.
func (e *EncryptedStorage) Delete(ctx context.Context, path string) error {
	return e.underlying.Delete(ctx, path)
}

