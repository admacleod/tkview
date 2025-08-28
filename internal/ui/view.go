package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/tree"
)

// View renders the model for display on the terminal.
func (m Model) View() string {
	boxHeight := 8

	topBox := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder(), true, true, true, true).
		Height(boxHeight).
		Width(m.width / 2) //nolint:mnd // Half width for top boxes.

	mainFrame := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder(), true, true, true).
		Height(m.height - boxHeight).
		Width(m.width)

	organisationBox := topBox.
		Render(m.renderOrganisations())
	agentBox := topBox.
		Render(m.renderAgents())
	executionBox := mainFrame.
		Render("Executions")

	topRow := lipgloss.JoinHorizontal(0, organisationBox, agentBox)

	frame := lipgloss.JoinVertical(0, topRow, executionBox)

	return frame
}

func (m Model) renderOrganisations() string {
	currentEnv, err := m.tkview.GetCurrentEnvironment()
	if err != nil {
		return err.Error()
	}

	t := tree.Root("Environments").
		ItemStyleFunc(func(children tree.Children, i int) lipgloss.Style {
			if children.At(i).Value() == currentEnv.Name {
				return lipgloss.NewStyle().
					Background(lipgloss.Color("240"))
			}

			return lipgloss.NewStyle()
		})

	for _, org := range m.orgs {
		orgTree := tree.Root(org.Name)

		for _, e := range org.Envs {
			orgTree.Child(e.Name)
		}

		t.Child(orgTree)
	}

	return t.String()
}

func (m Model) renderAgents() string {
	if len(m.agents) == 0 {
		return "No agents found"
	}

	agents := make([]string, 0, len(m.agents))
	for _, agent := range m.agents {
		agents = append(agents, agent.Name)
	}

	return strings.Join(agents, "\n")
}
