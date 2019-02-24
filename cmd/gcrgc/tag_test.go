package main

import "testing"

var containsTagTest = []struct {
	tag      Tag
	test     string
	expected bool
}{
	{Tag{}, "some tag", false},
	{Tag{Tags: []string{"not matching tag"}}, "some tag", false},
	{Tag{Tags: []string{"not matching tag", "matching tag"}}, "matching tag", true},
}

func TestContainsTag(t *testing.T) {
	var actual bool

	for _, test := range containsTagTest {
		actual = test.tag.ContainsTag(test.test)
		if actual != test.expected {
			t.Errorf("Expected ContainsTag to be %v, got %v instead", test.expected, actual)
		}
	}
}

var isTaggedTest = []struct {
	tag      Tag
	expected bool
}{
	{Tag{}, false},
	{Tag{Tags: []string{}}, false},
	{Tag{Tags: []string{"a tag"}}, true},
}

func TestIsTagged(t *testing.T) {
	var actual bool

	for _, test := range isTaggedTest {
		actual = test.tag.IsTagged()
		if actual != test.expected {
			t.Errorf("Expected IsTagged to be %v, got %v instead", test.expected, actual)
		}
	}
}
