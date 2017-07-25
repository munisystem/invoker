package command

import (
	"strings"
)

type ReadCommand struct {
	Meta
}

func (c *ReadCommand) Run(args []string) int {
	// Write your code here

	return 0
}

func (c *ReadCommand) Synopsis() string {
	return ""
}

func (c *ReadCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
