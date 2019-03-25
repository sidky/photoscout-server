package profile

import (
	"database/sql"
	"log"
	"os"
)

type User struct {
	UUID string
}

func (user User) BookmarkPhoto(photoId string) error {
	updateStatement := `
		INSERT INTO starred_photo(user_id, photo_id) VALUES($1, $2) ON CONFLICT 
		DO UPDATE SET updated_at = now() WHERE user_id = $1 AND photo_id = $2
	`
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Print("Unable to close connection: %v", err)
		}
	}()

	_, err = db.Exec(updateStatement, user.UUID, photoId)
	return nil
}