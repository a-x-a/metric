package security

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"io"
)

// PublicKey публичный RSA ключ. используется для штфрования данных.
type PublicKey *rsa.PublicKey

// NewPublicKey считывает публичный RSA ключ из файла.
func NewPublicKey(path string) (PublicKey, error) {
	rawKey, err := getRawKey(path)
	if err != nil {
		return nil, err
	}

	pub, err := x509.ParsePKIXPublicKey(rawKey.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, ErrNotSupportedFormatKey
	}

	return key, nil
}

// Encrypt шифрует полученное сообщение с использованием алгоритма RSA.
func Encrypt(srcMsg io.Reader, key PublicKey) (*bytes.Buffer, error) {
	// The message must be no longer than
	// the length of the public modulus
	// minus twice the hash length,
	// minus a further 2.
	bufSize := (*rsa.PublicKey)(key).Size() - 2*sha256.New().Size() - 2
	msg := new(bytes.Buffer)

	for {
		buf := make([]byte, bufSize)
		n, err := srcMsg.Read(buf)
		if n > 0 {
			ciphertext, encErr := rsa.EncryptOAEP(sha256.New(), rand.Reader, key, buf, nil)
			if encErr != nil {
				return nil, fmt.Errorf("security.Encrypt rsa.EncryptOAEP: %w", encErr)
			}

			msg.Write(ciphertext)
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("security.Encrypt srcMsg.Read: %w", err)
		}
	}

	return msg, nil
}
