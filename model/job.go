package model

type Job struct {
	MediaFileURI string `json:"media_file_uri"`
	OutputBucket string `json:"output_bucket"`
	LanguageCode string `json:"language_code"`
}
