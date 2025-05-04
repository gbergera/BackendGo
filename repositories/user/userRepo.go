package userRepo

import (
	"database/sql"
	"encoding/json"
	"ualabackend/entities/user"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Create(username string) error {
	_, err := r.DB.Exec(`
		INSERT INTO users (name, followers_id, following_id, feed)
		VALUES (?, '[]', '[]', '[]')
	`, username)
	return err
}

func (r *Repository) GetAll() ([]user.User, error) {
	query := `SELECT id, name, followers_id, following_id, feed FROM users ORDER BY id ASC`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		var followersID, followingID, feed sql.NullString

		if err := rows.Scan(&u.Id, &u.Name, &followersID, &followingID, &feed); err != nil {
			return nil, err
		}

		if followersID.Valid {
			u.Followers_id = json.RawMessage(followersID.String)
		}
		if followingID.Valid {
			u.Following_id = json.RawMessage(followingID.String)
		}
		if feed.Valid {
			u.Feed = json.RawMessage(feed.String)
		}

		users = append(users, u)
	}
	return users, nil
}

func (r *Repository) GetByID(id int) (*user.User, error) {
	query := `SELECT id, name, followers_id, following_id, feed FROM users WHERE id = ?`
	row := r.DB.QueryRow(query, id)

	var u user.User
	var followersID, followingID, feed sql.NullString

	if err := row.Scan(&u.Id, &u.Name, &followersID, &followingID, &feed); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if followersID.Valid {
		u.Followers_id = json.RawMessage(followersID.String)
	}
	if followingID.Valid {
		u.Following_id = json.RawMessage(followingID.String)
	}
	if feed.Valid {
		u.Feed = json.RawMessage(feed.String)
	}

	return &u, nil
}

func (r *Repository) Update(id int, newName string) error {
	query := `UPDATE users SET name = ? WHERE id = ?`
	_, err := r.DB.Exec(query, newName, id)
	return err
}

func (r *Repository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	return err
}
