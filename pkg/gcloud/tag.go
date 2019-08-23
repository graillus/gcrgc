package gcloud

// Tag represents an Image tag
type Tag struct {
	Digest    string    `json:"digest"`
	Tags      []string  `json:"tags"`
	Timestamp Timestamp `json:"timestamp"`
	IsRemoved bool
}

// Timestamp holds the Tag's date and time information
type Timestamp struct {
	Datetime string `json:"datetime"`
}

// ContainsTag checks if the image is tagged with a given tag
func (t Tag) ContainsTag(tag string) bool {
	for _, i := range t.Tags {
		if i == tag {
			return true
		}
	}

	return false
}

// IsTagged tells if the current Tag has at least one tag
func (t Tag) IsTagged() bool {
	return len(t.Tags) > 0
}
