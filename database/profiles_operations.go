package database

import (    
	"database/sql"
	_ "modernc.org/sqlite"
    "github.com/Kasbe14/PassMan/model" 
)

//profile Queryies
const (
    queryInsertProfile = `
    INSERT INTO profiles  (user_id, pro_hash, enc_pro_name, enc_pass, created_at, updated_at, lck, unlock_at)
    VALUES(?,?,?,?,?,?,?,?); 
    `               
    queryGetProfileByName = `
    SELECT user_id, pro_hash, enc_pro_name, enc_pass, created_at, updated_at, lck, unlock_at FROM profiles WHERE profile_hash = ?;
    `
    queryUserProfileCount = `
    SELECT COUNT(*) FROM profiles WHERE user_id = ?;
    `
    queryGetProfileNames = `
    SELECT enc_pro_name FROM profiles WHERE user_id = ?;
    `
)                           
                                

//inserts all the profile values and userid and assings profileid
func InsertProfile(db *sql.DB, p *model.Profile) error {
    
    result, err := db.Exec(queryInsertProfile,
                           p.UserID,
                           p.ProfileHash,
                           p.EncProfileName,
                           p.EncProfilePass,
                           p.CreatedAt,
                           p.UpdatedAt,
                           p.Locked,
                           p.UnlockAT,
                       )
    if err != nil {
        return err
    }
    p.ProfileID, err  = result.LastInsertId()
    if err != nil {
        return err
    }
	return nil
}

func GetProfileByName(db *sql.DB, profileHash string) (*model.Profile,error) {
    var p model.Profile
    err := db.QueryRow(queryGetProfileByName,profileHash).Scan(
           &p.UserID,
           &p.ProfileHash,
           &p.EncProfileName,
           &p.EncProfilePass,
           &p.CreatedAt,
           &p.UpdatedAt,
           &p.Locked,
           &p.UnlockAT)
  if err != nil {
      return  nil, err

  }
     return &p, nil
}
func GetUserProfileCount(db *sql.DB, userID int64) (int64,error) {
    var c int64
    err := db.QueryRow(queryUserProfileCount, userID).Scan(&c)
    if err != nil {
        return 0,err
    }
    return c,nil
}

func GetProfileNames(db *sql.DB, userID int64) ([][]byte,error) {
    var names [][]byte
    rows,err := db.Query(queryGetProfileNames,userID)
    if err != nil {
        return nil,err
    }
    defer rows.Close()
    for rows.Next() {
        var pEnc []byte
        if err := rows.Scan(&pEnc); err != nil {
            return nil, err
        }
        names = append(names,pEnc)
   
    }
    
    if err = rows.Err(); err != nil {
        return nil, err
    }
    return names,nil
}
