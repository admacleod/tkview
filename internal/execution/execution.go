// Package execution provides an abstraction for execution data models.
package execution

import (
	"tkview/internal/environment"
	"tkview/internal/organisation"
)

// ID is the unique identifier of a test workflow execution.
type ID string

// Execution is a tkview representation of a test workflow execution.
type Execution struct {
	ID   ID
	Name string
}

// Lister should return test workflow executions from a datasource.
type Lister interface {
	ListExecutions(organisationID organisation.ID, environmentID environment.ID) ([]Execution, error)
}
