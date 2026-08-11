//go:build linux

package ui

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/console"
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
)

const TIMEOUT_NOTIFICATION = 60000

// FALLBACK UNIVERSAL LINUX: direct call to D-Bus via dbus-send
func dbusCmd(title string, message string, icon string) models.Command {
	cmd := getDefaultCmd()
	cmd.Cmd = `dbus-send --session --dest=org.freedesktop.Notifications --type=method_call /org/freedesktop/Notifications org.freedesktop.Notifications.Notify`
	cmd.Cmd = fmt.Sprintf(`%s string:%s uint32:0`, cmd.Cmd, common.GetAppId()) // App name
	if !str.IsEmpty(icon) {
		cmd.Cmd = fmt.Sprintf(`%s string:%s`, cmd.Cmd, icon) // App icon
	}
	cmd.Cmd = fmt.Sprintf(`%s string:%s string:%s array:string: dict:string:subsetting: int32:%d`, cmd.Cmd, title, message, TIMEOUT_NOTIFICATION)
	return cmd
}

func notifySendCmd(title string, message string, icon string) models.Command {
	cmd := getDefaultCmd()
	cmd.Cmd = fmt.Sprintf(`notify-send "%s" "%s" -a %s -t %d`, title, message, common.GetAppId(), TIMEOUT_NOTIFICATION)
	if !str.IsEmpty(icon) {
		cmd.Cmd = fmt.Sprintf(`%s -i "%s"`, cmd.Cmd, icon)
	}
	return cmd
}

func kdialogCmd(title string, message string, icon string) models.Command {
	cmd := getDefaultCmd()
	cmd.Cmd = fmt.Sprintf(`kdialog --title "%s" --passivepopup "%s" %d`, title, message, TIMEOUT_NOTIFICATION)
	if !str.IsEmpty(icon) {
		cmd.Cmd = fmt.Sprintf(`%s --icon "%s"`, cmd.Cmd, icon)
	}
	return cmd
}

func notify(title string, message string, icon string) error {
	notifySendCmd := notifySendCmd(title, message, icon)
	dbusCmd := dbusCmd(title, message, icon)
	kdialog := kdialogCmd(title, message, icon)
	if !str.IsEmpty(console.WhichIgnoreError(notifySendCmd.Cmd)) {
		return exe.ExecRealTime(notifySendCmd)
	} else if !str.IsEmpty(console.WhichIgnoreError(dbusCmd.Cmd)) {
		return exe.ExecRealTime(dbusCmd)
	} else if !str.IsEmpty(console.WhichIgnoreError(kdialog.Cmd)) {
		return exe.ExecRealTime(kdialog)
	}
	return fmt.Errorf("dbus, notify-send and kdialog not found")
}

func zenitySelectDialogCmd(title string, isFile bool) models.Command {
	cmd := getDefaultCmd()
	cmd.IsAsync = false
	cmd.Cmd = fmt.Sprintf(`zenity --file-selection %s --title="%s"`, logic.Ternary(isFile, "", "--directory"), title)
	return cmd
}

func kdialogSelectDialogCmd(title string, isFile bool) models.Command {
	cmd := getDefaultCmd()
	cmd.IsAsync = false
	cmd.Cmd = fmt.Sprintf(`kdialog %s %s --title "%s"`, logic.Ternary(isFile, "--getopenfilename", "--getexistingdirectory"), system.HomeDir(), title)
	return cmd
}

func selectDialog(title string, isFile bool) models.Response[string] {
	response := models.NewResponse[string]()
	data, errZenity := exe.Exec(zenitySelectDialogCmd(title, isFile))
	if errZenity != nil {
		data, errKdialog := exe.Exec(kdialogSelectDialogCmd(title, isFile))
		if errKdialog != nil {
			response.Error = fmt.Errorf("zenity: %w; kdialog: %w", errZenity, errKdialog)
		} else {
			response.Data = data
		}
	} else {
		response.Data = data
	}
	return *response
}

func zenityDialogCmd(title string, message string, dialogType dialogType) models.Command {
	var flag string
	switch dialogType {
	case typeWarning:
		flag = "--warning"
	case typeError:
		flag = "--error"
	default:
		flag = "--info"
	}
	cmd := getDefaultCmd()
	cmd.Cmd = fmt.Sprintf(`zenity %s --title="%s" --text="%s"`, flag, title, message)
	return cmd
}

func kdialogDialogCmd(title string, message string, dialogType dialogType) models.Command {
	var flag string
	switch dialogType {
	case typeWarning:
		flag = "--warning"
	case typeError:
		flag = "--error"
	default:
		flag = "--info"
	}
	cmd := getDefaultCmd()
	cmd.Cmd = fmt.Sprintf(`kdialog --title "%s" %s "%s"`, title, flag, message)
	return cmd
}

func dialogBox(title string, message string, dialogType dialogType) error {
	if errZenity := exe.ExecRealTime(zenityDialogCmd(title, message, dialogType)); errZenity != nil {
		if errKdialog := exe.ExecRealTime(kdialogDialogCmd(title, message, dialogType)); errKdialog != nil {
			return fmt.Errorf("zenity: %w; kdialog: %w", errZenity, errKdialog)
		}
	}
	return nil
}
