package sync

import (
	"strings"

	"github.com/pkg/errors"
)

// TODO: Document that this function ignores the Old value for changes.
func Apply(base *Entry, changes []*Change) (*Entry, error) {
	// Create a mutable copy of base.
	result := base.copy()

	// Apply changes.
	for _, c := range changes {
		// Handle the special case of a root path.
		if c.Path == "" {
			result = c.New
			continue
		}

		// Crawl down the tree until there is only one component remaining - the
		// parent of the target location.
		parent := result
		components := strings.Split(c.Path, "/")
		for len(components) > 1 {
			child, ok := parent.Find(components[0])
			if !ok {
				return nil, errors.New("unable to resolve parent path")
			}
			parent = child
			components = components[1:]
		}

		// Depending on the new value, either set or remove the entry.
		if c.New == nil {
			if !parent.Remove(components[0]) {
				return nil, errors.New("unable to resolve path for deletion")
			}
		} else {
			parent.Insert(components[0], c.New)
		}
	}

	// Done.
	return result, nil
}
