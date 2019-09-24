package encryption

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryption(t *testing.T) {
	data := "super-secret-data"
	passphrase := "some-secure-passphrase"

	// Normal string test
	ciphertext, err := Encrypt([]byte(data), "shortpass")
	assert.Equal(t, err.Error(), "Passphrase should be at least 12 characters long.")

	ciphertext, err = Encrypt([]byte(data), passphrase)
	assert.Nil(t, err)

	plaintext, err := Decrypt(ciphertext, passphrase)
	assert.Nil(t, err)
	assert.Equal(t, data, string(plaintext))

	// B64 Test
	ciphertextB64, err := EncryptB64(data, "shortpass")
	assert.Equal(t, err.Error(), "Passphrase should be at least 12 characters long.")

	ciphertextB64, err = EncryptB64(data, passphrase)
	assert.Nil(t, err)

	plaintext, err = DecryptB64(ciphertextB64, passphrase)
	assert.Nil(t, err)
	assert.Equal(t, data, string(plaintext))
}
