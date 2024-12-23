package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"strings"
)

// LoadSymmetricKey loads a 32-byte symmetric key from the specified file path.
// If path is empty, it generates a new random key.
// The key file content can be either raw bytes or base64 encoded.
// It returns the key and any error encountered.
func LoadSymmetricKey(path string) ([]byte, error) {
	if path == "" {
		key := make([]byte, 32)
		if _, err := rand.Read(key); err != nil {
			return nil, fmt.Errorf("failed to generate symmetric key: %w", err)
		}
		return key, nil
	}
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(keyData)))
	if err == nil {
		keyData = decoded
	}
	if len(keyData) != 32 {
		return nil, fmt.Errorf("invalid key length: expected 32 bytes, got %d", len(keyData))
	}
	return keyData, nil
}

// LoadPrivateKey loads an Ed25519 private key from the specified PEM file path.
// If path is empty, it generates a new Ed25519 key pair and returns the private key.
// It returns the private key and any error encountered.
func LoadPrivateKey(path string) (ed25519.PrivateKey, error) {
	if path == "" {
		_, privateKey, err := ed25519.GenerateKey(nil)
		if err != nil {
			return nil, err
		}

		return privateKey, nil
	}
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	signPrivateKey, ok := privateKey.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not an ed25519 private key")
	}
	return signPrivateKey, nil
}

// PemKey converts an Ed25519 private key to PEM-encoded PKCS#8 format.
// It returns the PEM-encoded key string and any error encountered.
func PemKey(privateKey ed25519.PrivateKey) (string, error) {
	pkcs8Key, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", err
	}
	pemKey := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8Key,
	})

	return strings.TrimSpace(string(pemKey)), nil
}

func Base64(key []byte) string {
	return base64.StdEncoding.EncodeToString(key)
}

// Encrypt performs AES-GCM encryption on the input message using the provided symmetric key.
// It returns the base64-encoded ciphertext which includes the nonce.
// If encryption fails, it returns the original message unchanged.
func Encrypt(symmetricKey []byte, msg string) string {
	block, err := aes.NewCipher(symmetricKey)
	if err != nil {
		return msg
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return msg
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return msg
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(msg), nil)
	return base64.StdEncoding.EncodeToString(ciphertext)
}
