package workflow

import (
	"time"

	"tkview/internal/environment"
	"tkview/internal/organisation"
)

// ExecutionID is the unique identifier of a test workflow execution.
type ExecutionID string

// Execution is a tkview representation of a test workflow execution.
type Execution struct {
	ID        ExecutionID
	Name      string
	StartedAt time.Time
	Status    string
}

// ExecutionLister should return test workflow executions from a datasource.
type ExecutionLister interface {
	ListExecutions(orgID organisation.ID, envID environment.ID, id ID) ([]Execution, error)
}
