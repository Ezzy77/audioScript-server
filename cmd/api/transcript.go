package main

import (
	"github.com/Ezzy77/audioScript-server/internal"
	"github.com/gin-gonic/gin"
)

func (app *application) createTranscriptJobHandler(ctx *gin.Context) {
	// Parse request body
	result := internal.Job(app.awsConfig.AWS.AccessKey, app.awsConfig.AWS.SecretKey, app.awsConfig.AWS.Region)

	ctx.JSON(200, gin.H{
		"result": result,
	})
}
