package model

type Job struct {
	S3Uri        string `json:"s3_uri"`
	LanguageCode string `json:"language_code"`
}
