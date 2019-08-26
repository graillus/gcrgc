package gcloud

// Image represents a repository image
type Image struct {
	Digest    string    `json:"digest"`
	Tags      []string  `json:"tags"`
	Timestamp Timestamp `json:"timestamp"`
	IsRemoved bool
}

// Timestamp holds the image's date and time information
type Timestamp struct {
	Datetime string `json:"datetime"`
}

// ContainsTag checks if the image is tagged with a given tag
func (i Image) ContainsTag(tag string) bool {
	for _, t := range i.Tags {
		if t == tag {
			return true
		}
	}

	return false
}

// IsTagged tells if the image has at least one tag
func (i Image) IsTagged() bool {
	return len(i.Tags) > 0
}
