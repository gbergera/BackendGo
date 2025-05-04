package followRepo

import (
	"database/sql"
	"fmt"
	"time"
	"ualabackend/entities/follow"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Create(followerID, followedID int) error {
	_, err := r.DB.Exec("INSERT INTO follows (follower_id, followed_id) VALUES (?, ?)", followerID, followedID)
	if err != nil {
		fmt.Println("Error en INSERT:", err)
		return err
	}

	_, err = r.DB.Exec(`
		UPDATE users SET followers_id = JSON_ARRAY() WHERE id = ? AND followers_id IS NULL
	`, followedID)
	if err != nil {
		fmt.Println("Error inicializando followers_id:", err)
		return err
	}

	_, err = r.DB.Exec(`
		UPDATE users SET following_id = JSON_ARRAY() WHERE id = ? AND following_id IS NULL
	`, followerID)
	if err != nil {
		fmt.Println("Error inicializando following_id:", err)
		return err
	}

	_, err = r.DB.Exec(`
		UPDATE users SET followers_id = JSON_ARRAY_APPEND(followers_id, '$', ?) WHERE id = ?
	`, followerID, followedID)
	if err != nil {
		fmt.Println("Error en UPDATE followers_id:", err)
		return err
	}

	_, err = r.DB.Exec(`
		UPDATE users SET following_id = JSON_ARRAY_APPEND(following_id, '$', ?) WHERE id = ?
	`, followedID, followerID)
	if err != nil {
		fmt.Println("Error en UPDATE following_id:", err)
		return err
	}

	fmt.Println("Follow creado con éxito")
	return nil
}


func (r *Repository) GetAll() ([]follow.Follow, error) {
	query := `SELECT follower_id, followed_id FROM follows`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var follows []follow.Follow
	for rows.Next() {
		var f follow.Follow
		if err := rows.Scan(&f.FollowerID, &f.FollowedID); err != nil {
			return nil, err
		}
		f.Timestamp = time.Now() // opcional: si querés mostrar hora actual
		follows = append(follows, f)
	}

	return follows, nil
}

func (r *Repository) GetByIDs(followerID, followedID int) (*follow.Follow, error) {
	query := `SELECT follower_id, followed_id FROM follows WHERE follower_id = ? AND followed_id = ?`
	row := r.DB.QueryRow(query, followerID, followedID)

	var f follow.Follow
	if err := row.Scan(&f.FollowerID, &f.FollowedID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	f.Timestamp = time.Now()
	return &f, nil
}

func (r *Repository) Delete(followerID, followedID int) error {
	query := `DELETE FROM follows WHERE follower_id = ? AND followed_id = ?`
	_, err := r.DB.Exec(query, followerID, followedID)
	return err
}

func (r *Repository) GetFollowedByFollowerID(followerID int) ([]follow.Follow, error) {
	query := `SELECT follower_id, followed_id FROM follows WHERE follower_id = ?`
	rows, err := r.DB.Query(query, followerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var follows []follow.Follow
	for rows.Next() {
		var f follow.Follow
		if err := rows.Scan(&f.FollowerID, &f.FollowedID); err != nil {
			return nil, err
		}
		f.Timestamp = time.Now() // opcional
		follows = append(follows, f)
	}
	return follows, nil
}
