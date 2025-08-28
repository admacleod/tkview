// Package tkview is the location of the core business logic for the tkview program.
package tkview

import (
	"errors"
	"fmt"

	"tkview/internal/agent"
	"tkview/internal/environment"
	"tkview/internal/execution"
	"tkview/internal/organisation"
)

type client interface {
	agent.Lister
	environment.Lister
	execution.Lister
	organisation.Lister
}

// Organisation is a data structure that can be used to model the nested
// nature of organisations and environments.
type Organisation struct {
	organisation.Organisation

	Envs []environment.Environment
}

// TKView contains the core business logic for the tkview application.
// The default TKView can be used, however, it will return errors from
// every function. Instead use New() to create a new instance with
// a valid client implementation.
type TKView struct {
	client     client
	orgTree    []Organisation
	currentOrg organisation.ID
	currentEnv environment.ID
}

// New creates a new TKView using the passed client for accessing
// required data for running the application.
func New(c client) *TKView {
	return &TKView{
		client: c,
	}
}

var (
	errNoClient    = errors.New("no client")
	errNoOrgTree   = errors.New("organisations and environments are not populated")
	errNoOrgOrEnv  = errors.New("no organisation or environment is currently selected")
	errEnvNotFound = errors.New("environment not found")
)

// GetOrganisationTree updates and then returns the TKView organisation tree.
func (v *TKView) GetOrganisationTree() ([]Organisation, error) {
	if v.client == nil {
		return nil, errNoClient
	}

	orgs, err := v.client.ListOrganisations()
	if err != nil {
		return nil, fmt.Errorf("list organisations: %w", err)
	}

	v.orgTree = []Organisation{}
	for _, org := range orgs {
		t := Organisation{
			Organisation: org,
			Envs:         []environment.Environment{},
		}

		envs, err := v.client.ListEnvironments(org.ID)
		if err != nil {
			// If we get an error, just blank the current tree to avoid a parietal tree.
			v.orgTree = []Organisation{}

			return nil, fmt.Errorf("list environments for org %q: %w", org.Name, err)
		}

		t.Envs = envs

		v.orgTree = append(v.orgTree, t)
	}

	return v.orgTree, nil
}

// SelectEnvironment checks whether the passed environment ID is known within the current
// organisation tree and if it does, it will set it as the currently selected environment.
func (v *TKView) SelectEnvironment(envID environment.ID) error {
	if len(v.orgTree) == 0 {
		return errNoOrgTree
	}

	for _, org := range v.orgTree {
		for _, env := range org.Envs {
			if env.ID == envID {
				v.currentOrg = org.ID
				v.currentEnv = env.ID

				return nil
			}
		}
	}

	return fmt.Errorf("environment %q not currently known: %w", envID, errEnvNotFound)
}

// GetCurrentEnvironment returns the currently selected environment, unless the environment
// is no longer present in the organisation tree. This may be possible if stale data
// exists within the TKView model.
func (v *TKView) GetCurrentEnvironment() (environment.Environment, error) {
	if v.currentEnv == "" {
		return environment.Environment{}, errNoOrgOrEnv
	}

	for _, org := range v.orgTree {
		for _, env := range org.Envs {
			if env.ID == v.currentEnv {
				return env, nil
			}
		}
	}

	// This shouldn't be possible.
	// TODO: maybe try repopulating and then selecting again?
	return environment.Environment{}, fmt.Errorf("environment %q not currently known: %w", v.currentEnv, errEnvNotFound)
}

// GetAgents returns all agents belonging to the organisation parent of the currently selected environment.
// If no environment is currently selected, it will error.
func (v *TKView) GetAgents() ([]agent.Agent, error) {
	if v.client == nil {
		return nil, errNoClient
	}

	if v.currentOrg == "" {
		return nil, errNoOrgOrEnv
	}

	agents, err := v.client.ListAgents(v.currentOrg)
	if err != nil {
		return nil, fmt.Errorf("list agents: %w", err)
	}

	return agents, nil
}

// GetExecutions returns all executions belonging to the currently selected environment.
// If no environment is currently selected, it will error.
func (v *TKView) GetExecutions() ([]execution.Execution, error) {
	if v.client == nil {
		return nil, errNoClient
	}

	if v.currentEnv == "" {
		return nil, errNoOrgOrEnv
	}

	executions, err := v.client.ListExecutions(v.currentOrg, v.currentEnv)
	if err != nil {
		return nil, fmt.Errorf("list executions: %w", err)
	}

	return executions, nil
}
