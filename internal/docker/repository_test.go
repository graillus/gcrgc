package docker

import "testing"

func TestNewRepository(t *testing.T) {
	img := NewRepository("name")

	expected := "name"
	actual := img.Name
	if actual != expected {
		t.Errorf("Expected Repository.Name to be %s, got %s", expected, actual)
	}
}
