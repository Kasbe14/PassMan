package core

import (
	"crypto/hmac"
	"crypto/sha256"
	// "fmt"

	// "time"
	"encoding/base64"

	// "github.com/Kasbe14/PassMan/database"
	// "github.com/Kasbe14/PassMan/model"
)

//cretes a sha256hmac hash for profile blindindex using the masterkey
func  createProfileBlindHash(profileName  string, masterKey []byte) (string) {
    hashfucntion := hmac.New(sha256.New, masterKey)
    hashfucntion.Write([]byte(profileName))
    return base64.StdEncoding.EncodeToString(hashfucntion.Sum(nil))
}

//TODO

    
// func getDecryptedProfilePass()
// func getDecryptedProfile()
