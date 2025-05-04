package tweetRepo

import (
	"database/sql"
	"encoding/json"
	"time"
	"ualabackend/entities/tweet"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) Create(authorID int, message string) error {

	query := `INSERT INTO tweets (author_id, message, timestamp) VALUES (?, ?, ?)`
	result, err := r.DB.Exec(query, authorID, message, time.Now())
	if err != nil {
		return err
	}

	tweetID, err := result.LastInsertId()
	if err != nil {
		return err
	}


	var followersJSON string
	err = r.DB.QueryRow("SELECT followers_id FROM users WHERE id = ?", authorID).Scan(&followersJSON)
	if err != nil {
		return err
	}


	var followers []int
	err = json.Unmarshal([]byte(followersJSON), &followers)
	if err != nil {
		return err
	}

for _, followerID := range followers {
	_, err := r.DB.Exec(`
		UPDATE users 
		SET feed = JSON_ARRAY() 
		WHERE id = ? AND feed IS NULL
	`, followerID)
	if err != nil {
		return err
	}


	_, err = r.DB.Exec(`
		UPDATE users 
		SET feed = JSON_ARRAY_APPEND(feed, '$', ?) 
		WHERE id = ?
	`, tweetID, followerID)
	if err != nil {
		return err
	}
}

	return nil
}

func (r *Repository) GetAll() ([]tweet.Tweet, error) {
	query := `SELECT id, author_id, message, timestamp FROM tweets ORDER BY timestamp DESC`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []tweet.Tweet
	for rows.Next() {
		var t tweet.Tweet
		if err := rows.Scan(&t.Id, &t.Author_id, &t.Message, &t.Timestamp); err != nil {
			return nil, err
		}
		tweets = append(tweets, t)
	}
	return tweets, nil
}

func (r *Repository) GetByID(id int) (*tweet.Tweet, error) {
	query := `SELECT id, author_id, message, timestamp FROM tweets WHERE id = ?`
	row := r.DB.QueryRow(query, id)

	var t tweet.Tweet
	if err := row.Scan(&t.Id, &t.Author_id, &t.Message, &t.Timestamp); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *Repository) Update(id int, newMessage string) error {
	query := `UPDATE tweets SET message = ?, timestamp = ? WHERE id = ?`
	_, err := r.DB.Exec(query, newMessage, time.Now(), id)
	return err
}

func (r *Repository) Delete(id int) error {
	query := `DELETE FROM tweets WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	return err
}
