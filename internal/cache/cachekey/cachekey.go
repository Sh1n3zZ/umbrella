// Package cachekey centralises Redis key construction so callers never
// concatenate raw strings. Keys follow the "<business>:<entity>:<id>" layout.
package cachekey

import "fmt"

// Namespace declares a (business, entity) pair used to build typed Redis keys.
type Namespace struct {
	business string
	entity   string
}

// NewNamespace builds a Namespace from its business and entity segments.
// Both segments must be non-empty lowercase identifiers in the project's
// "<business>:<entity>:<id>" convention.
func NewNamespace(business, entity string) Namespace {
	return Namespace{business: business, entity: entity}
}

// Key returns "<business>:<entity>:<id>" for the given id.
func (n Namespace) Key(id string) string {
	return fmt.Sprintf("%s:%s:%s", n.business, n.entity, id)
}

// Prefix returns "<business>:<entity>:" for SCAN / DEL-by-pattern operations.
func (n Namespace) Prefix() string {
	return fmt.Sprintf("%s:%s:", n.business, n.entity)
}
