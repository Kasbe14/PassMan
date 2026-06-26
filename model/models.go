package model

// import "time"


type Users struct {
	UserID         int64  //
	Name           string // username
	PassHash       string //salted hash of the password
    EncryptSalt    []byte //salt for derived key for masterkey
	WrappedKeyPass []byte //masterkey locked by masterpassword [encrypted ciphertext AES-GCM nonce 12bytes appended at front]
	WrappedKeyRec  []byte //master key lockey by string given to user  nonce 12bytes appended at front
}

type Profile struct {
	ProfileID      int64
	UserID         int64  //foreign key from the user
	ProfileHash    string //hmac-Sha256 of profile name blind indexing
	EncProfileName []byte // profile name for a account eg gmail, instagram etc[encrypted ciphertext AES-GCM nonce 12bytes appended at front]
	EncProfilePass []byte // password for the profile [encypted ciphertext AES-GCM]
	CreatedAt      int64  //unix time
	UpdatedAt      int64 //unit time
	UnlockAT       int64 //locked password will be unlocked after unix epoch integer
	Locked         bool  //indicates locked password
}

//dto
type DecryptedProfile struct {
    Name string
    Password string
    CreatedAt  int64
    UpdatedAt int64
    Locked     bool
    UnlockAt   int64
}
