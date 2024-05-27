package model

type Transcription struct {
	AccountId string `json:"accountId"`
	JobName   string `json:"jobName"`
	Results   struct {
		Items []struct {
			Alternatives []struct {
				Confidence string `json:"confidence"`
				Content    string `json:"content"`
			} `json:"alternatives"`
			EndTime   string `json:"end_time"`
			StartTime string `json:"start_time"`
			Type      string `json:"type"`
		} `json:"items"`
		Transcripts []struct {
			Transcript string `json:"transcript"`
		} `json:"transcripts"`
	} `json:"results"`
}
