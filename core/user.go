package core

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argonMemoryUsage = 256 * 1024 //gets 256 bytes from ram for crptographic noise
	argonIterations  = 8          //number or iteration
	argonParallelism = 4          //use 4cores number of threads
	hashLenght       = 32         //final hash 32bytes
	saltLenght       = 16         //random salt of 16bytes
)

type Users struct {
	UserID int64  //
	Name   string // username
	Pass   []byte //salted hash of the password
	Answer []byte // answer to recovery/update password encrypted
	//tuning parameters for the hashing
	Salt              []byte // unique random data for salt 16b
	ArgonIteration    uint32
	ArgonMemory       uint32
	ArgonParrallelism uint8
}

func NewUser() *Users

// generates a salted hash and returns formatted string representation
func CreateSaltedHash(inputPassword string) (string, error) {
	salt := make([]byte, saltLenght)
	//random sal bytes
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate hashkey: %v", err)
	}
	//generating the salted hash
	hash := argon2.IDKey([]byte(inputPassword), salt, uint32(argonIterations), uint32(argonMemoryUsage), uint8(argonParallelism), uint32(hashLenght))
	//encoding b64 and returning the string format
	b64salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d$t=%d$p=%d$%s$%s",
		argon2.Version, argonMemoryUsage, argonIterations, argonParallelism, b64salt, b64Hash,
	)
	return encodedHash, nil
}

func AuthenticateUser(inputPass, storedHashEncoding string) (bool, error) {
	salt, err := parseSalt(storedHashEncoding)
	if err != nil {
		return false, fmt.Errorf("failed to authenticate user: %v", err)
	}
	storedHash, err := parseHash(storedHashEncoding)
	if err != nil {
		return false, fmt.Errorf("failed to authenticate user: %v", err)
	}
	///compute hash against the inputed password
	computedHash := argon2.IDKey([]byte(inputPass), salt, uint32(argonIterations), uint32(argonMemoryUsage), uint8(argonParallelism), uint32(hashLenght))
	//compare hash bytes against time attaks
	if subtle.ConstantTimeCompare(storedHash, computedHash) == 1 {
		return true, nil
	}
	return false, nil
}

// helper  parsing salt and hash
func parseSalt(encodedHash string) ([]byte, error) {
	hashData := strings.Split(encodedHash, "$")
	salt := hashData[6] //split creates empty string when input string starts with seprator
	saltBytes, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse salt: %v", err)
	}
	return saltBytes, nil
}
func parseHash(encodedHash string) ([]byte, error) {
	hashData := strings.Split(encodedHash, "$")
	hash := hashData[7]
	hashBytes, err := base64.RawStdEncoding.DecodeString(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to parse hash: %v", err)
	}
	return hashBytes, nil
}
