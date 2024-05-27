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
	"log"
	"net/url"
	"path"
	"time"
)

func Job(awsAccessKey, awsSecretKey, awsRegion string) model.Transcription {

	// load shared aws configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKey, awsSecretKey, "")),
		config.WithRegion(awsRegion),
	)
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	// Create an S3 client
	s3Client := s3.NewFromConfig(cfg)
	transcribeClient := transcribe.NewFromConfig(cfg)

	transcriptionJobName := "myTranscriptionJobElisio" + time.Now().Format("20060102150405")

	_, err = transcribeClient.StartTranscriptionJob(context.TODO(), &transcribe.StartTranscriptionJobInput{
		Media: &types.Media{
			MediaFileUri: aws.String("s3://myelisiobucket/tictactoe_trimmed.mp4"),
		},
		TranscriptionJobName: aws.String(transcriptionJobName),
		OutputBucketName:     aws.String("output99"),
		LanguageCode:         types.LanguageCode.Values("en-US")[0],
	})

	if err != nil {
		log.Fatalf("failed to start transcription job, %v", err)
	}

	for {
		output, err := transcribeClient.GetTranscriptionJob(context.TODO(), &transcribe.GetTranscriptionJobInput{
			TranscriptionJobName: aws.String(transcriptionJobName),
		})

		if err != nil {
			log.Fatalf("failed to get transcription job, %v", err)
		}

		switch output.TranscriptionJob.TranscriptionJobStatus {
		case types.TranscriptionJobStatusCompleted:
			// The job is completed, print the transcript
			//fmt.Println(output.TranscriptionJob.Transcript.TranscriptFileUri)

			transcriptFileUri := *output.TranscriptionJob.Transcript.TranscriptFileUri
			parsedUri, err := url.Parse(transcriptFileUri)
			if err != nil {
				log.Fatalf("failed to parse transcript file URI, %v", err)
			}

			transcriptFileName := path.Base(parsedUri.Path)

			// Get the object
			getObjectOutput, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
				Bucket: aws.String("output99"), // replace with your bucket name
				Key:    aws.String(transcriptFileName),
			})

			if err != nil {
				log.Fatalf("failed to get object, %v", err)
			}

			// Decode the JSON content of the object
			var transcript model.Transcription

			err = json.NewDecoder(getObjectOutput.Body).Decode(&transcript)
			if err != nil {
				log.Fatalf("failed to decode JSON, %v", err)
			}
			return transcript
		case types.TranscriptionJobStatusFailed:
			log.Fatalf("transcription job failed")
		case types.TranscriptionJobStatusInProgress:
			// The job is still in progress, wait for a while before checking the status again
			fmt.Println("Transcription job is still in progress")
			time.Sleep(5 * time.Second)
		default:
			log.Fatalf("unknown transcription job status: %v", output.TranscriptionJob.TranscriptionJobStatus)
		}
	}

}

func InitAwsServices(awsAccessKey, awsSecretKey, awsRegion string) (*s3.Client, *transcribe.Client, error) {
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

func StartJob(transcribeClient *transcribe.Client, s3Client *s3.Client, awsRegion, transcriptionJobName, mediaFileUri, outputBucketName, languageCode string) (model.Transcription, error) {
	transcript, err := processTranscriptionJob(transcribeClient, s3Client, transcriptionJobName, mediaFileUri, outputBucketName, languageCode)
	if err != nil {
		return model.Transcription{}, fmt.Errorf("failed to process transcription job: %v", err)
	}
	return transcript, nil
}

func processTranscriptionJob(transcribeClient *transcribe.Client, s3Client *s3.Client, transcriptionJobName, mediaFileUri, outputBucketName, languageCode string) (model.Transcription, error) {

	_, err := transcribeClient.StartTranscriptionJob(context.TODO(), &transcribe.StartTranscriptionJobInput{
		Media: &types.Media{
			MediaFileUri: aws.String(mediaFileUri),
		},
		TranscriptionJobName: aws.String(transcriptionJobName),
		OutputBucketName:     aws.String(outputBucketName),
		LanguageCode:         types.LanguageCode.Values(types.LanguageCode(languageCode))[0],
	})

	if err != nil {
		return model.Transcription{}, fmt.Errorf("failed to start transcription job: %v", err)
	}

	for {
		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
		defer cancel()

		output, err := transcribeClient.GetTranscriptionJob(ctx, &transcribe.GetTranscriptionJobInput{
			TranscriptionJobName: aws.String(transcriptionJobName),
		})

		if err != nil {
			return model.Transcription{}, fmt.Errorf("failed to get transcription job: %v", err)
		}

		switch output.TranscriptionJob.TranscriptionJobStatus {
		case types.TranscriptionJobStatusCompleted:
			transcriptFileUri := *output.TranscriptionJob.Transcript.TranscriptFileUri
			parsedUri, err := url.Parse(transcriptFileUri)
			if err != nil {
				return model.Transcription{}, fmt.Errorf("failed to parse transcript file URI: %v", err)
			}

			transcriptFileName := path.Base(parsedUri.Path)

			getObjectOutput, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
				Bucket: aws.String(outputBucketName),
				Key:    aws.String(transcriptFileName),
			})

			if err != nil {
				return model.Transcription{}, fmt.Errorf("failed to get object: %v", err)
			}

			defer getObjectOutput.Body.Close()

			var transcript model.Transcription

			if err := json.NewDecoder(getObjectOutput.Body).Decode(&transcript); err != nil {
				return model.Transcription{}, fmt.Errorf("failed to decode JSON: %v", err)
			}

			return transcript, nil

		case types.TranscriptionJobStatusFailed:
			return model.Transcription{}, fmt.Errorf("transcription job failed")

		case types.TranscriptionJobStatusInProgress:
			time.Sleep(5 * time.Second)

		default:
			return model.Transcription{}, fmt.Errorf("unknown transcription job status: %v", output.TranscriptionJob.TranscriptionJobStatus)
		}
	}
}
