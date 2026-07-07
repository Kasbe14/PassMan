package core

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
     // "github.com/Kasbe14/PassMan/model"
	"golang.org/x/crypto/argon2"
)

const (
	argonMemoryUsage = 256 * 1024 //gets 256 bytes from ram for crptographic noise
	argonIterations  = 10         //number or iteration
	argonParallelism = 4          //use 4cores number of threads
	hashLenght       = 32         //final hash 32bytes
	saltLenght       = 16         //random salt of 16bytes
)


// func NewUser(username, password, answer string) *Users {
// 	return &Users{
// 		Name: username,
// 		Pass: ,
// 	}
// }

// generates a salted hash and returns formatted string representation
func CreateSaltedHash(inputString string) (string /*, []byte*/, error) {
	salt := make([]byte, saltLenght)
	//random sal bytes
	if _, err := rand.Read(salt); err != nil {
		return "" /*, nil*/, fmt.Errorf("failed to generate salt: %v", err)
	}
	//generating the salted hash
	hash := argon2.IDKey([]byte(inputString), salt, uint32(argonIterations), uint32(argonMemoryUsage), uint8(argonParallelism), uint32(hashLenght))
	//encoding b64 and returning the string format
	b64salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d$t=%d$p=%d$%s$%s",
		argon2.Version, argonMemoryUsage, argonIterations, argonParallelism, b64salt, b64Hash,
	)
	return encodedHash /* salt,*/, nil
}

func AuthenticateUser(inputString, storedHashEncoding string) (bool, error) {
	argonIterations, argonMemoryUsage, argonParallelism, err := parseParameters(storedHashEncoding)
	if err != nil {
		return false, fmt.Errorf("failed to authenticate user: %v", err)
	}
	salt, err := parseSalt(storedHashEncoding)
	if err != nil {
		return false, fmt.Errorf("failed to authenticate user: %v", err)
	}
	storedHash, err := parseHash(storedHashEncoding)
	if err != nil {
		return false, fmt.Errorf("failed to authenticate user: %v", err)
	}
	///compute hash against the inputed password
	computedHash := argon2.IDKey([]byte(inputString), salt, argonIterations, argonMemoryUsage, argonParallelism, hashLenght)
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
func parseParameters(encodedHash string) (iter, mem uint32, thr uint8, err error) {
	parameters := strings.Split(encodedHash, "$")

	memStr := strings.TrimPrefix(parameters[3], "m=")
	m, err := strconv.ParseUint(memStr, 10, 32)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse memory parameter: %v", err)
	}
	iterStr := strings.TrimPrefix(parameters[4], "t=")
	i, err := strconv.ParseUint(iterStr, 10, 32)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse iteration parameter: %v", err)
	}
	thrStr := strings.TrimPrefix(parameters[5], "p=")
	t, err := strconv.ParseUint(thrStr, 10, 8)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to parse threads parameter: %v", err)
	}
	iter = uint32(i)
	mem = uint32(m)
	thr = uint8(t)

	return iter, mem, thr, nil
}
