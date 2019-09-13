package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	b64 "encoding/base64"
	"io"
	"strings"
)

// Encrypt data using AES
func Encrypt(data []byte, passphrase string) ([]byte, error) {
	gcm := setupCipher(data, passphrase)

	nonce := make([]byte, gcm.NonceSize())
	r := strings.NewReader(passphrase)
	if _, err := io.ReadFull(r, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Encrypt data using AES and encode to B64
func EncryptB64(data string, passphrase string) (string, error) {
	ciphertext, err := Encrypt([]byte(data), passphrase)
	if err != nil {
		return "", err
	}

	return b64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt data using AES
func Decrypt(data []byte, passphrase string) ([]byte, error) {
	gcm := setupCipher(data, passphrase)

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// Decode data to B64 and Decrypt using AES
func DecryptB64(data string, passphrase string) ([]byte, error) {
	dataB64, err := b64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	return Decrypt([]byte(dataB64), passphrase)
}

// AES Cipher initialization given a passphrase
func setupCipher(data []byte, passphrase string) cipher.AEAD {
	// Key hashed to sha256 for consistent 32byte length key for aes cipher
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	return gcm
}

// Used to hash keys in order to use 32byte string regardless of key length
func createHash(key string) string {
	hash := sha256.Sum256([]byte(key))
	return string(hash[:])
}
