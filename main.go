// TKView is an application for viewing TestKube operations from the commandline.
package main

import (
	"github.com/rivo/tview"
)

func main() {
	// Init main frames
	organisationBox := tview.NewBox().
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetTitle("Organisations")
	agentBox := tview.NewBox().
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetTitle("Agents")
	executionBox := tview.NewBox().
		SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetTitle("Executions")

	// Lay out the frames
	topFlex := tview.NewFlex()
	topFlex.AddItem(organisationBox, 0, 1, false)
	topFlex.AddItem(agentBox, 0, 1, false)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(topFlex, 8, 1, false)
	flex.AddItem(executionBox, 0, 1, true)

	// Start 'er up!
	if err := tview.NewApplication().SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
