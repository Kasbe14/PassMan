package core

import (
    "testing"
    "database/sql"
    _ "modernc.org/sqlite"
    // "github.com/Kasbe14/PassMan/model"
    "github.com/Kasbe14/PassMan/database"
)

//helper setup
func setUpDB(t *testing.T) *sql.DB {
      t.Helper()
      //shared samedatabse accorss the test, foreign key is on
      db, err := sql.Open("sqlite", "file::memory:?cache=shared&_pragma=foreign_keys(1)")
      if err != nil {
          t.Fatalf("failed to open test databse %v",err)
      }
      //only one connection
      db.SetMaxOpenConns(1)

      //initializing the schema
      err = database.InitializeSchema(db) 
      if err != nil {
          t.Fatalf("failed to initialize the test db with initial schema %v",err)
      }
      t.Cleanup(func (){
          // CloseDB(db)
          db.Close()
      })

      return db

}

func TestRegisterUser(t *testing.T) {
    db := setUpDB(t)
    //testing registeruser
    vault := NewVaultService(db)
    //New registration test
    userString, err := vault.RegisterUser("user1", "testPass1")

    if err != nil {
        t.Fatalf("RegisterUser failed to register user %v",err)
    }
    t.Logf("the user string for recovery return : %s",userString)
    //already existing user registration test
    _, err = vault.RegisterUser("user1","testPass1")
    if err == nil {
        t.Fatalf("RegisterUser failed: registration for already existing user successful")
    }
    t.Logf("error returned for already register user: %v",err)
}


func TestLoginUser(t *testing.T) {
    db := setUpDB(t)
    //testing registeruser
    vault := NewVaultService(db)
    _, err := vault.RegisterUser("user1Login", "testPass2")

    if err != nil {
        t.Fatalf("RegisterUser failed to register user %v",err)
    }
    masterkey, userID,err := vault.LoginUser("user1Login","testPass2")
    if err != nil {
        t.Fatalf("Login user failed: %v",err)
    }
    t.Logf("masterkey bytes returned [converted to string]: %s",string(masterkey))
    t.Logf(" returned [userID]: %d",userID)

    //false login
    _, _, err  = vault.LoginUser("falseUser","falsePass")
    if err == nil {
        t.Fatalf("login successful for a unregistered user")
    }
    t.Logf("error return for unregister user: %v",err)

    //Incorrect password
    _,_, err = vault.LoginUser("user1Login", "testpass2")
    if err == nil {
        t.Fatalf("login successful for a Incorrect password")
    }
    t.Logf("error returned for incorrect password: %v",err)
}

func TestAddNormalProfile(t *testing.T) {
    db := setUpDB(t)
    vault := NewVaultService(db)
    _, err := vault.RegisterUser("user1Login", "testPass2")

    if err != nil {
        t.Fatalf("RegisterUser failed to register user %v",err)
    }
    masterkey, userID,err := vault.LoginUser("user1Login","testPass2")
    if err != nil {
        t.Fatalf("Login user failed: %v",err)
    }
    //testing happy path
    err = vault.AddNormalProfile(userID,"netflix","testProfilePass",masterkey)
    if  err != nil {
        t.Fatalf("adding normal profile failed  %v",err)
    }
    t.Logf("the err return after adding [%v]",err)
    //testing sad path
     //wrong user id
    err = vault.AddNormalProfile(3,"netflix","testProfilePass",masterkey)
    if err == nil {
        t.Fatalf("add normal worked for non-existing user")
    }
     t.Logf("error for wrong user [%v]",err)

}


func TestGetProfileByName(t *testing.T) {
    db := setUpDB(t)
    vault := NewVaultService(db)
    _, err := vault.RegisterUser("user1Login", "testPass2")

    if err != nil {
        t.Fatalf("RegisterUser failed to register user %v",err)
    }
    masterkey, userID,err := vault.LoginUser("user1Login","testPass2")
    if err != nil {
        t.Fatalf("Login user failed: %v",err)
    }
    //testing happy path

    err = vault.AddNormalProfile(userID,"netflix","testProfilePass",masterkey)
    if  err != nil {
        t.Fatalf("adding normal profile failed")
    }
    decryptProfile, err := vault.GetProfileByName("netflix", masterkey)
    if err != nil {
        t.Fatalf("getProfileByName failed [%v]",err)
    }else{
        t.Logf("returned Profile %+v",*decryptProfile)
    }
    //sad path
    //wrong profile
    _, err = vault.GetProfileByName("netFlixl",masterkey)
    if err == nil {
        t.Fatalf("Get profile by name failed for incorrect profile name")
    }else {
        t.Logf("Error for Wrong Profile Name: %v",err)
    }
}
