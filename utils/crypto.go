package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
)

var (
	ErrDecodingFailed = errors.New("decoding failed")
)

func Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	return base64.StdEncoding.EncodeToString([]byte(plaintext)), nil
}

func Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("%w: invalid base64", ErrDecodingFailed)
	}
	return string(data), nil
}
