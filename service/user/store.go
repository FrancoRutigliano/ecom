package user

import (
	"database/sql"
	"fmt"

	typePayload "github.com/FrancoRutigliano/ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*typePayload.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	u := new(typePayload.User)

	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func scanRowIntoUser(rows *sql.Rows) (*typePayload.User, error) {
	user := new(typePayload.User)

	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) GetUserByID(id int) (*typePayload.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	u := new(typePayload.User)

	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func (s *Store) CreateUser(user typePayload.User) error {
	_, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password) VALUES (?, ?, ?, ?)", user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}
