package main

import (
	"github.com/mitchellh/cli"
	"github.com/munisystem/invoker/command"
)

func Commands(meta *command.Meta) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"apply": func() (cli.Command, error) {
			return &command.ApplyCommand{
				Meta: *meta,
			}, nil
		},
		"list": func() (cli.Command, error) {
			return &command.ListCommand{
				Meta: *meta,
			}, nil
		},
		"read": func() (cli.Command, error) {
			return &command.ReadCommand{
				Meta: *meta,
			}, nil
		},

		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Meta:     *meta,
				Version:  Version,
				Revision: Revision,
				Name:     Name,
			}, nil
		},
	}
}
