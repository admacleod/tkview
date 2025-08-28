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
	"tkview/internal/execution"
	"tkview/internal/organisation"

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
	_ execution.Lister    = Client{}
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

const listExecutionPath = "%s/organizations/%s/environments/%s/agent/test-workflow-executions"

var errResponseCode = errors.New("unexpected HTTP status code")

// ListExecutions returns all executions under the passed organisation and environment,
// sadly both are required due to the testkube API.
func (c Client) ListExecutions(orgID organisation.ID, envID environment.ID) ([]execution.Execution, error) {
	url := fmt.Sprintf(listExecutionPath, c.url, orgID, envID)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create list executions request to %q: %w", url, err)
	}

	req.Header.Add("Authorization", "Bearer "+c.token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("doing list executions request to %q: %w", url, err)
	}

	defer func() {
		// Too late to handle this error, and we don't really care.
		_ = res.Body.Close()
	}()

	switch res.StatusCode {
	case http.StatusRequestTimeout:
		// There may be an issue at the agent, this is fine, but nothing will be returned.
		return nil, nil
	case http.StatusOK:
		// Everything is fine and working as expected!
		break
	default:
		return nil, fmt.Errorf("list executions request to %q returned status %d: %w", url, res.StatusCode, errResponseCode)
	}

	var result testkube.TestWorkflowExecutionsResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode list executions response body: %w", err)
	}

	ret := make([]execution.Execution, 0, len(result.Results))
	for _, e := range result.Results {
		ret = append(ret, execution.Execution{
			ID:        execution.ID(e.Id),
			Name:      e.Name,
			StartedAt: e.ScheduledAt,
			Status:    string(*e.Result.Status),
		})
	}

	return ret, nil
}
