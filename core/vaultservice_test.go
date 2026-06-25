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
      db, err := sql.Open("sqlite", ":memory:?_foreign_keys=ON")
      if err != nil {
          t.Fatalf("failed to open test databse %v",err)
      }

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
    masterkey, err := vault.LoginUser("user1Login","testPass2")
    if err != nil {
        t.Fatalf("Login user failed: %v",err)
    }
    t.Logf("masterkey bytes returned [converted to string]: %s",string(masterkey))

    //false login
    _, err  = vault.LoginUser("falseUser","falsePass")
    if err == nil {
        t.Fatalf("login successful for a unregistered user")
    }
    t.Logf("error return for unregister user: %v",err)

    //Incorrect password
    _, err = vault.LoginUser("user1Login", "testpass2")
    if err == nil {
        t.Fatalf("login successful for a Incorrect password")
    }
    t.Logf("error returned for incorrect password: %v",err)
}
