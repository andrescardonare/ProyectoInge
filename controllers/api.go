package controllers

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

func PostToApi(c echo.Context) error {
	// Prepare form data for POST request
	formData := []byte("origin=value1&destination=value2&fecha=value3") // Modify with your actual form data
	req, err := http.NewRequest("POST", "localhost:8000/predict/", bytes.NewBuffer(formData))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error creating request")
	}

	// Set necessary headers
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error sending request")
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error reading response")
	}

	// Process the HTML response as a string
	htmlResponse := string(body)

	// Send the HTML response back to the client
	return c.HTML(http.StatusOK, htmlResponse)
}
