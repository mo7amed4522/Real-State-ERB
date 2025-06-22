package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

type EncryptionService struct {
	key []byte
}

func NewEncryptionService(secretKey string) *EncryptionService {
	// Generate a deterministic 32-byte key from the secret
	hash := sha256.Sum256([]byte(secretKey))
	return &EncryptionService{
		key: hash[:],
	}
}

// EncryptText encrypts a string and returns base64 encoded result
func (e *EncryptionService) EncryptText(text string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptText decrypts a base64 encoded encrypted string
func (e *EncryptionService) DecryptText(encryptedText string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// EncryptFile encrypts a file and saves it to a new location
func (e *EncryptionService) EncryptFile(sourcePath, destPath string) (string, error) {
	// Read source file
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return "", err
	}
	defer sourceFile.Close()

	// Create destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer destFile.Close()

	// Generate a random encryption key for this file
	fileKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, fileKey); err != nil {
		return "", err
	}

	// Encrypt the file key with our master key
	encryptedFileKey, err := e.EncryptText(hex.EncodeToString(fileKey))
	if err != nil {
		return "", err
	}

	// Create cipher
	block, err := aes.NewCipher(fileKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Write nonce to destination file
	destFile.Write(nonce)

	// Create buffer for reading
	buffer := make([]byte, 4096)
	for {
		n, err := sourceFile.Read(buffer)
		if err != nil && err != io.EOF {
			return "", err
		}
		if n == 0 {
			break
		}

		// Encrypt the chunk
		ciphertext := gcm.Seal(nil, nonce, buffer[:n], nil)
		destFile.Write(ciphertext)
	}

	return encryptedFileKey, nil
}

// DecryptFile decrypts a file and saves it to a new location
func (e *EncryptionService) DecryptFile(sourcePath, destPath, encryptedFileKey string) error {
	// Decrypt the file key
	fileKeyHex, err := e.DecryptText(encryptedFileKey)
	if err != nil {
		return err
	}

	fileKey, err := hex.DecodeString(fileKeyHex)
	if err != nil {
		return err
	}

	// Open source file
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Create cipher
	block, err := aes.NewCipher(fileKey)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Read nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(sourceFile, nonce); err != nil {
		return err
	}

	// Read and decrypt the rest of the file
	buffer := make([]byte, 4096+gcm.Overhead())
	for {
		n, err := sourceFile.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		// Decrypt the chunk
		plaintext, err := gcm.Open(nil, nonce, buffer[:n], nil)
		if err != nil {
			return err
		}

		destFile.Write(plaintext)
	}

	return nil
}

// GenerateRandomKey generates a random encryption key
func GenerateRandomKey() (string, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
} 