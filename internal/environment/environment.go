// Package environment provides an abstraction for environment data models.
package environment

import (
	"tkview/internal/organisation"
)

// ID is the unique identifier of an environment.
type ID string

// Environment is a tkview representation of an environment.
type Environment struct {
	ID   ID
	Name string
}

// Lister should return environments from a datasource.
type Lister interface {
	ListEnvironments(organisationID organisation.ID) ([]Environment, error)
}
