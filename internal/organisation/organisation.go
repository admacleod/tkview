// Package organisation provides an abstraction for organisation data models.
package organisation

// ID is the unique identifier of an organisation.
type ID string

// Organisation is a tkview representation of an organisation.
type Organisation struct {
	ID   ID
	Name string
}

// Lister should return organisations from a datasource.
type Lister interface {
	ListOrganisations() ([]Organisation, error)
}
