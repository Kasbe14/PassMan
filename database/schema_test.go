package database

import (
    "database/sql"
    _ "modernc.org/sqlite"
    "testing"
    
    "github.com/Kasbe14/PassMan/model"
)

//helper setup
func setUpDB(t *testing.T) *sql.DB {
      t.Helper()
      db, err := sql.Open("sqlite", ":memory:?_foreign_keys=ON")
      if err != nil {
          t.Fatalf("failed to open test databse %v",err)
      }

      //initializing the schema
      err = InitializeSchema(db) 
      if err != nil {
          t.Fatalf("failed to initialize the test db with initial schema %v",err)
      }
      t.Cleanup(func (){
          // CloseDB(db)
          db.Close()
      })

      return db

}

//creating test db and initiallizing with the schema
func TestInitializeSchema(t *testing.T) {
      //opening and creating database with the  
      db, err := sql.Open("sqlite", ":memory:?_foreign_keys=ON")
      if err != nil {
          t.Fatalf("failed to open test databse %v",err)
      }

      //initializing the schema
      err = InitializeSchema(db) 
      if err != nil {
          t.Fatalf("failed to initialize the test db with initial schema %v",err)
      }
      t.Cleanup(func (){
          // CloseDB(db)
          db.Close()
      })
}
       
     

func TestInsertUser(t  *testing.T) {
    //setting up the initialize schema 
     db := setUpDB(t)

     testUser := &model.Users{
         Name: "Saurabh", 
         PassHash: "$argon2id$v=19$....",
         EncryptSalt: []byte("fakesalt"),
         WrappedKeyPass: []byte("fake_encrypted_bytes_pass"),
         WrappedKeyRec: []byte("fake_encrypted_bytes_rec"),
     }

     //testing the insert funtion
     err := InsertUser(db,testUser)
     if err != nil {
         t.Fatalf("InsertUser failed %v",err)
     }
     //cross check if insert was done
     var savedName string
     var savedHash  string
     err = db.QueryRow("SELECT username, master_hash FROM users WHERE username=?", "Saurabh").Scan(&savedName, &savedHash)
     if err != nil {
         if  err == sql.ErrNoRows {
             t.Fatalf("InsertUser completed but no rows found in the database")
         }
         t.Fatalf("failed to query test database: %v",err)
     }
     //verify datao
     if savedName != testUser.Name {
         t.Errorf("expected username %s, got %s",testUser.Name,savedName)
     }
     if savedHash != testUser.PassHash {
         t.Errorf("pass has mismatch: courruption")
     }
}

func TestCheckUserExist(t *testing.T) {

     db := setUpDB(t)
     testUser1 := &model.Users{
         Name: "Saurabh", 
         PassHash: "$argon2id$v=19$....",
         EncryptSalt: []byte("fakesalt1"),
         WrappedKeyPass: []byte("fake_encrypted_bytes_pass"),
         WrappedKeyRec: []byte("fake_encrypted_bytes_rec"),
     }
     testUser2 := &model.Users{
         Name: "Vivek", 
         PassHash: "$argon2id$v=20$....",
         EncryptSalt: []byte("fakesalt2"),
         WrappedKeyPass: []byte("fake_encrypted_bytes_pass_of_vi"),
         WrappedKeyRec: []byte("fake_encrypted_bytes_rec_of_vi"),
     }
     err := InsertUser(db,testUser1)
     if err != nil {
         t.Fatalf("InsertUser failed %v",err)
     }
     err = InsertUser(db,testUser2)
     if err != nil {
         t.Fatalf("InsertUser failed %v",err)
     }
     //testing the checkuser
     
     exits, err := CheckUserExist(db, testUser1.Name)
     if !exits && err == nil {
         t.Errorf("CheckUserExist failed expected  user to exits but got not exists")
     }
     if err != nil {
         t.Fatalf("CheckUserExist failed %v",err)
     }
     exits, err = CheckUserExist(db, "NonExistUser")
     if exits && err ==  nil {
         t.Errorf("CheckUserExist failed expected user to not exists but got exist")
     }
     if err != nil {
         t.Fatalf("CheckUserExist failed %v",err)
     }
    
}
// TODO ::test getusercredentials
