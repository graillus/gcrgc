package main

import "testing"

func TestContainsTag(t *testing.T) {
	var tag Tag
	var expected bool
	var actual bool

	tag = Tag{}
	expected = false
	actual = tag.ContainsTag("some tag")
	if actual != expected {
		t.Errorf("Expected ContainsTag to be %v, got %v instead", expected, actual)
	}

	tag = Tag{Tags: []string{"not matching tag"}}
	expected = false
	actual = tag.ContainsTag("matching tag")
	if actual != expected {
		t.Errorf("Expected ContainsTag to be %v, got %v instead", expected, actual)
	}

	tag = Tag{Tags: []string{"not matching tag", "matching tag"}}
	expected = true
	actual = tag.ContainsTag("matching tag")
	if actual != expected {
		t.Errorf("Expected ContainsTag to be %v, got %v instead", expected, actual)
	}
}
func TestIsTagged(t *testing.T) {
	var tag Tag
	var expected bool
	var actual bool

	tag = Tag{}
	expected = false
	actual = tag.IsTagged()
	if actual != expected {
		t.Errorf("Expected IsTagged to be %v, got %v instead", expected, actual)
	}

	tag = Tag{Tags: []string{"a tag"}}
	expected = true
	actual = tag.IsTagged()
	if actual != expected {
		t.Errorf("Expected IsTagged to be %v, got %v instead", expected, actual)
	}
}
