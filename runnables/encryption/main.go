package main

import (
	"fmt"

	b64 "encoding/base64"

	enc "github.com/scrambledeggs/booky-go-helpers/encryption"
)

// Test run for generating encrypted keys
// TODO: Transform to runnable instance that can accept parameters
func main() {
	value := "super-secret-key"
	passphrase := "some-secure-passphrase"

	fmt.Println("Starting the application...")
	fmt.Printf("Encrypting `%s` with passphrase `%s`", value, passphrase)

	ciphertext, _ := enc.Encrypt([]byte(value), passphrase)
	fmt.Printf("\nEncrypted byte: %x\n", ciphertext)
	fmt.Printf("Encrypted string: %s\n", string(ciphertext))

	ciphertextB64 := b64.StdEncoding.EncodeToString(ciphertext)
	fmt.Printf("Encrypted b64: %s\n", ciphertextB64)

	plaintext, _ := enc.Decrypt(ciphertext, passphrase)
	fmt.Printf("Decrypted using byte value: %s\n", plaintext)

	fmt.Println("B64 Usages...")
	ciphertextB64, _ = enc.EncryptB64(value, passphrase)
	fmt.Printf("Encrypted b64: %s\n", ciphertextB64)

	plaintext, _ = enc.DecryptB64(ciphertextB64, passphrase)
	fmt.Printf("Decrypted using b64 value: %s\n", plaintext)
}
