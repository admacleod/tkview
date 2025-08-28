package ui

import (
	"time"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/table"
	"github.com/charmbracelet/lipgloss/v2/tree"
)

// View renders the model for display on the terminal.
func (m Model) View() string {
	frame := lipgloss.JoinVertical(0,
		lipgloss.JoinHorizontal(0,
			m.renderOrganisations(),
			m.renderAgents(),
		),
		m.renderExecutions(),
	)

	return frame
}

func (m Model) renderOrganisations() string {
	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, true, true, true).
		Height(m.topBoxHeight).
		Width(m.width / m.topBoxCount)

	if m.focused == viewOrgs {
		box = box.BorderStyle(lipgloss.DoubleBorder())
	}

	currentEnv, err := m.tkview.GetCurrentEnvironment()
	if err != nil {
		return box.Render(err.Error())
	}

	t := tree.Root("(O)rganisations and Environments").
		ItemStyleFunc(func(children tree.Children, i int) lipgloss.Style {
			if children.At(i).Value() == currentEnv.Name {
				return lipgloss.NewStyle().
					Background(lipgloss.BrightBlue).
					Foreground(lipgloss.White)
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

	return box.Render(t.String())
}

func (m Model) renderAgents() string {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		Height(m.topBoxHeight).
		Width(m.width/m.topBoxCount).
		Wrap(true).
		Headers("(A)gents", "Type", "Version", "LastSeen")

	if m.focused == viewAgents {
		t.Border(lipgloss.DoubleBorder())
	}

	if len(m.agents) == 0 {
		t.Row("No agents found", "", "", "")
		m.padAgentTable(t)

		return t.Render()
	}

	for _, agent := range m.agents {
		lastSeenFormat := time.RFC822
		// If last seen today only show the time.
		if agent.LastSeen.Format(time.DateOnly) == time.Now().Format(time.DateOnly) {
			lastSeenFormat = "15:04 MST"
		}

		t.Row(agent.Name, agent.Type, agent.Version, agent.LastSeen.Format(lastSeenFormat))
	}

	if len(m.agents) < m.topBoxTableRows {
		m.padAgentTable(t)
	}

	return t.Render()
}

func (m Model) padAgentTable(t *table.Table) {
	blankRow := []string{"", "", "", ""}
	for range m.topBoxTableRows - len(m.agents) {
		t.Row(blankRow...)
	}
}

func (m Model) renderExecutions() string {
	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, true, true).
		Height(m.height - m.topBoxHeight).
		Width(m.width)

	if m.focused == viewExecutions {
		box = box.BorderStyle(lipgloss.DoubleBorder())
	}

	return box.Render("(E)xecutions")
}
