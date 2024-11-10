package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

var AuthError = errors.New("Unauthorized")

type loginData struct {
	email          string
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

func register(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return c.String(http.StatusMethodNotAllowed, "Ilegal")
	}

	username := c.FormValue("username")
	password := c.FormValue("password")
	if len(username) < 8 || len(password) < 8 {
		return c.String(http.StatusNotAcceptable, "Usuario/Contraseña debe ser mayor a 8 caracteres")
	}

	if _, ok := users[username]; ok {
		return c.String(http.StatusConflict, "Usuario ya existe")
	}

	hashedPassword, _ := hashPassword(password)
	users[username] = Login{
		HashedPassword: hashedPassword,
	}
	return c.String(http.StatusOK, "Registro exitoso")
}

func login(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return c.String(http.StatusMethodNotAllowed, "Ilegal")
	}

	username := c.FormValue("username")
	password := c.FormValue("password")

	user, ok := users[username]

	if !ok || !checkHashPassword(password, user.HashedPassword) {
		return c.String(http.StatusUnauthorized, "Invalid username or password")
	}

	// assign tokens and export
	sessionToken := generateToken(32)
	csrfToken := generateToken(32) // cross site request forgery

	c.SetCookie(&http.Cookie{
		Name:     "sessionToken",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	c.SetCookie(&http.Cookie{
		Name:     "csrfToken",
		Value:    csrfToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
	})

	// store session token to db
	user.SessionToken = sessionToken
	user.CSRFToken = csrfToken
	users[username] = user

	return c.String(http.StatusOK, "Login successful")
}

func logout(c echo.Context) error {
	if err := authorize(c); err != nil {
		return c.String(http.StatusUnauthorized, "No autorizado")
	}

	// Clear cookies for session and CSRF tokens
	c.SetCookie(&http.Cookie{
		Name:     "sessionToken",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})

	c.SetCookie(&http.Cookie{
		Name:     "csrfToken",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: false,
	})

	username := c.FormValue("username")
	user := users[username]
	user.SessionToken = ""
	user.CSRFToken = ""
	users[username] = user

	return c.String(http.StatusOK, "Logout successful")
}

func protected(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return c.String(http.StatusMethodNotAllowed, "Metodo de solicitud ilegal")
	}

	if err := authorize(c); err != nil {
		return c.String(http.StatusUnauthorized, "No autorizado")
	}

	username := c.FormValue("username")
	return c.String(http.StatusOK, fmt.Sprintf("CSRF validation successful, user: %s", username))
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkHashPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func authorize(c echo.Context) error {
	username := c.FormValue("username")
	user, ok := users[username]
	if !ok {
		return fmt.Errorf("AuthError")
	}

	sessionToken, err := c.Cookie("sessionToken")
	if err != nil || sessionToken.Value == "" || sessionToken.Value != user.SessionToken {
		return fmt.Errorf("AuthError")
	}

	csrf := c.Request().Header.Get("X-CSRF-Token")
	if csrf != user.CSRFToken || csrf == "" {
		return fmt.Errorf("AuthError")
	}

	return nil
}
