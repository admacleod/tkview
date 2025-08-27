// Package auth provides a model for authenticating with a Testkube instance.
package auth

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbletea/v2"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

const (
	clientID  = "testkube-cloud-cli"
	localhost = "127.0.0.1"
)

func server(oidcProviderURL string, port int) (tea.Cmd, string, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, oidcProviderURL)
	if err != nil {
		return nil, "", fmt.Errorf("new OIDC provider from %q: %w", oidcProviderURL, err)
	}

	httpListenURL := net.JoinHostPort(localhost, strconv.Itoa(port))

	oauth2Config := oauth2.Config{
		ClientID:    clientID,
		Endpoint:    provider.Endpoint(),
		RedirectURL: fmt.Sprintf("http://%s/callback", httpListenURL),
		Scopes:      []string{oidc.ScopeOpenID, "profile", "email", "offline_access"},
	}

	cmd := serveCmd(httpListenURL, oauth2Config, provider)
	authURL := oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)

	return cmd, authURL, nil
}

func serveCmd(httpListenURL string, oauth2Config oauth2.Config, provider *oidc.Provider) tea.Cmd {
	return func() tea.Msg {
		var token *oauth2.Token
		// Start a local server to handle the callback from the OIDC provider.
		srv := &http.Server{
			Addr:              httpListenURL,
			ReadHeaderTimeout: time.Minute,
		}
		srv.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// This endpoint should only be called once, and then the server should exit.
			// This call is only to release ListenAndServe when the function is complete.
			defer func(srv *http.Server, ctx context.Context) {
				// We don't care about errors here, there isn't anything that can be done.
				_ = srv.Shutdown(ctx)
			}(srv, r.Context())

			code := r.URL.Query().Get("code")
			if code == "" {
				http.Error(w, "Authorization failed. Missing code.", http.StatusUnauthorized)

				return
			}

			// Get the oauth token.
			var err error

			token, err = oauth2Config.Exchange(r.Context(), code)
			if err != nil {
				http.Error(w, "Token exchange failed", http.StatusInternalServerError)

				return
			}
			// Validate the oauth token.
			if _, err := provider.Verifier(&oidc.Config{ClientID: clientID}).Verify(r.Context(), token.AccessToken); err != nil {
				http.Error(w, "Token verification failed", http.StatusInternalServerError)

				return
			}

			if _, err := fmt.Fprint(w,
				"<script>window.close()</script>",
				"TKView is now successfully authenticated. Go back to the terminal to continue.",
			); err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		})

		err := srv.ListenAndServe()
		switch {
		case errors.Is(err, http.ErrServerClosed):
			// This is expected.
			break
		case err != nil:
			return errMsg(fmt.Errorf("listen for oauth2 callback on %q: %w", httpListenURL, err))
		}

		return tokenMsg(token)
	}
}
