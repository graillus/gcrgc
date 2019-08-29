package docker

// Image represents a repository image
type Image struct {
	Digest    string
	Tags      []string
	IsRemoved bool
}

// NewImage creates a new Image
func NewImage(digest string, tags []string) *Image {
	return &Image{digest, tags, false}
}

// HasTag checks if the image is tagged with a given tag
func (i Image) HasTag(tag string) bool {
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
