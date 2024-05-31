package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Enabling the release mode
	gin.SetMode(gin.ReleaseMode)

	// Setting the port
	port := 8888
	formattedPort := fmt.Sprintf(":%d", port)

	router := gin.Default()

	// Struct to represent the response
	type ImgURL struct {
		URL string `json:"ImgURL"`
	}

	// POST endpoint to handle image uploads
	router.POST("/upload", func(c *gin.Context) {
		// Retrieve the file from the form
		file, err := c.FormFile("img")
		if err != nil {
			errLog := fmt.Sprintf("Bad request. Error: %s", err)
			c.String(http.StatusBadRequest, errLog)
			return
		}

		// Define the destination path
		dst := filepath.Join("images", file.Filename)

		// Check if the file already exists
		if _, err := os.Stat(dst); err == nil {
			// File exists, generate a new file name
			timestamp := time.Now().Unix()
			ext := filepath.Ext(file.Filename)
			name := file.Filename[:len(file.Filename)-len(ext)]
			newFilename := fmt.Sprintf("%s_%d%s", name, timestamp, ext)
			dst = filepath.Join("images", newFilename)
		}

		// Save the file to the destination
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Upload failed: %s", err.Error()))
			return
		}

		// Create the response
		url := fmt.Sprintf("http://localhost%s/images/%s", formattedPort, filepath.Base(dst))
		response := ImgURL{
			URL: url,
		}

		c.JSON(http.StatusOK, response)
	})

	// Serve uploaded files
	router.Static("/images", "./images")

	// Start the server
	router.Run(formattedPort)
}
