package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"iac/utils/errors"
	"log/slog"

	"golang.org/x/crypto/scrypt"
)

func generateSalt() []byte {
	salt := make([]byte, 16)
	rand.Read(salt)
	return salt
}

func GenerateSaltB64() string {
	salt := generateSalt()
	return base64.StdEncoding.EncodeToString(salt)
}

func generateKey(plaintextBytes []byte, salt []byte) []byte {
	keyBytes, err := scrypt.Key(plaintextBytes, salt, 32768, 8, 1, 32)
	if err != nil {
		slog.Error("Key generation failed", "error", err)
		return nil
	}
	return keyBytes
}

func GenerateKeyB64(plaintext string, saltBase64 string) (string, error) {
	saltBytes, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		return "", errors.New("failed to decode base64 salt", "error", err)
	}

	plaintextBytes := []byte(plaintext)
	keyBytes := generateKey(plaintextBytes, saltBytes)
	keyBase64 := base64.StdEncoding.EncodeToString(keyBytes)
	return keyBase64, nil
}

func generateIV() []byte {
	iv := make([]byte, 12)
	rand.Read(iv)
	return iv
}

func GenerateIVB64() string {
	iv := generateIV()
	return base64.StdEncoding.EncodeToString(iv)
}

func prepareAES256GCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("failed to create AES cipher", "error", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("failed to create GCM cipher mode", "error", err)
	}

	return gcm, nil
}

func encryptAES256GCM(plainBytes []byte, key []byte, iv []byte) ([]byte, error) {
	gcm, err := prepareAES256GCM(key)
	if err != nil {
		return nil, err
	}

	cipherBytes := gcm.Seal(nil, iv, plainBytes, nil)
	return cipherBytes, nil
}

func decryptAES256GCM(ciphertext []byte, key []byte, iv []byte) ([]byte, error) {
	gcm, err := prepareAES256GCM(key)
	if err != nil {
		return nil, err
	}

	plaintextBytes, err := gcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return nil, errors.New("decryption failed", "error", err)
	}

	return plaintextBytes, nil
}

func encryptAES256GCMB64(plaintext string, keyBase64 string, ivBase64 string) (string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return "", errors.New("failed to decode base64 key", "error", err)
	}

	ivBytes, err := base64.StdEncoding.DecodeString(ivBase64)
	if err != nil {
		return "", errors.New("failed to decode base64 iv", "error", err)
	}

	plainBytes := []byte(plaintext)
	cipherBytes, err := encryptAES256GCM(plainBytes, keyBytes, ivBytes)
	if err != nil {
		slog.Error("Encryption failed", "error", err)
		return "", err
	}

	cipherBase64 := base64.StdEncoding.EncodeToString(cipherBytes)
	return cipherBase64, nil
}

func decryptAES256GCMB64(ciphertextBase64 string, keyBase64 string, ivBase64 string) (string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return "", errors.New("failed to decode base64 key", "error", err)
	}

	ivBytes, err := base64.StdEncoding.DecodeString(ivBase64)
	if err != nil {
		return "", errors.New("failed to decode base64 iv", "error", err)
	}

	cipherBytes, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", errors.New("failed to decode base64 ciphertext", "error", err)
	}

	plainBytes, err := decryptAES256GCM(cipherBytes, keyBytes, ivBytes)
	if err != nil {
		slog.Error("Decryption failed", "error", err)
		return "", err
	}

	return string(plainBytes), nil
}
