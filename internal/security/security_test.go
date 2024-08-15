package security

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewKey(t *testing.T) {
	require := require.New(t)

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(err)

	privKeyByte, err := x509.MarshalPKCS8PrivateKey(key)
	require.NoError(err)

	privPemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privKeyByte,
	}

	fpriv, err := os.CreateTemp(os.TempDir(), "*.private.pem")
	require.NoError(err)

	defer os.Remove(fpriv.Name())
	defer fpriv.Close()

	err = pem.Encode(fpriv, privPemBlock)
	require.NoError(err)

	pathToPrivKey := fpriv.Name()

	privKey, err := NewPrivateKey(pathToPrivKey)
	require.NoError(err)
	require.True(key.Equal((*rsa.PrivateKey)(privKey)))

	pubKeyByte, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	require.NoError(err)

	pubPemBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyByte,
	}

	fpub, err := os.CreateTemp(os.TempDir(), "*.public.pem")
	require.NoError(err)

	defer os.Remove(fpub.Name())
	defer fpub.Close()

	err = pem.Encode(fpub, pubPemBlock)
	require.NoError(err)

	pathToPubKey := fpub.Name()

	pubKey, err := NewPublicKey(pathToPubKey)
	require.NoError(err)
	require.True(key.PublicKey.Equal((*rsa.PublicKey)(pubKey)))
}

func Example() {
	secretMessage := []byte("A very secret message")
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	_, err = buf.Write(secretMessage)
	if err != nil {
		log.Fatal(err)
	}

	encryptBuf, err := Encrypt(&buf, &key.PublicKey)
	if err != nil {
		log.Fatal(err)
	}

	decryptBuf, err := Decrypt(encryptBuf, key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(bytes.Trim(decryptBuf.Bytes(), "\x00")))

	// Output:
	// A very secret message
}
