package types

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"io"
	mathRand "math/rand"
	"time"
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

	chunkEncMsg := chunkMessage(encMsg, size)
	chunkKey := chunkMessage(key[:], size)

	return mergeKeyMsgChunks(chunkEncMsg, chunkKey)
}

func DecryptFromChunks(data [][]byte) ([]byte, error) {
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

func chunkMessage(msg []byte, size int) [][]byte {
	chunkSize := (len(msg) + size - 1) / size
	var divided [][]byte

	for i := 0; i < len(msg); i += chunkSize {
		end := i + chunkSize
		if end > len(msg) {
			end = len(msg)
		}
		divided = append(divided, msg[i:end])
	}

	return divided
}

func mergeKeyMsgChunks(msg, key [][]byte) ([][]byte, error) {
	if len(msg) != len(key) {
		return nil, errors.New("message and key chunks have different lengths")
	}

	merged := make([][]byte, len(msg))
	for i := 0; i < len(msg); i += 1 {
		//todo: check performance
		merged[i] = make([]byte, 0, len(key[i])+len(msg[i])+1)
		merged[i] = append(merged[i], byte(i))
		merged[i] = append(merged[i], key[i]...)
		merged[i] = append(merged[i], msg[i]...)
	}

	return merged, nil
}

func extractDataFromChunk(chunks [][]byte, keySize KeySize) (message []byte, key []byte) {
	keyPrefixLen := (keySize.Len() + len(chunks) - 1) / len(chunks)
	for i := 0; i < len(chunks); i += 1 {
		if chunks[i][0] != byte(i) {
			panic("unsorted chunks")
		}
		if i == len(chunks)-1 {
			keyPrefixLen = keySize.Len() - i*keyPrefixLen
		}
		chunkWithoutIndex := chunks[i][1:]
		key = append(key, chunkWithoutIndex[:keyPrefixLen]...)
		message = append(message, chunkWithoutIndex[keyPrefixLen:]...)
	}
	return
}

// Temporary not in use yet.
func createChunkLink(chunks [][]byte) [][]byte {
	if len(chunks) < 2 {
		return chunks
	}

	used := genUsedBytesMap(chunks)
	linkByte := genLinkByte(used, true)

	for i := 0; i < len(chunks)-1; i += 1 {
		chunks[i+1] = prependByte(chunks[i+1], linkByte)
		if i < len(chunks)-2 {
			chunks[i] = append(chunks[i], linkByte)
			used[linkByte] = struct{}{}
			linkByte = genLinkByte(used, true)
		}
	}

	return chunks
}

func genUsedBytesMap(chunks [][]byte) map[byte]struct{} {
	firstChunk := chunks[0]
	lastChunk := chunks[len(chunks)-1]

	occupied := make(map[byte]struct{}, len(chunks)+2)
	occupied[firstChunk[0]] = struct{}{}
	occupied[lastChunk[len(lastChunk)-1]] = struct{}{}

	return occupied
}

func genLinkByte(used map[byte]struct{}, seed bool) byte {
	if seed {
		mathRand.Seed(time.Now().UnixNano())
	}

	bytes := make([]byte, 0, len(used)*2)
	var _, _ = mathRand.Read(bytes)
	for _, b := range bytes {
		if _, ok := used[b]; ok {
			continue
		}
		return b
	}
	return genLinkByte(used, false)
}

func prependByte(x []byte, y byte) []byte {
	x = append(x, 0)
	copy(x[1:], x)
	x[0] = y
	return x
}
