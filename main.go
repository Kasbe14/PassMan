package main

import (
	// "github.com/Kasbe14/PassMan/core"
	"log"

	"github.com/Kasbe14/PassMan/database"
)

func main() {

	/*db, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("successfully connected to the SQLite database")*/
	// var userPass string
	// // fakePass := "sunStone14"
	// fmt.Print("Enter password: ")
	// _, err := fmt.Scanf("%s\n", &userPass)
	// if err != nil {
	// 	fmt.Printf("error:%s\n", err)
	// 	return
	// }

	// hash, err := core.CreateSaltedHash(userPass)
	// if err != nil {
	// 	log.Fatalf("failed to create user: %v", err)
	// }
	// fmt.Printf("user created with pass hash: %s\n", hash)
	// var enterePassword string
	// fmt.Print("Enter password to authenticate:")
	// _, err = fmt.Scanf("%s", &enterePassword)
	// if err != nil {
	// 	fmt.Printf("error:%s", err)
	// 	return
	// }
	// Valid, err := core.AuthenticateUser(enterePassword, hash)
	// if err != nil {
	// 	log.Fatalf("user authentication failed: %v", err)
	// }
	// if !Valid && err == nil {
	// 	log.Fatalf("user authentication failed:")
	// }
	// if Valid {
	// 	fmt.Println("user authentication success")
	// }
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer database.CloseDB(db)
	if err := database.InitializeSchema(db); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

}
