// Package workflow provides an abstraction for workflow data models.
package workflow

import (
	"time"

	"tkview/internal/environment"
	"tkview/internal/organisation"
)

// ID is the unique identifier of a workflow.
type ID string

// Workflow is a tkview representation of a workflow.
type Workflow struct {
	ID                  ID
	Name                string
	LastExecutionAt     time.Time
	LastExecutionStatus string
}

// Lister should return workflows from a datasource.
type Lister interface {
	ListWorkflows(orgID organisation.ID, envID environment.ID) ([]Workflow, error)
}
