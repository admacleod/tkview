// TKView is an application for viewing TestKube operations from the commandline.
package main

import (
	"flag"
	"log"
	"os"

	"tkview/internal/testkube"
	"tkview/internal/tkview"
	"tkview/internal/ui"

	"github.com/charmbracelet/bubbletea/v2"
)

func main() {
	var (
		token string
		url   string
	)

	flag.StringVar(&token, "token", "", "API Token")
	flag.StringVar(&url, "url", "http://localhost:8099", "URL")
	flag.Parse()

	if token == "" {
		log.Println("You must provide an API Token")
		os.Exit(1)
	}

	client := testkube.New(url, token)
	tk := tkview.New(client)

	m := ui.NewModel(tk)

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		panic(err)
	}
}
