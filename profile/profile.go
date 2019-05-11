package profile

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type User struct {
	UUID string
}

func (user User) BookmarkPhoto(photoId string) error {
	updateStatement := `
		INSERT INTO bookmarks(user_id, photo_id) VALUES($1, $2) ON CONFLICT (user_id, photo_id)
		DO UPDATE SET updated_at = now() WHERE bookmarks.user_id = $1 AND bookmarks.photo_id = $2
	`
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Printf("Unable to close connection: %v", err)
		}
	}()

	_, err = db.Exec(updateStatement, user.UUID, photoId)
	return err
}

func (user *User) GetBookmarkedPhotos() ([]string, error) {
	query := `SELECT photo_id FROM bookmarks WHERE user_id = $1`

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Print("Unable to close connection: %v", err)
		}
	}()

	rows, err := db.Query(query, user.UUID)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("Unable to close row: %v", err)
		}
	}()

	photoIds := make([]string, 0)

	for rows.Next() {
		var photoId string
		err := rows.Scan(&photoId)
		if err != nil {
			return nil, err
		}

		photoIds = append(photoIds, photoId)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return photoIds, nil
}

func (user *User) GetBookmarkedPhotosFromSet(photoIds []string) ([]string, error) {
	if len(photoIds) == 0 {
		return nil, nil
	}
	query := "SELECT photo_id FROM bookmarks WHERE user_id = ? AND photo_id IN (?" + strings.Repeat(",?", len(photoIds) - 1) + ")"
	args := []interface{}{user.UUID}

	sqlArgs := make([]interface{}, len(photoIds))
	for n, photoId := range photoIds {
		sqlArgs[n] = photoId
	}
	args = append(args, sqlArgs...)

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Print("Unable to close connection: %v", err)
		}
	}()

	rows, err := db.Query(query, args)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("Unable to close row: %v", err)
		}
	}()

	bookmarkedIds := make([]string, 0)

	for rows.Next() {
		var photoId string
		err := rows.Scan(&photoId)
		if err != nil {
			return nil, err
		}

		bookmarkedIds = append(bookmarkedIds, photoId)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return bookmarkedIds, nil
}

func (user *User) IsBookmarkedPhoto(photoId string) (bool, error) {
	query := "SELECT COUNT(photo_id) FROM bookmarks WHERE user_id = $1 AND photo_id = $2"
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return false, err
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Print("Unable to close connection: %v", err)
		}
	}()

	rows, err := db.Query(query, user.UUID, photoId)
	if err != nil {
		return false, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Printf("Unable to close row: %v", err)
		}
	}()

	var found bool

	if rows.Next() {
		var c int64
		err := rows.Scan(&c)
		if err != nil {
			return false, err
		}
		if c != 0 {
			found = true
		} else {
			found = false
		}
	}
	if rows.Err() != nil {
		return false, rows.Err()
	}
	return found, nil
}

