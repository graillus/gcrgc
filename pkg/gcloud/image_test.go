package gcloud

import "testing"

func TestNewImage(t *testing.T) {
	img := NewImage("name")

	expected := "name"
	actual := img.Name
	if actual != expected {
		t.Errorf("Expected Image.Name to be %s, got %s", expected, actual)
	}
}

var containsTests = []struct {
	imgs     []Image
	test     string
	expected bool
}{
	{
		[]Image{},
		"some image",
		false,
	},
	{
		[]Image{
			Image{Name: "first"},
			Image{Name: "second"},
			Image{Name: "third"},
		},
		"some image",
		false,
	},
	{
		[]Image{
			Image{Name: "first"},
			Image{Name: "second"},
			Image{Name: "third"},
		},
		"third",
		true,
	},
}

func TestContainsImage(t *testing.T) {
	var actual bool

	for _, test := range containsTests {
		actual = ContainsImage(test.test, test.imgs)
		if actual != test.expected {
			t.Errorf("Expected ContainsImage to be %v, got %v instead", test.expected, actual)
		}
	}
}
