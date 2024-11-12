package controllers

import (
	"ProyectoInge/templates"
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

func Register(c echo.Context) error {
	// verifies method
	if c.Request().Method != http.MethodPost {
		return c.String(http.StatusMethodNotAllowed, "Ilegal")
	}

	// get forms
	username := c.FormValue("username")
	password := c.FormValue("password")

	if len(username) < 8 || len(password) < 8 {
		return c.String(http.StatusNotAcceptable, "Usuario/Contraseña debe ser mayor a 8 caracteres")
	}

	// checks if user exists
	if usr, _ := getUserByUsername(username); usr != nil {
		return c.String(http.StatusOK, "Usuario ya existe!")
	}

	hashed, _ := hashPassword(password)
	// registers user
	err := createUserDB(username, hashed)
	if err != nil {
		return c.String(http.StatusOK, "Usuario ya existe!")
	}

	return c.String(http.StatusOK, "Registro exitoso")
}

func Login(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return c.String(http.StatusMethodNotAllowed, "Ilegal")
	}

	username := c.FormValue("username")
	password := c.FormValue("password")

	// demo mode
	if username == "demo mode" && password == "123456789" {
		return c.String(http.StatusOK, "login successful")
	}

	user, err := getUserByUsername(username)

	if err != nil || !checkHashPassword(password, user.HashedPassword) {
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
	_ = updateTokensDB(username, sessionToken, csrfToken)

	return c.String(http.StatusOK, "login successful")
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

	return c.String(http.StatusOK, "Logout successful")
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionToken, err := c.Cookie("sessionToken")
		if err != nil || sessionToken.Value == "" {
			return c.String(http.StatusUnauthorized, "No autorizado")
		}

		// Aquí puedes agregar lógica adicional para verificar el token en la base de datos
		/*
				user, err := getUserBySessionToken(sessionToken.Value)
				if err != nil || user == nil {
					return c.String(http.StatusUnauthorized, "No autorizado")
				}


			// Almacena información del usuario en el contexto para uso posterior
			c.Set("user", user)
		*/
		return next(c)
	}
}

func Protected(c echo.Context) error {
	/*
		if c.Request().Method != http.MethodPost {
			return c.String(http.StatusMethodNotAllowed, "Metodo de solicitud ilegal")
		}
	*/

	if err := authorize(c); err != nil {
		return c.String(http.StatusUnauthorized, "No autorizado")
	}

	return templates.Protected().Render(context.Background(), c.Response().Writer)
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
	return nil
	/*
		username := c.FormValue("username")
		hash, _ := hashPassword(c.FormValue("password"))
		user, err := getUserByUsername(username)
		if err != nil {
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
	*/

}
