package security

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"io"
)

// PrivateKey приватный RSA ключ, используется дл расшифровки данных.
type PrivateKey *rsa.PrivateKey

// NewPrivateKey считывает приватный RSA ключ из файла.
func NewPrivateKey(path string) (PrivateKey, error) {
	rawKey, err := getRawKey(path)
	if err != nil {
		return nil, err
	}

	priv, err := x509.ParsePKCS8PrivateKey(rawKey.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrNotSupportedFormatKey
	}

	return key, nil
}

// Decrypt дешифрует полученное сообщение с использованием алгоритма RSA.
func Decrypt(srcMsg io.Reader, key PrivateKey) (*bytes.Buffer, error) {
	// The message must be no longer than
	// the length of the public modulus.
	bufSize := key.PublicKey.Size()
	msg := new(bytes.Buffer)

	for {
		buf := make([]byte, bufSize)
		n, err := srcMsg.Read(buf)
		if n > 0 {
			plaintext, decErr := rsa.DecryptOAEP(sha256.New(), nil, key, buf, nil)
			if decErr != nil {
				return nil, fmt.Errorf("security.Decrypt rsa.DecryptOAEP: %w", decErr)
			}

			msg.Write(plaintext)
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("security.Decrypt srcMsg.Read: %w", err)
		}
	}

	return msg, nil
}
