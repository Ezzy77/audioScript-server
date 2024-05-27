package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func (app *application) createTranscriptJobHandler(ctx *gin.Context) {
	// Parse request body
	file, err := ctx.FormFile("file")
	if err != nil {
		return
	}
	// Print the file type, size, and name
	fmt.Printf("File type: %v\n", file.Header.Get("Content-Type"))
	fmt.Printf("File size: %v\n", file.Size)
	fmt.Printf("File name: %v\n", file.Filename)

	ctx.JSON(200, gin.H{
		"message": "File received successfully",
	})
}
