package ui

import (
	"github.com/charmbracelet/bubbles/v2/key"
)

type keyMap struct {
	Quit               key.Binding
	Next               key.Binding
	Prev               key.Binding
	FocusOrganisations key.Binding
	FocusAgents        key.Binding
	FocusExecutions    key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		Quit:               key.NewBinding(key.WithKeys("ctrl+c")),
		Next:               key.NewBinding(key.WithKeys("down")),
		Prev:               key.NewBinding(key.WithKeys("up")),
		FocusOrganisations: key.NewBinding(key.WithKeys("shift+o")),
		FocusAgents:        key.NewBinding(key.WithKeys("shift+a")),
		FocusExecutions:    key.NewBinding(key.WithKeys("shift+e")),
	}
}
