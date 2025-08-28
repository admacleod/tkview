package ui

import (
	"github.com/charmbracelet/bubbles/v2/key"
)

type keyMap struct {
	Quit key.Binding
	Next key.Binding
	Prev key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		Quit: key.NewBinding(key.WithKeys("ctrl+c")),
		Next: key.NewBinding(key.WithKeys("down")),
		Prev: key.NewBinding(key.WithKeys("up")),
	}
}
