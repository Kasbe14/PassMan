package core

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
    "time"
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
    finalWrappedKeyPass,  err := encryptData(masterKey, aesKey[:])

    //Creating user recovery string to encypt the masterKey 
    userRecString, err := generateUserRecoveryString()
    if err != nil {
        return "", fmt.Errorf("failed to create user: %v",err)
    }
    userRecStringAsKey := sha256.Sum256([]byte(userRecString))
    defer Wipe(userRecStringAsKey[:])

    finalWrappedKeyRec, err := encryptData(masterKey, userRecStringAsKey[:])
    if err != nil {
        return "",fmt.Errorf("failed to create user: %v",err)
    }

	newUser := &model.Users{
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


func (vs *VaultService) LoginUser(username, userpassword string) ([]byte,int64,error) {
	exists, _ := database.CheckUserExist(vs.db, username)
	if !exists {
        return nil,0 ,ErrUserNotFound
	}
    //userAuth 
	storedEncodedHashed,encryptSalt,wrapKey,userID ,err := database.GetUserCredentials(vs.db, username)
    defer Wipe(encryptSalt)
    defer Wipe(wrapKey)
    if err != nil {
        return nil,0,fmt.Errorf("failed to login %v",err)
    }
	valid, err := AuthenticateUser(userpassword, storedEncodedHashed)
	if !valid && err == nil {
		return nil,0, ErrIncorrectPassword
	}
	if err != nil {
		return nil,0,fmt.Errorf("failed to login: %v", err)
	}
    //deriving the master key a using the input password
    masterkey, err := unwrapMasKey(encryptSalt,wrapKey,userpassword)
    if err != nil {
        return nil, 0,fmt.Errorf("failed to login: can't retrive key")
    }
     
	return masterkey,userID,nil
}

//TODO -> tui will genrate pas every 15 sec dicewar when done signal then call this functton
//normal profile [unlocked]
func (v *VaultService) AddLockedProfile(days, userID int64, proName,proPass string, mKey []byte) error {
    //TODO check if already exists
    lock := true
    unlockAT := time.Now().Unix() + (days * 86400)
    createdAt := time.Now().Unix()
    updatedAt := time.Now().Unix()
    proBlindHash := createProfileBlindHash(proName,mKey)
    encryptedProName,err := encryptData([]byte(proName),mKey)
    if err != nil {
        return fmt.Errorf("failed to add profile: %v",err)
    }
    encryptedProPass,err := encryptData([]byte(proPass),mKey)
    defer Wipe([]byte(proName))
    if err != nil {
        return fmt.Errorf("failed to add profile: %v",err)
    }
    profile := &model.Profile{
        UserID:         userID,        
        ProfileHash:   proBlindHash, 
        EncProfileName:encryptedProName, 
        EncProfilePass:encryptedProPass, 
        CreatedAt:     createdAt, 
        UpdatedAt:     updatedAt, 
        UnlockAT:      unlockAT, 
        Locked:        lock,
     }
     err = database.InsertProfile(v.db,profile)
     if err != nil {
         return fmt.Errorf("failed to add profile: %v",err)
     }
     return nil
}

func (v *VaultService) AddNormalProfile(userID int64, proName,proPass string, mKey []byte) error {
    //TODO checkProfileExists
     proBlindHash := createProfileBlindHash(proName,mKey)
     encryptedProName,err := encryptData([]byte(proName),mKey)
     if err != nil {
         return fmt.Errorf("failed to add profile: %v",err)
     }
     encryptedProPass,err := encryptData([]byte(proPass),mKey)
     defer Wipe([]byte(proName))
     if err != nil {
         return fmt.Errorf("failed to add profile: %v",err)
     }
     createdAt := time.Now().Unix()
     updatedAt := time.Now().Unix()
     locked := false
     profile := &model.Profile{
         UserID:         userID,        
         ProfileHash:   proBlindHash, 
         EncProfileName:encryptedProName, 
         EncProfilePass:encryptedProPass, 
         CreatedAt:     createdAt, 
         UpdatedAt:     updatedAt, 
         UnlockAT:      0, 
         Locked:        locked,
      }
      err = database.InsertProfile(v.db,profile)
      if err != nil {
          return fmt.Errorf("failed to add profile: %v",err)
      }
      return nil
} 

func (v *VaultService) GetProfileByName(profileName string, mKey []byte) (*model.DecryptedProfile,error) {
    searchHash := createProfileBlindHash(profileName,mKey)
    pEnc, err := database.GetProfileByName(v.db,searchHash)
    if err != nil {
        if errors.Is(err,sql.ErrNoRows) {
        return nil, fmt.Errorf("profile not found enter correct name")
        }
        return nil, fmt.Errorf("failed to get profile %v",err)
    }
    //decrypt the data returned
    pName , err:= decryptData(pEnc.EncProfileName,mKey) 
    if err != nil {
        return nil, fmt.Errorf("failed to get profile %v",err)
    }
    var pPass string
    now := time.Now().Unix()
    if pEnc.Locked &&  now < pEnc.UnlockAT {
        timeLeft := time.Duration(pEnc.UnlockAT - now) * time.Second
        pPass = fmt.Sprintf("Locked Until %s", timeLeft.String())
    }else {
        pass, err := decryptData(pEnc.EncProfilePass,mKey)
        if err != nil {
            return nil, fmt.Errorf("failed to get profile %v",err)
        }
        pPass = string(pass)
    }

    //dto for user
    p := &model.DecryptedProfile{
        Name: string(pName),
        Password: pPass,
        CreatedAt: pEnc.CreatedAt,
        UpdatedAt: pEnc.UpdatedAt,
        Locked:    pEnc.Locked,
        UnlockAt:  pEnc.UnlockAT,
    }
   
    return p ,nil
   
}
 
//returns all the profile names of the user
func (v *VaultService) GetProfileNameList(userID int64, mKey []byte) ([]string,error) {
       encNames, err := database.GetProfileNames(v.db, userID)
       if err != nil {
           return nil, fmt.Errorf("failed to get profile names")
       }
       var decNames []string
       for _, encName := range encNames {
           decName ,err:= decryptData(encName,mKey)
           if err != nil {
               return  nil, fmt.Errorf("failed to get profile names")
           }
           decNames = append(decNames,string(decName))
       }
       return decNames, nil
}

