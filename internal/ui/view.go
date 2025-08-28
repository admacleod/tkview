package ui

import (
	"fmt"
	"strings"
	"time"

	"tkview/internal/tkview"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/table"
	"github.com/charmbracelet/lipgloss/v2/tree"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
)

// View renders the model for display on the terminal.
func (m Model) View() string {
	frame := lipgloss.JoinVertical(0,
		lipgloss.JoinHorizontal(0,
			m.renderOrganisations(),
			m.renderAgents(),
		),
		m.renderWorkflows(),
	)

	return frame
}

func (m Model) renderOrganisations() string {
	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, true, true, true).
		Height(m.topBoxHeight).
		Width(m.width / m.topBoxCount)

	if m.focused == viewEnvs {
		box = box.BorderStyle(lipgloss.DoubleBorder())
	}

	currentEnv, err := m.tkview.GetCurrentEnvironment()
	if err != nil {
		return box.Render(err.Error())
	}

	t := tree.Root("(E)nvironments").
		ItemStyleFunc(func(children tree.Children, i int) lipgloss.Style {
			if children.At(i).Value() == currentEnv.Name {
				return lipgloss.NewStyle().
					Background(lipgloss.BrightBlue).
					Foreground(lipgloss.White)
			}

			return lipgloss.NewStyle()
		})

	for _, org := range m.orgs {
		orgTree := tree.Root(org.Name)

		for _, e := range org.Envs {
			orgTree.Child(e.Name)
		}

		t.Child(orgTree)
	}

	return box.Render(t.String())
}

func (m Model) renderAgents() string {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		Height(m.topBoxHeight).
		Width(m.width/m.topBoxCount).
		Wrap(true).
		Headers("(A)gents", "Type", "Version", "LastSeen")

	if m.focused == viewAgents {
		t.Border(lipgloss.DoubleBorder())
	}

	if len(m.agents) == 0 {
		t.Row("No agents found", "", "", "")
		m.padAgentTable(t)

		return t.Render()
	}

	for _, agent := range m.agents {
		t.Row(agent.Name, agent.Type, agent.Version, renderTime(agent.LastSeen))
	}

	m.padAgentTable(t)

	return t.Render()
}

func (m Model) padAgentTable(t *table.Table) {
	blankRow := []string{"", "", "", ""}
	for range m.topBoxHeight - m.tableBorderHeight - len(m.agents) {
		t.Row(blankRow...)
	}
}

func (m Model) renderWorkflows() string {
	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, true, true, true).
		Height(m.height - m.topBoxHeight).
		Width(m.width)

	if m.focused == viewWorkflows {
		box = box.BorderStyle(lipgloss.DoubleBorder())
	}

	currentWorkflow, err := m.tkview.GetCurrentWorkflow()
	if err != nil {
		return box.Render(err.Error())
	}

	title := "(W)orkflows"
	if m.focused == viewWorkflows {
		title += " | (s)tart"
	}

	t := tree.Root(title).
		ItemStyleFunc(func(children tree.Children, i int) lipgloss.Style {
			if children.At(i).Value() == m.renderWorkflow(currentWorkflow) {
				return lipgloss.NewStyle().
					Background(lipgloss.BrightBlue).
					Foreground(lipgloss.White)
			}

			return lipgloss.NewStyle()
		})

	for _, workflow := range m.workflows {
		workflowTree := tree.Root(m.renderWorkflow(workflow))

		_, expanded := m.expandedWorkflows[workflow.ID]

		if expanded {
			for _, execution := range workflow.Executions {
				executionTree := tree.Root(m.renderExecution(execution))
				workflowTree.Child(executionTree)
			}
		}

		t.Child(workflowTree)
	}

	return box.Render(t.String())
}

func (m Model) renderWorkflow(workflow tkview.Workflow) string {
	t := renderTime(workflow.LastExecutionAt)

	timePadding := ""
	if len(t) < len(time.RFC822) {
		timePadding = strings.Join(make([]string, len(time.RFC822)-len(t)), " ")
	}

	return fmt.Sprintf("%s\t%s%s\t%s",
		renderStatus(workflow.LastExecutionStatus),
		timePadding,
		t,
		workflow.Name,
	)
}

func renderTime(t time.Time) string {
	format := time.RFC822
	// If time is today, only show the time.
	if t.Format(time.DateOnly) == time.Now().Format(time.DateOnly) {
		format = "15:04 MST"
	}

	return t.Format(format)
}

func renderStatus(s string) string {
	status := "❓"

	switch testkube.TestWorkflowStatus(s) {
	case testkube.QUEUED_TestWorkflowStatus,
		testkube.PENDING_TestWorkflowStatus:
		status = "🚶"
	case testkube.STARTING_TestWorkflowStatus,
		testkube.SCHEDULING_TestWorkflowStatus,
		testkube.PAUSING_TestWorkflowStatus,
		testkube.RESUMING_TestWorkflowStatus,
		testkube.STOPPING_TestWorkflowStatus:
		status = "⏳"
	case testkube.RUNNING_TestWorkflowStatus:
		status = "🔄"
	case testkube.PAUSED_TestWorkflowStatus:
		status = "⏸️"
	case testkube.ABORTED_TestWorkflowStatus,
		testkube.CANCELED_TestWorkflowStatus:
		status = "🛑"
	case testkube.PASSED_TestWorkflowStatus:
		status = "✅"
	case testkube.FAILED_TestWorkflowStatus:
		status = "❌"
	}

	return status
}

func (m Model) renderExecution(execution tkview.Execution) string {
	t := renderTime(execution.StartedAt)

	timePadding := ""
	if len(t) < len(time.RFC822) {
		timePadding = strings.Join(make([]string, len(time.RFC822)-len(t)), " ")
	}

	return fmt.Sprintf("%s\t%s%s\t%s",
		renderStatus(execution.Status),
		timePadding,
		t,
		execution.Name,
	)
}
