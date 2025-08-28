// Package testkube contains a client for a testkube API that fulfills the required interfaces
// that abstract working with tkview.
package testkube

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"tkview/internal/agent"
	"tkview/internal/environment"
	"tkview/internal/organisation"
	"tkview/internal/workflow"

	"github.com/kubeshop/testkube/cmd/kubectl-testkube/commands/agents"
	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/cloud/client"
)

// Client is a testkube API client.
// The default client cannot be used safely.
// Instead, initialise a client using New().
type Client struct {
	url   string
	token string
}

var (
	_ agent.Lister        = Client{}
	_ environment.Lister  = Client{}
	_ workflow.Lister     = Client{}
	_ organisation.Lister = Client{}
)

// New creates a valid testkube API client.
// The passed url should be the base url at which the API can be accessed.
// The token should be a valid token string for authenticating with the API.
func New(url, token string) Client {
	return Client{
		url:   url,
		token: token,
	}
}

// ListOrganisations returns all discoverable organisations using the provided API token.
func (c Client) ListOrganisations() ([]organisation.Organisation, error) {
	oo, err := client.NewOrganizationsClient(c.url, c.token).List()
	if err != nil {
		return nil, fmt.Errorf("list organisations at %q: %w", c.url, err)
	}

	ret := make([]organisation.Organisation, 0, len(oo))
	for _, o := range oo {
		ret = append(ret, organisation.Organisation{
			ID:   organisation.ID(o.Id),
			Name: o.Name,
		})
	}

	return ret, nil
}

// ListEnvironments returns all discoverable environments under the passed organisation using the provided API token.
func (c Client) ListEnvironments(organisationID organisation.ID) ([]environment.Environment, error) {
	ee, err := client.NewEnvironmentsClient(c.url, c.token, string(organisationID)).List()
	if err != nil {
		return nil, fmt.Errorf("list environments at %q for %q: %w", c.url, organisationID, err)
	}

	ret := make([]environment.Environment, 0, len(ee))
	for _, e := range ee {
		ret = append(ret, environment.Environment{
			ID:   environment.ID(e.Id),
			Name: e.Name,
		})
	}

	return ret, nil
}

// ListAgents returns all discoverable agents under the passed organisation using the provided API token.
func (c Client) ListAgents(organisationID organisation.ID) ([]agent.Agent, error) {
	aa, err := client.NewAgentsClient(c.url, c.token, string(organisationID)).List()
	if err != nil {
		return nil, fmt.Errorf("list agents at %q for %q: %w", c.url, organisationID, err)
	}

	ret := make([]agent.Agent, 0, len(aa))
	for _, a := range aa {
		agentType, err := agents.GetCliAgentType(a.Type)
		if err != nil {
			agentType = "Unknown"
		}

		ret = append(ret, agent.Agent{
			ID:       agent.ID(a.ID),
			Name:     a.Name,
			Type:     agentType,
			Version:  a.Version,
			LastSeen: *a.AccessedAt,
		})
	}

	return ret, nil
}

const (
	listWorkflowPath  = "%s/organizations/%s/environments/%s/agent/test-workflow-with-executions"
	listExecutionPath = "%s/organizations/%s/environments/%s/agent/test-workflows/%s/executions"
)

// ListWorkflows returns all workflows under the passed organisation and environment.
// Sadly, all parameters are required due to the testkube API.
func (c Client) ListWorkflows(orgID organisation.ID, envID environment.ID) ([]workflow.Workflow, error) {
	url := fmt.Sprintf(listWorkflowPath, c.url, orgID, envID)

	var result []testkube.TestWorkflowWithExecutionSummary
	if err := c.callTestKubeAPI(url, &result); err != nil {
		return nil, fmt.Errorf("call testkube api: %w", err)
	}

	ret := make([]workflow.Workflow, 0, len(result))
	for _, w := range result {
		if w.Workflow == nil {
			// Skip this broken workflow.
			continue
		}

		// Might not be any executions.
		var lastExecutionAt time.Time

		var lastExecutionStatus string

		if w.LatestExecution != nil {
			lastExecutionAt = w.LatestExecution.ScheduledAt
			if w.LatestExecution.Result != nil && w.LatestExecution.Result.Status != nil {
				lastExecutionStatus = string(*w.LatestExecution.Result.Status)
			}
		}

		ret = append(ret, workflow.Workflow{
			ID:                  workflow.ID(w.Workflow.Name),
			Name:                w.Workflow.Name,
			LastExecutionAt:     lastExecutionAt,
			LastExecutionStatus: lastExecutionStatus,
		})
	}

	return ret, nil
}

// ListExecutions returns all executions under the passed organisation, environment, and workflow.
// Sadly, all parameters are required due to the testkube API.
func (c Client) ListExecutions(orgID organisation.ID, envID environment.ID, id workflow.ID) ([]workflow.Execution, error) {
	url := fmt.Sprintf(listExecutionPath, c.url, orgID, envID, id)

	var result testkube.TestWorkflowExecutionsResult
	if err := c.callTestKubeAPI(url, &result); err != nil {
		return nil, fmt.Errorf("call testkube api: %w", err)
	}

	ret := make([]workflow.Execution, 0, len(result.Results))
	for _, e := range result.Results {
		status := "unknown"
		if e.Result != nil && e.Result.Status != nil {
			status = string(*e.Result.Status)
		}

		ret = append(ret, workflow.Execution{
			ID:        workflow.ExecutionID(e.Id),
			Name:      e.Name,
			StartedAt: e.ScheduledAt,
			Status:    status,
		})
	}

	return ret, nil
}

var errResponseCode = errors.New("unexpected HTTP status code")

func (c Client) callTestKubeAPI(url string, result any) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create request to %q: %w", url, err)
	}

	req.Header.Add("Authorization", "Bearer "+c.token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("doing request to %q: %w", url, err)
	}

	defer func() {
		// Too late to handle this error, and we don't really care.
		_ = res.Body.Close()
	}()

	switch res.StatusCode {
	case http.StatusRequestTimeout:
		// There may be an issue at the agent, this is fine, but nothing will be returned.
		return nil
	case http.StatusOK:
		// Everything is fine and working as expected!
		break
	default:
		return fmt.Errorf("request to %q returned status %d: %w", url, res.StatusCode, errResponseCode)
	}

	if err := json.NewDecoder(res.Body).Decode(result); err != nil {
		return fmt.Errorf("decode response body: %w", err)
	}

	return nil
}
