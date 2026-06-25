package core

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
    "github.com/Kasbe14/PassMan/model"

	"github.com/Kasbe14/PassMan/database"
)

var (
	ErrUserAlreadyExist    = errors.New("user already exist")
	ErrNotFound            = errors.New("record not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrProfileAlreadyExist = errors.New("profile already exist")
	ErrProfileNotFound     = errors.New("profile not found")
	ErrIncorrectPassword   = errors.New("incorrect password")
)

// provides database operations to the tui
type VaultService struct {
	db *sql.DB
}

func NewVaultService(db *sql.DB) *VaultService {
    return &VaultService{
        db: db,
    }
}

func (vs *VaultService) RegisterUser(username, inputpass  string) (string,error) {
	// 	//todo checkuserexist
	exists, _ := database.CheckUserExist(vs.db, username)
	if exists {
        return "",ErrUserAlreadyExist
	}
	saltedHash, err := CreateSaltedHash(inputpass)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %v", err)
	}
	masterKey, err := randomKey()
	if err != nil {
		return "",fmt.Errorf("failed to create user: %v", err)
	}
	//another saltedhash as derived key to encrypt the masterKey[key that encrypts data]
	saltEncrypt,aesKey, err := createAESKey(inputpass)	
    defer Wipe(aesKey) 
    defer Wipe(saltEncrypt)
    wrapPassKeyBytes, nonce, err := encryptData(masterKey, aesKey[:])
	// nonce appended wrapp key
	finalWrappedKeyPass := append(nonce, wrapPassKeyBytes...)

    //Creating user recovery string to encypt the masterKey 
    userRecString, err := generateUserRecoveryString()
    if err != nil {
        return "", fmt.Errorf("failed to create user: %v",err)
    }
    userRecStringAsKey := sha256.Sum256([]byte(userRecString))
    defer Wipe(userRecStringAsKey[:])
    wrapRecStringKeyBytes, nonce, err := encryptData(masterKey, userRecStringAsKey[:])
    if err != nil {
        return "",fmt.Errorf("failed to create user: %v",err)
    }
    finalWrappedKeyRec := append(nonce, wrapRecStringKeyBytes...)

	newUser := model.Users{
		Name:           username,
		PassHash:       saltedHash,
        EncryptSalt:     saltEncrypt,
		WrappedKeyPass: finalWrappedKeyPass,
        WrappedKeyRec: finalWrappedKeyRec,
	}

	err = database.InsertUser(vs.db, newUser)
	if err != nil {
		return "",fmt.Errorf("failed to create user: %v", err)
	}
    return userRecString, nil
}


func (vs *VaultService) LoginUser(username, userpassword string) ([]byte,error) {
	exists, _ := database.CheckUserExist(vs.db, username)
	if !exists {
        return nil, ErrUserNotFound
	}
    //userAuth 
	storedEncodedHashed,encryptSalt,wrapKey ,err := database.GetUserCredentials(vs.db, username)
    defer Wipe(encryptSalt)
    if err != nil {
        return nil,fmt.Errorf("failed to login %v",err)
    }
	valid, err := AuthenticateUser(userpassword, storedEncodedHashed)
	if !valid && err == nil {
		return nil, ErrIncorrectPassword
	}
	if err != nil {
		return nil,fmt.Errorf("failed to login: %v", err)
	}
    //deriving the master key a using the input password
    masterkey, err := unwrapMasKey(encryptSalt,wrapKey,userpassword)
    if err != nil {
        return nil, fmt.Errorf("failed to login: can't retrive key")
    }
     
	return masterkey,nil
}
