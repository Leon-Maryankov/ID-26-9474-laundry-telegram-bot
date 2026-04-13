package users

import (
	"database/sql"
	"fmt"
)

type Users struct {
	database *sql.DB
}

type User struct {
	Id int `json:"id"`
	TelegramId int `json:"telegramId"`
	Name string `json:"name"`
	Surname string `json:"surname"`
	Email string `json:"email"`
}

func NewUsers(database *sql.DB) (*Users, error) {
	_, err := database.Exec(
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER UNIQUE NOT NULL PRIMARY KEY,
			telegramId INTEGER UNIQUE NOT NULL,
			name TEXT NOT NULL,
			surname TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL
		)`,
	)
	if err != nil {
		return nil, fmt.Errorf("users.NewUsers: failed to initialize database: %w", err)
	}

	result := &Users{
		database,
	}

	return result, nil
}

func (users *Users) AddUser(user *User) error {
	_, err := users.database.Exec(
		`INSERT INTO users (telegramId, name, surname, email) values (?, ?, ?, ?);`,
		user.TelegramId,
		user.Name, user.Surname,
		user.Email,
	)
	if err != nil {
		return fmt.Errorf("users.addUser: failed to add user: %w", err)
	}

	return nil
}

func (users *Users) GetUserByTelegramId(telegramId int) (*User, error) {
	result := &User{}
	err := users.database.QueryRow(
		`SELECT * FROM users WHERE telegramId = ?;`,
		telegramId,
	).Scan(
		&result.Id,
		&result.TelegramId,
		&result.Name, &result.Surname,
		&result.Email,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("users.GetUserByTelegramId: failed to query row: %w", err)
	}

	return result, nil
}

func (users *Users) GetUserById(id int) (*User, error) {
	result := &User{}
	err := users.database.QueryRow(
		`SELECT * FROM users WHERE id = ?;`,
		id,
	).Scan(
		&result.Id,
		&result.TelegramId,
		&result.Name, &result.Surname,
		&result.Email,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("users.GetUserByTelegramId: failed to query row: %w", err)
	}

	return result, nil
}
