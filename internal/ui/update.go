//nolint:ireturn // tea.Model and tea.Msg are commonly returned from update functions and commands.
package ui

import (
	"fmt"

	"tkview/internal/agent"
	"tkview/internal/environment"
	"tkview/internal/execution"
	"tkview/internal/tkview"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbletea/v2"
)

type errMsg error
type focusMsg view
type orgTreeMsg []tkview.Organisation
type envMsg environment.ID
type agentsMsg []agent.Agent
type executionsMsg []execution.Execution

func errCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return errMsg(err)
	}
}

func focusCmd(v view) tea.Cmd {
	return func() tea.Msg {
		return focusMsg(v)
	}
}

func switchEnvCmd(id environment.ID) tea.Cmd {
	return func() tea.Msg {
		return envMsg(id)
	}
}

// Update responds to tea messages to create a new model implementation.
//
//nolint:cyclop,funlen // This is going to be a big function because it handles all the UI update logic.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Basic messages for the general good behaviour of the program.
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg.Key(), m.keyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg.Key(), m.keyMap.Next):
			switch m.focused {
			case viewOrgs:
				return m, m.nextOrgEnv
			case viewAgents:
			case viewExecutions:
			}
		case key.Matches(msg.Key(), m.keyMap.Prev):
			switch m.focused {
			case viewOrgs:
				return m, m.prevOrgEnv
			case viewAgents:
			case viewExecutions:
			}
		case key.Matches(msg.Key(), m.keyMap.FocusOrganisations):
			return m, focusCmd(viewOrgs)
		case key.Matches(msg.Key(), m.keyMap.FocusAgents):
			return m, focusCmd(viewAgents)
		case key.Matches(msg.Key(), m.keyMap.FocusExecutions):
			return m, focusCmd(viewExecutions)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		return m, nil
	case errMsg:
		// TODO: display errors
		panic(msg)
	case orgTreeMsg:
		m.orgs = msg

		// Select the initial environment.
		// TODO: this is dangerous, the first org might not have any environments!
		return m, switchEnvCmd(msg[0].Envs[0].ID)
	case envMsg:
		err := m.tkview.SelectEnvironment(environment.ID(msg))
		if err != nil {
			return m, errCmd(err)
		}

		// After the environment changes, load the agents and executions again.
		return m, tea.Batch(
			m.loadAgents,
			m.loadExecutions,
		)
	case agentsMsg:
		m.agents = msg

		return m, nil
	case executionsMsg:
		m.executions = msg

		return m, nil
	case focusMsg:
		m.focused = view(msg)

		return m, nil
	}

	return m, nil
}

func (m Model) getOrgTree() tea.Msg {
	t, err := m.tkview.GetOrganisationTree()
	if err != nil {
		return errMsg(fmt.Errorf("get organisation tree: %w", err))
	}

	return orgTreeMsg(t)
}

func (m Model) nextOrgEnv() tea.Msg {
	currentEnv, err := m.tkview.GetCurrentEnvironment()
	if err != nil {
		return errMsg(fmt.Errorf("get current environment: %w", err))
	}

	for i, o := range m.orgs {
		for j, e := range o.Envs {
			if e.ID == currentEnv.ID {
				switch {
				case i+1 == len(m.orgs) && j+1 == len(o.Envs):
					// Nowhere to go such loop back to the start.
					// TODO: this is dangerous, the first org might not have any environments!
					return envMsg(m.orgs[0].Envs[0].ID)
				case j+1 == len(o.Envs):
					// All out of envs in this org, go to the next org and select the first environment there.
					// TODO: this is dangerous, the next org might not have any environments!
					return envMsg(m.orgs[i+1].Envs[0].ID)
				default:
					// Go to the next environment in this org.
					// This should be safe because of the length checks that have occurred above.
					return envMsg(m.orgs[i].Envs[j+1].ID)
				}
			}
		}
	}
	// This shouldn't happen, but could.
	// TODO: probably do something here
	return nil
}

func (m Model) prevOrgEnv() tea.Msg {
	currentEnv, err := m.tkview.GetCurrentEnvironment()
	if err != nil {
		return errMsg(fmt.Errorf("get current environment: %w", err))
	}

	for i, o := range m.orgs {
		for j, e := range o.Envs {
			if e.ID == currentEnv.ID {
				switch {
				case i-1 < 0 && j-1 < 0:
					// Nowhere to go such loop back to the bottom.
					// TODO: this is dangerous, the last org might not have any environments!
					return envMsg(m.orgs[len(m.orgs)-1].Envs[len(m.orgs[len(m.orgs)-1].Envs)-1].ID)
				case j-1 < 0:
					// All out of envs in this org, go to the previous org and select the last environment there.
					// TODO: this is dangerous, the previous org might not have any environments!
					return envMsg(m.orgs[i-1].Envs[len(m.orgs[i-1].Envs)-1].ID)
				default:
					// Go to the previous environment in this org.
					// This should be safe because of the checks that have occurred above.
					return envMsg(m.orgs[i].Envs[j-1].ID)
				}
			}
		}
	}
	// This shouldn't happen, but could.
	// TODO: probably do something here
	return nil
}

func (m Model) loadAgents() tea.Msg {
	agents, err := m.tkview.GetAgents()
	if err != nil {
		return errMsg(fmt.Errorf("get agents: %w", err))
	}

	return agentsMsg(agents)
}

func (m Model) loadExecutions() tea.Msg {
	executions, err := m.tkview.GetExecutions()
	if err != nil {
		return errMsg(fmt.Errorf("get executions: %w", err))
	}

	return executionsMsg(executions)
}
