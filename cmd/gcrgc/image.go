package main

// Image represents a docker image inside repository
type Image struct {
	Name string `json:"name"`
}

// NewImage returns a new instance of Image
func NewImage(name string) *Image {
	return &Image{name}
}

// ContainsImage checks if an image is present in an array of Image structs
func ContainsImage(name string, images []Image) bool {
	for _, item := range images {
		if name == item.Name {
			return true
		}
	}

	return false
}
