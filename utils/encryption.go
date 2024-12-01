package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

var key = []byte("0^$D1A=LL-oRC4E}bld7Kd?]s7sFJ,HZ") // Ensure this is 32 bytes (for AES-256)

// Encrypt encrypts data using AES-CFB mode.
func Encrypt(data []byte) ([]byte, error) {
	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create ciphertext slice with room for IV and data
	ciphertext := make([]byte, aes.BlockSize+len(data))

	// Generate a random IV and store it in the first block of the ciphertext
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// Create a stream encrypter and encrypt the data
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

// Decrypt decrypts data using AES-CFB mode.
func Decrypt(data []byte) ([]byte, error) {
	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Ensure the ciphertext length is valid
	if len(data) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	// Extract the IV from the first block of the ciphertext
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	// Create a stream decrypter and decrypt the data
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)

	return data, nil
}
