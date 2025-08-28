// Package tkview is the location of the core business logic for the tkview program.
package tkview

import (
	"errors"
	"fmt"

	"tkview/internal/agent"
	"tkview/internal/environment"
	"tkview/internal/organisation"
	"tkview/internal/workflow"
)

type client interface {
	agent.Lister
	environment.Lister
	organisation.Lister
	workflow.ExecutionLister
	workflow.Lister
}

// Organisation is a data structure that can be used to model the nested
// nature of organisations and environments.
type Organisation struct {
	organisation.Organisation

	Envs []environment.Environment
}

// Execution is a data structure that can be used to model the nested
// nature of executions and execution steps.
type Execution struct {
	workflow.Execution

	// TODO: steps
	// Steps []workflow.Steps
}

// Workflow is a data structure that can be used to model the nested
// nature of workflows and executions.
type Workflow struct {
	workflow.Workflow

	Executions []Execution
}

// TKView contains the core business logic for the tkview application.
// The default TKView can be used, however, it will return errors from
// every function. Instead use New() to create a new instance with
// a valid client implementation.
type TKView struct {
	client          client
	orgTree         []Organisation
	workflowTree    []Workflow
	currentOrg      organisation.ID
	currentEnv      environment.ID
	currentWorkflow workflow.ID
}

// New creates a new TKView using the passed client for accessing
// required data for running the application.
func New(c client) *TKView {
	return &TKView{
		client: c,
	}
}

var (
	errNoClient         = errors.New("no client")
	errNoOrgTree        = errors.New("organisations and environments are not populated")
	errNoOrgOrEnv       = errors.New("no organisation or environment is currently selected")
	errEnvNotFound      = errors.New("environment not found")
	errNoWorkflowTree   = errors.New("workflows are not populated")
	errWorkflowNotFound = errors.New("workflow not found")
	errNoWorkflow       = errors.New("no workflow is currently selected")
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

// SelectEnvironment checks whether the passed environment ExecutionID is known within the current
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

// GetWorkflowTree returns all workflows and executions belonging to the
// currently selected environment and organisation.
// If no organisation or environment is currently selected, it will error.
func (v *TKView) GetWorkflowTree() ([]Workflow, error) {
	if v.client == nil {
		return nil, errNoClient
	}

	if v.currentOrg == "" || v.currentEnv == "" {
		return nil, errNoOrgOrEnv
	}

	workflows, err := v.client.ListWorkflows(v.currentOrg, v.currentEnv)
	if err != nil {
		return nil, fmt.Errorf("list workflows: %w", err)
	}

	ret := make([]Workflow, 0, len(workflows))
	for _, w := range workflows {
		wf := Workflow{
			Workflow: w,
		}

		// Merge in any already known executions to prevent overfetching.
		for _, w := range v.workflowTree {
			if w.ID == wf.ID {
				wf.Executions = w.Executions

				break
			}
		}

		ret = append(ret, wf)
	}

	v.workflowTree = ret

	return ret, nil
}

// SelectWorkflow sets the passed workflow as the currently selected workflow.
// Upon selection the executions for the passed workflow will be populated and
// so subsequent calls to GetWorkflowTree will have these included.
func (v *TKView) SelectWorkflow(workflowID workflow.ID) error {
	if len(v.workflowTree) == 0 {
		return errNoWorkflowTree
	}

	found := false

	for _, w := range v.workflowTree {
		if w.ID == workflowID {
			v.currentWorkflow = w.ID
			found = true
		}
	}

	if !found {
		return errWorkflowNotFound
	}

	executions, err := v.client.ListExecutions(v.currentOrg, v.currentEnv, v.currentWorkflow)
	if err != nil {
		return fmt.Errorf("list executions for workflow %q: %w", workflowID, err)
	}

	ee := make([]Execution, 0, len(executions))
	for _, e := range executions {
		ee = append(ee, Execution{
			Execution: e,
		})
	}

	// Add the executions to the existing tree.
	for i, w := range v.workflowTree {
		if w.ID == workflowID {
			v.workflowTree[i].Executions = ee

			break
		}
	}

	return nil
}

// GetCurrentWorkflow returns the currently selected workflow, unless the workflow
// is no longer present in the workflow tree. This may be possible if stale data
// exists within the TKView model.
func (v *TKView) GetCurrentWorkflow() (Workflow, error) {
	if v.currentWorkflow == "" {
		return Workflow{}, errNoWorkflow
	}

	for _, w := range v.workflowTree {
		if w.ID == v.currentWorkflow {
			return w, nil
		}
	}

	// This shouldn't be possible.
	// TODO: maybe try repopulating and then selecting again?
	return Workflow{}, fmt.Errorf("workflow %q not currently known: %w", v.currentWorkflow, errWorkflowNotFound)
}
