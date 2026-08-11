package ui

import (
	"golangutils/pkg/models"
)

func InfoNofity(message string, icon string) error {
	return notify("Information", message, icon)
}

func WarnNofity(message string, icon string) error {
	return notify("Warnning", message, icon)
}

func ErrorNofity(message string, icon string) error {
	return notify("Error", message, icon)
}

func OkNofity(message string, icon string) error {
	return notify("Success", message, icon)
}

func SelectFile(title string) models.Response[string] {
	return selectDialog(title, true)
}

func SelectFolder(title string) models.Response[string] {
	return selectDialog(title, false)
}

func InfoDialog(title string, message string) error {
	return dialogBox(title, message, typeInfo)
}

func WarnDialog(title string, message string) error {
	return dialogBox(title, message, typeWarning)
}

func ErrorDialog(title string, message string) error {
	return dialogBox(title, message, typeError)
}
