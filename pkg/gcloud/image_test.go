package gcloud

import "testing"

var containsTagTest = []struct {
	img      Image
	test     string
	expected bool
}{
	{Image{}, "some tag", false},
	{Image{Tags: []string{"not matching tag"}}, "some tag", false},
	{Image{Tags: []string{"not matching tag", "matching tag"}}, "matching tag", true},
}

func TestContainsTag(t *testing.T) {
	var actual bool

	for _, test := range containsTagTest {
		actual = test.img.ContainsTag(test.test)
		if actual != test.expected {
			t.Errorf("Expected ContainsTag to be %v, got %v instead", test.expected, actual)
		}
	}
}

var isTaggedTest = []struct {
	img      Image
	expected bool
}{
	{Image{}, false},
	{Image{Tags: []string{}}, false},
	{Image{Tags: []string{"a tag"}}, true},
}

func TestIsTagged(t *testing.T) {
	var actual bool

	for _, test := range isTaggedTest {
		actual = test.img.IsTagged()
		if actual != test.expected {
			t.Errorf("Expected IsTagged to be %v, got %v instead", test.expected, actual)
		}
	}
}
