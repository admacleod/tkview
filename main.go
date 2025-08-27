// TKView is an application for viewing TestKube operations from the commandline.
package main

import (
	"github.com/charmbracelet/bubbletea/v2"
)

func main() {
	m := NewModel()

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		panic(err)
	}
}
