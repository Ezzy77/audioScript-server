package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Ezzy77/audioScript-server/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/transcribe"
	"github.com/aws/aws-sdk-go-v2/service/transcribe/types"
	"io"
	"log"
	"net/url"
	"path"
	"time"
)

func InitAwsServices(awsAccessKey,
	awsSecretKey,
	awsRegion string,
) (*s3.Client, *transcribe.Client, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKey, awsSecretKey, "")),
		config.WithRegion(awsRegion),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load configuration: %v", err)
	}
	// Create an S3 client
	s3Client := s3.NewFromConfig(cfg)

	// Create a Transcribe client
	transcribeClient := transcribe.NewFromConfig(cfg)

	return s3Client, transcribeClient, nil
}

func StartJob(
	transcribeClient *transcribe.Client,
	s3Client *s3.Client, awsRegion,
	transcriptionJobName,
	mediaFileUri,
	outputBucketName,
	languageCode string,
) (model.Transcription, error) {
	transcript, err := processTranscriptionJob(
		transcribeClient,
		s3Client,
		transcriptionJobName,
		mediaFileUri,
		outputBucketName,
		languageCode)
	if err != nil {
		return model.Transcription{}, fmt.Errorf("failed to process transcription job: %v", err)
	}
	return transcript, nil
}

func processTranscriptionJob(transcribeClient *transcribe.Client, s3Client *s3.Client, transcriptionJobName, mediaFileUri, outputBucketName, languageCode string) (model.Transcription, error) {
	return processTranscriptionJobImpl(transcribeClient, s3Client, transcriptionJobName, mediaFileUri, outputBucketName, languageCode)
}

func processTranscriptionJobImpl(transcribeClient *transcribe.Client, s3Client *s3.Client, transcriptionJobName, mediaFileUri, outputBucketName, languageCode string) (model.Transcription, error) {
	// Create a context with a timeout that will be canceled manually once we're finished with it.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	// Start the transcription job
	_, err := transcribeClient.StartTranscriptionJob(ctx, &transcribe.StartTranscriptionJobInput{
		TranscriptionJobName: aws.String(transcriptionJobName),
		LanguageCode:         types.LanguageCode.Values(types.LanguageCode(languageCode))[0],
		Media:                &types.Media{MediaFileUri: aws.String(mediaFileUri)},
		OutputBucketName:     aws.String(outputBucketName),
	})

	if err != nil {
		return model.Transcription{}, fmt.Errorf("failed to start transcription job: %v", err)
	}

	// Utilize the AWS SDK's built-in retry functionality to monitor the job's status.
	for {
		// This select statement helps us respect the context's deadline
		select {
		case <-ctx.Done():
			return model.Transcription{}, ctx.Err()
		default:
			// If the context hasn't been canceled or exceeded its deadline, continue checking the transcription job status
			output, err := transcribeClient.GetTranscriptionJob(ctx, &transcribe.GetTranscriptionJobInput{
				TranscriptionJobName: aws.String(transcriptionJobName),
			})

			if err != nil {
				return model.Transcription{}, fmt.Errorf("failed to get transcription job: %v", err)
			}

			// Handle the job status accordingly
			switch output.TranscriptionJob.TranscriptionJobStatus {
			case types.TranscriptionJobStatusCompleted:
				transcriptFileUri := *output.TranscriptionJob.Transcript.TranscriptFileUri
				parsedUri, err := url.Parse(transcriptFileUri)
				if err != nil {
					return model.Transcription{}, fmt.Errorf("failed to parse transcript file URI: %v", err)
				}

				transcriptFileName := path.Base(parsedUri.Path)

				getObjectOutput, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
					Bucket: aws.String(outputBucketName),
					Key:    aws.String(transcriptFileName),
				})

				if err != nil {
					return model.Transcription{}, fmt.Errorf("failed to get object: %v", err)
				}

				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
						log.Println(err)
					}
				}(getObjectOutput.Body)

				var transcript model.Transcription

				if err := json.NewDecoder(getObjectOutput.Body).Decode(&transcript); err != nil {
					return model.Transcription{}, fmt.Errorf("failed to decode JSON: %v", err)
				}

				return transcript, nil

			case types.TranscriptionJobStatusFailed:
				return model.Transcription{}, fmt.Errorf("transcription job failed")

			case types.TranscriptionJobStatusInProgress:
				fmt.Println("Transcription progress...")
				time.Sleep(5 * time.Second)

			default:
				return model.Transcription{}, fmt.Errorf("unknown transcription job status: %v", output.TranscriptionJob.TranscriptionJobStatus)
			}
		}
	}
}
