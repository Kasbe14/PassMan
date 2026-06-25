package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	// "crypto/subtle"
	"encoding/hex"
	"fmt"
	"io"
	"runtime"

	"golang.org/x/crypto/argon2"
)

const (
	masterKeySize = 32 //32byteskey
    recStringByteSize = 16
)

func encryptData(data []byte, key []byte) ([]byte, []byte, error) {
	//cipher blockW
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt: %v", err)
	}
	aegcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt: %v", err)
	}
	//nonce
	nonce := make([]byte, aegcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt data: %v", err)
	}
	encryptedData := aegcm.Seal(nil, nonce, []byte(data), nil)

	return encryptedData, nonce, nil

}

func decryptData(combinedData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %v", err)
	}
	aegcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %v", err)
	}
    nonceSize := aegcm.NonceSize()
    if len(combinedData) < nonceSize{
        return nil,fmt.Errorf("failed to decrypt: encypted data too short to contain nonce")
    }
    nonce :=  combinedData[:nonceSize]
    cipherText := combinedData[nonceSize:]
	decryptedData, err := aegcm.Open(nil, nonce, []byte(cipherText), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %v", err)
	}
	return decryptedData, nil

}

// returns 32 bytes key
func randomKey() ([]byte, error) {
	//encryption key for profiles and sensitive data
	masterKey := make([]byte, masterKeySize)
	if _, err := rand.Read(masterKey); err != nil {
		return nil, fmt.Errorf("failed to create random key: %v", err)
	}
	return masterKey, nil
}
//retuns 32hex charachters string
func generateUserRecoveryString() (string, error) {
    buf := make([]byte, recStringByteSize) 
    _, err := rand.Read(buf)
    if err != nil {
        return "",err
    }
    return hex.EncodeToString(buf), nil
}
//returns the salt and aes key to encrypt the masterkey
func createAESKey(inputPassword string) ([]byte, []byte, error) {
	salt := make([]byte, saltLenght)
	//random sal bytes
	if _, err := rand.Read(salt); err != nil {
		return nil,nil, fmt.Errorf("failed to generate salt: %v", err)
	}
	//generating the salted hash as aes key to wrap masterkey
    aesKey := argon2.IDKey([]byte(inputPassword), salt, uint32(argonIterations), uint32(argonMemoryUsage), uint8(argonParallelism), uint32(hashLenght))
    return salt, aesKey, nil
}

//returns masterkey
func unwrapMasKey(salt,wrappedKey []byte, userpass string) ([]byte,error) {
    aesKey := argon2.IDKey([]byte(userpass), salt, uint32(argonIterations), uint32(argonMemoryUsage), uint8(argonParallelism), uint32(hashLenght))
    masterkey, err := decryptData(wrappedKey,aesKey)
    if err != nil {
        return nil,fmt.Errorf("failed to derive key check username")
    }
	return masterkey, nil
}

// wipe the critcal data from ram instanlty
func Wipe(b []byte) {
    if  len(b) == 0 {
        return
    }
    for i := range b {
        b[i] = 0
    }
    _ = b[0]
    runtime.KeepAlive(b)
}
