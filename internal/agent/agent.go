// Package agent provides an abstraction for agent data models.
package agent

import (
	"time"

	"tkview/internal/organisation"
)

// ID is the unique identifier of an agent.
type ID string

// Agent is a tkview representation of an agent.
type Agent struct {
	ID       ID
	Name     string
	Type     string
	Version  string
	LastSeen time.Time
}

// Lister should return agents from a datasource.
type Lister interface {
	ListAgents(organisationID organisation.ID) ([]Agent, error)
}
