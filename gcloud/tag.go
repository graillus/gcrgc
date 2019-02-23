package gcloud

// Tag represents an Image tag
type Tag struct {
	Digest    string    `json:"digest"`
	Tags      []string  `json:"tags"`
	Timestamp Timestamp `json:"timestamp"`
}

// Timestamp holds the Tag's date and time information
type Timestamp struct {
	Datetime string `json:"datetime"`
}
