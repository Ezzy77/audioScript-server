package main

import (
	"context"
	"fmt"
	"github.com/Ezzy77/audioScript-server/internal"
	"github.com/Ezzy77/audioScript-server/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func (app *application) createTranscriptJobHandler(ctx *gin.Context) {

	jobName := "transcriptJob" + time.Now().Format("20060102150405")

	var job model.Job
	err := ctx.ShouldBindJSON(&job)
	fmt.Println("job is here...", job)
	if err != nil {
		fmt.Println("Error binding request body: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transcript, err := internal.StartJob(
		app.transcribeClient,
		app.s3Client,
		app.awsConfig.Region,
		app.awsConfig.OutputBucket,
		jobName,
		job.S3Uri,
		job.LanguageCode,
	)
	if err != nil {
		fmt.Println("Error starting job: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"result": transcript,
	})
}

func (app *application) uploadMedia(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//defer func(src multipart.File) {
	//	err := src.Close()
	//	if err != nil {
	//
	//	}
	//}(src)

	// Create a unique file name
	originalName := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))
	fileName := fmt.Sprintf("%s_%d%s", originalName, time.Now().UnixNano(), path.Ext(file.Filename))

	// Upload the file to S3
	_, err = app.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(app.awsConfig.UploadBucket),
		Key:    aws.String(fileName),
		Body:   src,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate the presigned URL for the uploaded file
	presignClient := s3.NewPresignClient(app.s3Client)
	input := &s3.GetObjectInput{
		Bucket: aws.String(app.awsConfig.OutputBucket),
		Key:    aws.String(fileName),
	}
	presignedURL, err := presignClient.PresignGetObject(context.Background(), input, func(o *s3.PresignOptions) {
		o.Expires = 15 * time.Minute
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Generate S3 URI
	s3Uri := fmt.Sprintf("s3://%s/%s", app.awsConfig.UploadBucket, fileName)

	ctx.JSON(http.StatusOK, gin.H{
		"url":    presignedURL,
		"s3_uri": s3Uri,
	})
}
