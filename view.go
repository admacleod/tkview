package main

import (
	"github.com/charmbracelet/lipgloss/v2"
)

// View renders the model for display on the terminal.
func (m Model) View() string {
	// If we aren't authed yet, then let the auth model handle this.
	if m.auth.Token() == nil {
		return m.auth.View()
	}

	boxHeight := 8

	topBox := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder(), true, true, true, true).
		Height(boxHeight).
		Width(m.width / 2) //nolint:mnd // Half width for top boxes.

	mainFrame := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder(), true, true, true).
		Height(m.height - boxHeight).
		Width(m.width)

	organisationBox := topBox.Render("Organisations")
	agentBox := topBox.Render("Agents")
	executionBox := mainFrame.Render("Executions")

	topRow := lipgloss.JoinHorizontal(0, organisationBox, agentBox)

	frame := lipgloss.JoinVertical(0, topRow, executionBox)

	return frame
}
