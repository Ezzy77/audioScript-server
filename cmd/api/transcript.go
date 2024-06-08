package main

import (
	"fmt"
	"github.com/Ezzy77/audioScript-server/internal"
	"github.com/gin-gonic/gin"
	"time"
)

func (app *application) createTranscriptJobHandler(ctx *gin.Context) {
	//result := internal.Job(app.awsConfig.AWS.AccessKey, app.awsConfig.AWS.SecretKey, app.awsConfig.AWS.Region)
	jobName := "myTranscriptionJobElisio" + time.Now().Format("20060102150405")
	MediaFileUri := app.awsConfig.AWS.MediaFileUri
	outputBucketName := app.awsConfig.AWS.OutputBucket
	languageCode := "en-US"
	s3Client, transcribeService, err := internal.InitAwsServices(
		app.awsConfig.AWS.AccessKey,
		app.awsConfig.AWS.SecretKey,
		app.awsConfig.AWS.Region,
	)
	if err != nil {

	}
	job, err := internal.StartJob(
		transcribeService,
		s3Client,
		app.awsConfig.AWS.Region,
		jobName,
		MediaFileUri,
		outputBucketName,
		languageCode,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx.JSON(200, gin.H{
		"result": job,
	})
}
