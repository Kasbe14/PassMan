package core

import "time"

type Profile struct {
	ProfileID      int64
	UserID         int64  //foreign key from the user
	ProfileHash    []byte //hmac-Sha256 of profile name blind indexing
	EncProfileName []byte // profile name for a account eg gmail, instagram etc[encrypted ciphertext AES-GCM nonce 12bytes appended at front]
	EncProfilePass []byte // password for the profile [encypted ciphertext AES-GCM]
	CreatedAt      time.Time
	UpdatedAt      time.Time
	UnlockAT       int64 //locked password will be unlocked after unix epoch integer
	Locked         bool  //indicates locked password
}
