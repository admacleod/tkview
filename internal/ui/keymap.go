package ui

import (
	"github.com/charmbracelet/bubbles/v2/key"
)

type keyMap struct {
	Quit              key.Binding
	Next              key.Binding
	Prev              key.Binding
	Select            key.Binding
	FocusEnvironments key.Binding
	FocusAgents       key.Binding
	FocusWorkflows    key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		Quit:              key.NewBinding(key.WithKeys("ctrl+c")),
		Next:              key.NewBinding(key.WithKeys("down")),
		Prev:              key.NewBinding(key.WithKeys("up")),
		Select:            key.NewBinding(key.WithKeys("enter")),
		FocusEnvironments: key.NewBinding(key.WithKeys("shift+e")),
		FocusAgents:       key.NewBinding(key.WithKeys("shift+a")),
		FocusWorkflows:    key.NewBinding(key.WithKeys("shift+w")),
	}
}
