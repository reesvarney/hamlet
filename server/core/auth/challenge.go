package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"log"
	"os"
	"strconv"
)

var challenge_size int

func init() {
	// Get the challenge byte size from the ENV
	env_challenge_size, isPresent := os.LookupEnv("AUTH_CHALLENGE_SIZE")
	if isPresent {
		conv_challenge_size, err := strconv.Atoi(env_challenge_size)
		if err != nil {
			log.Fatal("AUTH_CHALLENGE_SIZE environment variable could not be cast to integer")
		}
		challenge_size = conv_challenge_size
		return
	}
	// else just use the default
	challenge_size = 1024

}

func GenerateChallenge(pub *rsa.PublicKey) (challenge []byte, encoded []byte, err error) {
	buf := make([]byte, challenge_size)
	_, err = rand.Read(buf)
	if err != nil {
		return buf, []byte(""), errors.New("challenge_bytes_not_generated")
	}
	// Not sure if there is any need for a label at this point, since the data is not long-lived
	label := []byte("")
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, buf, label)
	if err != nil {
		return buf, ciphertext, errors.New("challenge_bytes_not_encodeable")
	}
	return buf, ciphertext, nil
}
