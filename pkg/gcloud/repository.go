package gcloud

// Repository represents a docker image inside repository
type Repository struct {
	Name   string `json:"name"`
	Images []Image
}

// NewRepository returns a new instance of Repository
func NewRepository(name string) *Repository {
	return &Repository{name, []Image{}}
}
