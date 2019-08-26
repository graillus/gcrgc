package gcrgc

// Registry
type Registry interface {
	ContainsRepository(name string) bool
}
