// TKView is an application for viewing TestKube operations from the commandline.
package main

import (
	"tkview/internal/model"

	"github.com/charmbracelet/bubbletea/v2"
)

func main() {
	m := model.Model{}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		panic(err)
	}
}
