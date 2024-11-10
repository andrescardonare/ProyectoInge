package controllers

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/microsoft/go-mssqldb"
	"log"
	"os"
	"time"
)

var db *sql.DB

type User struct {
	ID             int       `gorm:"primaryKey;autoIncrement"`
	Username       string    `gorm:"size:50;unique;not null"`
	HashedPassword string    `gorm:"type:nvarchar(MAX);not null"`
	SessionToken   *string   `gorm:"type:nvarchar(MAX)"` // nullable
	CSRFToken      *string   `gorm:"type:nvarchar(MAX)"` // nullable
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func DBconnection() {
	var server = os.Getenv("DB_SERVER")
	var port = 1433
	var user = os.Getenv("DB_USER")
	var password = os.Getenv("DB_PASSWORD")
	var database = os.Getenv("DB_DATABASE")

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)
	var err error

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("DB Connection Established!")
}

func createUser(username, hashedPassword string) (int, error) {
	query := `
        INSERT INTO users (username, hashed_password, created_at, updated_at)
        OUTPUT INSERTED.id
        VALUES (@p1, @p2, @p3, @p4);
    `
	var userID int
	err := db.QueryRow(query, username, hashedPassword, time.Now(), time.Now()).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}
	return userID, nil
}

// Retrieve a user by username
func getUserByUsername(username string) (*User, error) {
	query := `
        SELECT id, username, hashed_password, session_token, csrf_token, created_at, updated_at
        FROM users
        WHERE username = @p1;
    `
	row := db.QueryRow(query, username)
	user := User{}
	err := row.Scan(
		&user.ID, &user.Username, &user.HashedPassword,
		&user.SessionToken, &user.CSRFToken,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}
	return &user, nil
}

// Update a user's session and CSRF tokens
func updateUserTokens(userID int, sessionToken, csrfToken string) error {
	query := `
        UPDATE users
        SET session_token = @p1, csrf_token = @p2, updated_at = @p3
        WHERE id = @p4;
    `
	_, err := db.Exec(query, sessionToken, csrfToken, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update tokens: %w", err)
	}
	return nil
}

// Delete a user by ID
func deleteUser(userID int) error {
	query := `DELETE FROM users WHERE id = @p1;`
	_, err := db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
