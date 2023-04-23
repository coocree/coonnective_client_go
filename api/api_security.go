package api

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"strings"
)

// SecurityModel ApiSecurity fornece métodos para realizar operações de criptografia e geração de valores aleatórios.
type SecurityModel struct{}

func Security() SecurityModel {
	return SecurityModel{}
}

// RandomLetters randomLetters gera uma string aleatória de letras minúsculas com o tamanho especificado.
func (s *SecurityModel) RandomLetters(length int) string {
	possible := "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	rand.Read(b)
	for i := 0; i < length; i++ {
		b[i] = possible[b[i]%byte(len(possible))]
	}
	return string(b)
}

// RandomNumbers randomNumbers gera uma string aleatória de números com o tamanho especificado.
func (s *SecurityModel) RandomNumbers(length int) string {
	possible := "0123456789"
	b := make([]byte, length)
	rand.Read(b)
	for i := 0; i < length; i++ {
		b[i] = possible[b[i]%byte(len(possible))]
	}
	return string(b)
}

// RandomBytes randomBytes gera uma string aleatória de bytes com o tamanho especificado.
func (s *SecurityModel) RandomBytes(length int) string {
	possible := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	b := make([]byte, length)
	rand.Read(b)
	for i := 0; i < length; i++ {
		b[i] = possible[b[i]%byte(len(possible))]
	}
	return string(b)
}

// Base64URLEncode base64URLEncode codifica uma string em base64URL.
func (s *SecurityModel) Base64URLEncode(str string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(str))
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(encoded, "+", "-"), "/", "_"), "=", "")
}

// EncodeSha256 encodeSha256 codifica uma string usando SHA-256.
func (s *SecurityModel) EncodeSha256(value string) string {
	hash := sha256.Sum256([]byte(value))
	return fmt.Sprintf("%x", hash)
}

// EncodeSha1 encodeSha1 codifica uma string usando SHA-1.
func (s *SecurityModel) EncodeSha1(value string) string {
	hash := sha1.Sum([]byte(value))
	return fmt.Sprintf("%x", hash)
}

// EncodeMD5 encodeMD5 codifica uma string usando MD5.
func (s *SecurityModel) EncodeMD5(value string) string {
	hash := md5.Sum([]byte(value))
	return fmt.Sprintf("%x", hash)
}

// UidSha1 uidSha1 gera um UID no formato "XXXX-XXXX-XXXX-XXXX" com base em uma string de entrada,
// codificando-a usando SHA-1 e convertendo-a em maiúsculas.
func (s *SecurityModel) UidSha1(value string) string {
	hash := sha1.Sum([]byte(value))
	uid := fmt.Sprintf("%x", hash)
	return fmt.Sprintf("%s-%s-%s-%s", uid[0:4], uid[4:8], uid[8:12], uid[12:16])
}

// uid gera um UID no formato "XXXX-XXXX-XXXX-XXXX".
func (s *SecurityModel) Uid() string {
	length := 16
	possible := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	rand.Read(b)
	for i := 0; i < length; i++ {
		b[i] = possible[b[i]%byte(len(possible))]
	}
	return fmt.Sprintf("%s-%s-%s-%s", string(b[0:4]), string(b[4:8]), string(b[8:12]), string(b[12:16]))
}

// encrypt criptografa uma string usando AES com a senha fornecida.
func (s *SecurityModel) EncryptWithBase64(value string, password string) (string, error) {
	salt := make([]byte, 8)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	key := pbkdf2.Key([]byte(password), salt, 10000, 32, sha256.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(value))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(value))

	return base64.URLEncoding.EncodeToString(append(salt, ciphertext...)), nil
}

// decrypt descriptografa uma string criptografada usando AES com a senha fornecida.
func (s *SecurityModel) DecryptWithBase64(ciphertext string, password string) (string, error) {
	ciphertextBytes, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(ciphertextBytes) < aes.BlockSize+8 {
		return "", errors.New("invalid ciphertext")
	}

	salt := ciphertextBytes[:8]
	key := pbkdf2.Key([]byte(password), salt, 10000, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := ciphertextBytes[8 : 8+aes.BlockSize]
	ciphertextBytes = ciphertextBytes[8+aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return string(ciphertextBytes), nil
}

// encrypt criptografa uma string usando AES com a senha fornecida.
func (s *SecurityModel) Encrypt(value string, password string) (string, error) {
	salt := make([]byte, 8)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}
	key := pbkdf2.Key([]byte(password), salt, 10000, 32, sha256.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(value))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(value))
	salt = append(salt, ciphertext...)
	return string(salt), nil
}

// decrypt descriptografa uma string criptografada usando AES com a senha fornecida.
func (s *SecurityModel) Decrypt(value string, password string) (string, error) {
	ciphertext := []byte(value)
	if len(ciphertext) < aes.BlockSize+8 {
		return "", errors.New("invalid ciphertext")
	}
	salt := ciphertext[:8]
	key := pbkdf2.Key([]byte(password), salt, 10000, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := ciphertext[8 : 8+aes.BlockSize]
	ciphertextBytes := ciphertext[8+aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return string(ciphertextBytes), nil
}
