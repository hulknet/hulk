package types

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"io"
)

type KeySize int

func (k KeySize) Len() int {
	switch k {
	case Key256:
		return 32
	case Key512:
		return 64
	default:
		return 0
	}
}

const (
	Key256 KeySize = iota
	Key512
)

func EncryptToParts(message []byte, size int) ([][]byte, error) {
	key, err := generateKey(message, Key256)
	if err != nil {
		return nil, err
	}
	encMsg, err := encryptMessage(message, key[:])
	if err != nil {
		return nil, err
	}

	chunkEncMsg := chunkData(encMsg, size)
	chunkKey := chunkData(key[:], size)

	return mergeChunk(chunkEncMsg, chunkKey)
}

func DecryptFromParts(data [][]byte) ([]byte, error) {
	encMsg, key := extractDataFromChunk(data, Key256)
	return decryptMessage(encMsg, key)
}

func generateKey(data []byte, keySize KeySize) ([]byte, error) {
	switch keySize {
	case Key256:
		key := sha256.Sum256(data)
		return key[:], nil
	case Key512:
		key := sha512.Sum512(data)
		return key[:], nil
	default:
		return nil, errors.New("unknown key size")
	}
}

func encryptMessage(message []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	encrypted := aesGCM.Seal(nonce, nonce, message, nil)
	return encrypted, nil
}

func decryptMessage(encrypted []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]

	message, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return message, err
}

func chunkData(data []byte, size int) [][]byte {
	chunkSize := (len(data) + size - 1) / size
	var divided [][]byte

	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		divided = append(divided, data[i:end])
	}

	return divided
}

func mergeChunk(data, key [][]byte) ([][]byte, error) {
	if len(data) != len(key) {
		return nil, errors.New("data and key chunks have different lengths")
	}

	merged := make([][]byte, len(data))
	for i := 0; i < len(data); i += 1 {
		merged[i] = append(key[i], data[i]...)
	}

	return merged, nil
}

func extractDataFromChunk(chunk [][]byte, keySize KeySize) (message []byte, key []byte) {
	keyPrefixLen := (keySize.Len() + len(chunk) - 1) / len(chunk)
	for i := 0; i < len(chunk); i += 1 {
		if i == len(chunk)-1 {
			keyPrefixLen = keySize.Len() - i*keyPrefixLen
		}
		key = append(key, chunk[i][:keyPrefixLen]...)
		message = append(message, chunk[i][keyPrefixLen-1:]...)
	}
	return
}
