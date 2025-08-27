package auth

import (
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbletea/v2"
	"golang.org/x/oauth2"
)

type errMsg error
type tokenMsg *oauth2.Token

// Update processes messages to affect the underlying model.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case errMsg:
		panic(msg)
	case tea.KeyMsg:
		if key.Matches(msg.Key(), key.NewBinding(key.WithKeys("enter"))) && m.authURL.Focused() {
			var err error

			cmd, m.authCodeURL, err = server(m.authURL.Value(), m.listenPort)
			if err != nil {
				panic(err)
			}

			return m, cmd
		}
	case tokenMsg:
		m.token = msg

		return m, nil
	}

	// Update submodels.
	m.authURL, cmd = m.authURL.Update(msg)

	return m, cmd
}
