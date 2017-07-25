package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestReadCommand_implement(t *testing.T) {
	var _ cli.Command = &ReadCommand{}
}
