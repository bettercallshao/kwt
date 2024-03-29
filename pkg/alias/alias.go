package alias

import (
	"unicode"

	"github.com/golang-collections/collections/set"
)

// Store is a set
type Store struct {
	set *set.Set
}

// New creates a new alias store
func New() Store {
	store := Store{set.New()}
	Avoid(store, []string{"h"})
	return store
}

// Avoid given alias
func Avoid(store Store, avoid []string) {
	for _, item := range avoid {
		store.set.Insert(item)
	}
}

// Pick a alias from name and remember in store
func Pick(store Store, name string) []string {
	if len(name) > 1 {
		for _, r := range name {
			if unicode.IsLetter(r) {
				c := string(r)
				if !store.set.Has(c) {
					store.set.Insert(c)
					return []string{c}
				}
			}
		}
	}
	return []string{}
}
