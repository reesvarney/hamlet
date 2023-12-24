package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func ParsePubKey(pub string) (rsa.PublicKey, []byte, error) {
	publicKeyBytes := []byte(pub)

	//  Decode the PEM block
	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		return rsa.PublicKey{}, block.Bytes, errors.New("invalid_public_key")
	}
	// Parse the public key
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return rsa.PublicKey{}, block.Bytes, err
	}
	// Assert the parsed key type
	publicKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return rsa.PublicKey{}, block.Bytes, errors.New("public_key_not_parseable")
	}
	return *publicKey, block.Bytes, nil

}