package auth

import (
	"github.com/charmbracelet/bubbles/v2/textinput"
	"golang.org/x/oauth2"
)

const listenPort = 8090

// Model is mostly a tea.Model implementation.
// The underlying oauth2 token can be accessed using the Token function.
type Model struct {
	authURL     textinput.Model
	authCodeURL string
	listenPort  int
	token       *oauth2.Token
}

// NewModel creates a new auth Model that is ready for use.
func NewModel() Model {
	authURL := textinput.New()
	authURL.Placeholder = "http://localhost:5556"
	authURL.Focus()
	authURL.CharLimit = 256

	return Model{
		authURL:    authURL,
		listenPort: listenPort,
	}
}

// Token returns the underlying authentication token.
// If the model is not currently authenticated, then
// the returned token will be nil.
func (m Model) Token() *oauth2.Token {
	return m.token
}
