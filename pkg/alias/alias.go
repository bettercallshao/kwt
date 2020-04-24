package alias

import (
	"github.com/golang-collections/collections/set"
)

// Store is a set
type Store struct {
	set *set.Set
}

// New creates a new alias store
func New() Store {
	return Store{set.New()}
}

// Avoid given alias
func Avoid(store Store, avoid []string) {
	for _, item := range avoid {
		store.set.Insert(item)
	}
}

// Pick a alias from name and remember in store
func Pick(store Store, name string) []string {
	for _, r := range name {
		c := string(r)
		if !store.set.Has(c) {
			store.set.Insert(c)
			return []string{c}
		}
	}
	return []string{}
}
