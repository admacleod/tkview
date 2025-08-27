package main

import (
	"tkview/internal/auth"

	"github.com/charmbracelet/bubbles/v2/textinput"
	"github.com/charmbracelet/bubbletea/v2"
)

// Model defines our Elm Architecture model for use in a tea program.
type Model struct {
	width, height int
	keyMap        keyMap
	auth          auth.Model
}

// NewModel creates a new Model.
// It will not be initialised and so Init should be called
// before first use to ensure that everything operates as expected.
func NewModel() Model {
	return Model{
		width:  0,
		height: 0,
		keyMap: defaultKeyMap(),
		auth:   auth.NewModel(),
	}
}

// Init initialises the model so that it is ready for use.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
	)
}
