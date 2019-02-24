package main

import (
	"testing"
)

func TestTotal(t *testing.T) {
	var report *report
	var expected int
	var actual int

	report = newReport()
	expected = 0
	actual = report.Total()
	if actual != expected {
		t.Errorf("Expected Total to be %d, got %d instead", expected, actual)
	}
}

func TestTotalDeleted(t *testing.T) {
	var report *report
	var expected int
	var actual int

	report = newReport()
	expected = 0
	actual = report.TotalDeleted()
	if actual != expected {
		t.Errorf("Expected TotalDeleted to be %d, got %d instead", expected, actual)
	}
}

func TestReportTag(t *testing.T) {
	var report *report
	var tag Tag
	var expectedTotal int
	var expectedTotalDeleted int
	var actualTotal int
	var actualTotalDeleted int

	report = newReport()
	tag = Tag{}
	report.reportTag(tag)
	expectedTotal = 1
	actualTotal = report.Total()
	if actualTotal != expectedTotal {
		t.Errorf("Expected Total to be %d, got %d instead", expectedTotal, actualTotal)
	}

	expectedTotalDeleted = 0
	actualTotalDeleted = report.TotalDeleted()
	if actualTotal != expectedTotal {
		t.Errorf("Expected TotalDeleted to be %d, got %d instead", expectedTotalDeleted, actualTotalDeleted)
	}

	report = newReport()
	tag = Tag{}
	tag.IsRemoved = true
	report.reportTag(tag)
	expectedTotal = 1
	actualTotal = report.Total()
	if actualTotal != expectedTotal {
		t.Errorf("Expected Total to be %d, got %d instead", expectedTotal, actualTotal)
	}

	expectedTotalDeleted = 1
	actualTotalDeleted = report.TotalDeleted()
	if actualTotal != expectedTotal {
		t.Errorf("Expected TotalDeleted to be %d, got %d instead", expectedTotalDeleted, actualTotalDeleted)
	}
}
