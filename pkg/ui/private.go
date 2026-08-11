package ui

import "golangutils/pkg/models"

func getDefaultCmd() models.Command {
	return models.Command{Verbose: verbose, IsAsync: true}
}
