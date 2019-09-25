package main

import (
	"fmt"
	"log"
	"os"

	enc "github.com/scrambledeggs/booky-go-common/encryption"
)

// Test run for generating encrypted keys
func main() {
	var (
		value      string
		passphrase string
	)

	argsLen := len(os.Args)
	if argsLen < 2 {
		fmt.Printf("No Arguments passed. Using sample values instead")
		value = "some-super-secret-data"
		passphrase = "some-super-secret-key"
	} else if argsLen == 3 {
		value = os.Args[1]
		passphrase = os.Args[2]
	} else {
		fmt.Println("Unsupported number of params passed.")
		fmt.Println("Usage: go run main.go <value-to-encrypt> <passphrase>")
		return
	}

	// ENCRYPT
	fmt.Printf("Encrypting `%s` with passphrase `%s`\n", value, passphrase)
	ciphertext, err := enc.Encrypt([]byte(value), passphrase)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Encrypted byte: %x\n", ciphertext)
	fmt.Printf("Encrypted string: %s\n", string(ciphertext))

	ciphertextB64, err := enc.EncryptB64(value, passphrase)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Encrypted b64: %s\n", ciphertextB64)

	// DECRYPT
	plaintext, err := enc.Decrypt(ciphertext, passphrase)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Decrypted using normal Decryption (in:byte out:byte): %s\n", plaintext)

	plaintextB64, err := enc.DecryptB64(ciphertextB64, passphrase)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Decrypted using b64 Decryption (in:b64 out:byte): %s\n", plaintextB64)
	fmt.Printf("Decrypted using b64 Decryption (in:b64 out:string): %s\n", string(plaintextB64))
}
