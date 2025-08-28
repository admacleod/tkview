package ui

import (
	"tkview/internal/agent"
	"tkview/internal/execution"
	"tkview/internal/tkview"

	"github.com/charmbracelet/bubbles/v2/textinput"
	"github.com/charmbracelet/bubbletea/v2"
)

type view int

const (
	viewOrgs view = iota
	viewAgents
	viewExecutions
)

// Model defines our Elm Architecture model for use in a tea program.
type Model struct {
	width, height int
	keyMap        keyMap
	tkview        *tkview.TKView
	focused       view
	orgs          []tkview.Organisation
	agents        []agent.Agent
	executions    []execution.Execution
}

// NewModel creates a new Model.
// It will not be initialised and so Init should be called
// before first use to ensure that everything operates as expected.
func NewModel(tkview *tkview.TKView) Model {
	return Model{
		width:  0,
		height: 0,
		keyMap: defaultKeyMap(),
		tkview: tkview,
	}
}

// Init initialises the model so that it is ready for use.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.getOrgTree,
	)
}
