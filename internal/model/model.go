// Package model contains the main model and program loop of tkview.
package model

import (
	"os"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// Model defines our Elm Architecture model for use in a tea program.
type Model struct {
	width, height int
}

// Init initialises the model so that it is ready for use.
func (Model) Init() tea.Cmd {
	return nil
}

// Update responds to tea messages to create a new model implementation.
//
//nolint:ireturn // This is required to fulfill the tea.Model interface.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var quitKey = key.NewBinding(key.WithKeys("ctrl+c"))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg.Key(), quitKey) {
			os.Exit(0)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

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

	organisationBox := topBox.Render("Organisations")
	agentBox := topBox.Render("Agents")
	executionBox := mainFrame.Render("Executions")

	topRow := lipgloss.JoinHorizontal(0, organisationBox, agentBox)

	frame := lipgloss.JoinVertical(0, topRow, executionBox)

	return frame
}
