package ui

import (
	"tkview/internal/agent"
	"tkview/internal/tkview"
	"tkview/internal/workflow"

	"github.com/charmbracelet/bubbles/v2/textinput"
	"github.com/charmbracelet/bubbletea/v2"
)

type view int

const (
	viewEnvs view = iota
	viewAgents
	viewWorkflows
)

const (
	uiTopBoxCount     = 2
	uiTopBoxMaxHeight = 8
	// 3 for the top, header, and bottom, borders.
	// 1 for the header line.
	uiTableBorderHeight = 4
)

// Model defines our Elm Architecture model for use in a tea program.
type Model struct {
	width, height     int
	topBoxCount       int
	topBoxHeight      int
	tableBorderHeight int
	keyMap            keyMap
	tkview            *tkview.TKView
	focused           view
	orgs              []tkview.Organisation
	agents            []agent.Agent
	workflows         []tkview.Workflow
	expandedWorkflows map[workflow.ID]struct{}
}

// NewModel creates a new Model.
// It will not be initialised and so Init should be called
// before first use to ensure that everything operates as expected.
func NewModel(tkview *tkview.TKView) Model {
	return Model{
		width:             0,
		height:            0,
		topBoxCount:       uiTopBoxCount,
		topBoxHeight:      uiTopBoxMaxHeight,
		tableBorderHeight: uiTableBorderHeight,
		keyMap:            defaultKeyMap(),
		tkview:            tkview,
		expandedWorkflows: make(map[workflow.ID]struct{}),
	}
}

// Init initialises the model so that it is ready for use.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.getOrgTree,
	)
}
