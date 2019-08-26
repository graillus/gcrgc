package gcloud

import "testing"

func TestNewRegistry(t *testing.T) {
	reg := NewRegistry("name")

	expected := "name"
	actual := reg.Name
	if actual != expected {
		t.Errorf("Expected Registry.Name to be %s, got %s", expected, actual)
	}

	if len(reg.Repositories) != 0 {
		t.Errorf("Expected Registry.Repositories to be empty, got %d element", len(reg.Repositories))
	}
}

var containsTests = []struct {
	repos    []Repository
	test     string
	expected bool
}{
	{
		[]Repository{},
		"some image",
		false,
	},
	{
		[]Repository{
			Repository{Name: "first"},
			Repository{Name: "second"},
			Repository{Name: "third"},
		},
		"some image",
		false,
	},
	{
		[]Repository{
			Repository{Name: "first"},
			Repository{Name: "second"},
			Repository{Name: "third"},
		},
		"third",
		true,
	},
}

func TestContainsRepository(t *testing.T) {
	var actual bool

	for _, test := range containsTests {
		registry := Registry{"", test.repos}
		actual = registry.ContainsRepository(test.test)
		if actual != test.expected {
			t.Errorf("Expected ContainsRepository to be %v, got %v instead", test.expected, actual)
		}
	}
}
