package main

import "testing"

func TestContainsImage(t *testing.T) {
	var imgs []Image
	var expected bool
	var actual bool

	imgs = []Image{}
	expected = false
	actual = ContainsImage("some image", imgs)
	if actual != expected {
		t.Errorf("Expected ContainsImage to be %v, got %v instead", expected, actual)
	}

	imgs = []Image{
		Image{Name: "first"},
		Image{Name: "second"},
		Image{Name: "third"},
	}
	expected = false
	actual = ContainsImage("some image", imgs)
	if actual != expected {
		t.Errorf("Expected ContainsImage to be %v, got %v instead", expected, actual)
	}

	imgs = []Image{
		Image{Name: "first"},
		Image{Name: "second"},
		Image{Name: "third"},
	}
	expected = true
	actual = ContainsImage("third", imgs)
	if actual != expected {
		t.Errorf("Expected ContainsImage to be %v, got %v instead", expected, actual)
	}
}
