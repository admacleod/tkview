package auth

import (
	"fmt"
)

// View renders the UI for authentication.
func (m Model) View() string {
	if m.authCodeURL == "" {
		return fmt.Sprintf(`You aren't logged in.
Enter the authentication URL for your Testkube installation below:

%s
`, m.authURL.View())
	}

	return "Please visit the following URL to authenticate:\n\n" + m.authCodeURL
}
