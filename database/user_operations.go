package database

import (
	"database/sql"
	"fmt"

	"github.com/Kasbe14/PassMan/model"
	_ "modernc.org/sqlite"
)

//queries
const  (
    queryCheckUserExist = `
     SELECT EXISTS (SELECT 1 FROM users WHERE username = ?) 
    `
    queryGetUserCreds = `
    SELECT master_hash, encrypt_salt, WrappedKeyPass FROM users WHERE username = ? LIMIT 1
    `
    queryInsertUser = `
    INSERT INTO users (username, master_hash, encrypt_salt, WrappedKeyPass, WrappedKeyRec) VALUES (?,?,?,?,?)
    `
)

func CheckUserExist(db *sql.DB, username string) (bool, error) {
    var exist bool
    err := db.QueryRow(queryCheckUserExist, username).Scan(&exist)
    if err != nil {
        return false, err
    }
	return exist, nil
}

//returns the stored salted hash and encrypt salt for login and wrapped key for deriving key
func GetUserCredentials(db *sql.DB, username string) (string,[]byte, []byte,error) {
    var saltedHash string
    var encryptSalt []byte
    var wrapKey      []byte
    err := db.QueryRow(queryGetUserCreds, username).Scan(&saltedHash,&encryptSalt,&wrapKey)
    if err != nil {
        return "",nil, nil,fmt.Errorf("failed to get user Credentials %v",err)
    }
	return saltedHash,encryptSalt,wrapKey, nil
}
func InsertUser(db *sql.DB, user model.Users) error {
    result, err := db.Exec(queryInsertUser,user.Name,user.PassHash,user.EncryptSalt,user.WrappedKeyPass,user.WrappedKeyRec)
    if err != nil {
        return err
    }
    user.UserID, err  = result.LastInsertId()
    if err != nil {
        return err
    }
	return nil
}
