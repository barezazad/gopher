package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"gopher/internal/core"
	"io/ioutil"
)

// encryption
func Encrypt(e *core.Engine, origData []byte) (string, error) {

	keyFile, err := ioutil.ReadFile(e.Environments.RSA.CipherPublicKey)
	if err != nil {
		return "", errors.New("error in reading public key")
	}

	var publicKey = keyFile

	//Decryption of public key in pem format
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return "", errors.New("public key error")
	}

	// Parsing public key
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// Type Asserts
	pub := pubInterface.(*rsa.PublicKey)

	//encryption
	var data []byte
	if data, err = rsa.EncryptPKCS1v15(rand.Reader, pub, origData); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// decrypt
func Decrypt(e *core.Engine, ciphertext string) ([]byte, error) {

	keyFile, err := ioutil.ReadFile(e.Environments.RSA.CipherPrivateKey)
	if err != nil {
		return nil, errors.New("error in reading private key")
	}

	var privateKey = keyFile

	cipText, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, errors.New("error in decoding string")
	}

	//decrypt
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}

	//Analysis of private key in PKCS1 format
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// decrypt
	return rsa.DecryptPKCS1v15(rand.Reader, priv, cipText)
}
