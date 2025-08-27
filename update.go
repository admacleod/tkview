package main

import (
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbletea/v2"
)

// Update responds to tea messages to create a new model implementation.
//
//nolint:ireturn // This is required to fulfill the tea.Model interface.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Basic messages for the general good behaviour of the program.
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg.Key(), m.keyMap.Quit) {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		return m, nil
	}

	// Update submodels.
	var cmd tea.Cmd

	m.auth, cmd = m.auth.Update(msg)

	return m, cmd
}
