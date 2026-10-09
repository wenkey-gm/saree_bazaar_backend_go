package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"strings"
)

func HashPassword(password string) (string, error) {

	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	shash := pbkdf2.Key([]byte(password), salt, 4096, 32, sha256.New)

	hashedPW := fmt.Sprintf("%s.%s", hex.EncodeToString(shash), hex.EncodeToString(salt))

	return hashedPW, nil
}

func ComparePasswords(storedPassword string, suppliedPassword string) (bool, error) {
	pwsalt := strings.Split(storedPassword, ".")
	if len(pwsalt) != 2 {
		return false, fmt.Errorf("Unable to verify user password")
	}

	// check supplied password salted with hash
	salt, err := hex.DecodeString(pwsalt[1])

	if err != nil {
		return false, fmt.Errorf("Unable to verify user password")
	}

	// Must use the same derivation as HashPassword.
	shash := pbkdf2.Key([]byte(suppliedPassword), salt, 4096, 32, sha256.New)

	return subtle.ConstantTimeCompare([]byte(hex.EncodeToString(shash)), []byte(pwsalt[0])) == 1, nil
}
