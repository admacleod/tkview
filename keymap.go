package main

import (
	"github.com/charmbracelet/bubbles/v2/key"
)

type keyMap struct {
	Quit key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		Quit: key.NewBinding(key.WithKeys("ctrl+c")),
	}
}
