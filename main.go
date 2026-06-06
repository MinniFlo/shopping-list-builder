package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
	p := tea.NewProgram(initialModel())
	if m, err := p.Run(); err != nil {
		fmt.Printf("There has been an error: %v", err)
		fmt.Printf("The last model state was: %v", m)
		os.Exit(1)
	}
}
