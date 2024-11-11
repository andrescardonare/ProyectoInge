package controllers

import (
	"errors"
	"fmt"
	_ "github.com/microsoft/go-mssqldb"
	"gorm.io/driver/sqlserver"
	_ "gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"log"
	"os"
)

var db *gorm.DB //*sql.DB

type user_struct struct {
	gorm.Model
	Username       string `gorm:"size:50;unique;not null"`
	HashedPassword string `gorm:"type:nvarchar(MAX);not null"`
	SessionToken   string `gorm:"type:nvarchar(MAX)"`
	CSRFToken      string `gorm:"type:nvarchar(MAX)"`
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
	db, err = gorm.Open(sqlserver.Open(connString), &gorm.Config{})
	if err != nil {
		fmt.Println("errordbconn")
		log.Fatal("Error creating connection pool: ", err.Error())
	}

	db.AutoMigrate(&user_struct{})

	fmt.Println("DB Connection Established!")
}

func createUserDB(username, hashedPassword string) error {
	var usr user_struct
	result := db.Where("username = ?", username).First(&usr)
	if result.Error == nil {
		// Usuario ya existe
		fmt.Println("User already exists")
		return errors.New("user already exists")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Otro tipo de error ocurrió
		return fmt.Errorf("error retrieving user: %w", result.Error)
	}

	// Crear el nuevo usuario
	if err := db.Create(&user_struct{Username: username, HashedPassword: hashedPassword}).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	fmt.Println("User created successfully")
	return nil
}

func updateTokensDB(username, newSessionToken, newCSRFToken string) error {
	// Attempt to update only the specified fields for the user
	result := db.Model(&user_struct{}).Where("username = ?", username).Updates(map[string]interface{}{
		"SessionToken": newSessionToken,
		"CSRFToken":    newCSRFToken,
	})

	// Check if an error occurred during the update
	if result.Error != nil {
		return fmt.Errorf("failed to update tokens: %w", result.Error)
	}

	// Ensure at least one record was updated
	if result.RowsAffected == 0 {
		return errors.New("no user found with the specified username")
	}

	fmt.Println("Tokens updated successfully")
	return nil
}

// Retrieve a user by username
func getUserByUsername(username string) (*user_struct, error) {
	var user user_struct

	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		fmt.Println("result.Error")
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// No user found with this username
			fmt.Println("result.Error2")
			return nil, errors.New("user not found")
		}
		// Some other error occurred
		return nil, result.Error
	}
	// User found
	fmt.Println("user found")
	return &user, nil
}

/*
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
*/
