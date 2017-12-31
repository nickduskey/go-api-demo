package models

import "database/sql"

// User represents a user
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// GetUser fetches a user from the db
func (u *User) GetUser(db *sql.DB) error {
	return db.QueryRow("SELECT name, price FROM users WHERE id=$1",
		u.ID).Scan(&u.Username, &u.Password)
}

// UpdateUser updates a user in the db
func (u *User) UpdateUser(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE users SET username=$1, password=$2 where id=$3",
			u.Username, u.Password, u.ID)

	return err
}

// DeleteUser deletes a user
func (u *User) DeleteUser(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users WHERE id=$1, u.ID")

	return err
}

// CreateUser inserts a user in the db
func (u *User) CreateUser(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO users(name, price) VALUES($1, $2) RETURNING id",
		u.Username, u.Password).Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}

// GetUsers retrieves users from the db
func GetUsers(db *sql.DB, start, count int) ([]User, error) {
	rows, err := db.Query(
		"SELECT id, username, password FROM users LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Password); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
