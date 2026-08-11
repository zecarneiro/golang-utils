//go:build windows

package ui

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/exe"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/shell"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
)

var (
	scriptTemplate = `[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] > $null
$Template = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastImageAndText02)
$RawXml = [xml] $Template.GetXml()
($RawXml.toast.visual.binding.text|Where-Object {$_.id -eq "1"}).AppendChild($RawXml.CreateTextNode(%s)) > $null
($RawXml.toast.visual.binding.text|Where-Object {$_.id -eq "2"}).AppendChild($RawXml.CreateTextNode(%s)) > $null
%s
$SerializedXml = New-Object Windows.Data.Xml.Dom.XmlDocument
$SerializedXml.LoadXml($RawXml.OuterXml)
Write-Host $SerializedXml.ToString()
$Toast = [Windows.UI.Notifications.ToastNotification]::new($SerializedXml)
$Toast.Tag = "%s"
$Toast.Group = "%s"
$Toast.ExpirationTime = [DateTimeOffset]::Now.AddMinutes(1)
$Notifier = [Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("%s")
$Notifier.Show($Toast);
`
	scriptSelectDialogTemplate = `Add-Type -AssemblyName System.Windows.Forms
$dialog = New-Object System.Windows.Forms.FolderBrowserDialog
$dialog.Description = '%s'
$dialog.ShowNewFolderButton = $true
$dialog.InitialDirectory = [Environment]::GetFolderPath('%s')
if ($dialog.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) {
	%s
}
`
	scriptDialogBoxTemplate = `Add-Type -AssemblyName System.Windows.Forms
[System.Windows.Forms.MessageBox]::Show('%s', '%s', [System.Windows.Forms.MessageBoxButtons]::OK, [System.Windows.Forms.MessageBoxIcon]::%s)
`
)

func notify(title string, message string, icon string) error {
	iconScriptData := ""
	if !str.IsEmpty(icon) {
		iconScriptData = fmt.Sprintf(`($RawXml.toast.visual.binding.image|Where-Object {$_.id -eq "1"}).SetAttribute('src', %s) > $null`, icon)
	}
	script := fmt.Sprintf(scriptTemplate, title, message, iconScriptData, common.GetAppId(), common.GetAppId(), common.GetAppId())
	return exe.ExecRealTime(models.Command{Cmd: shell.GetPowershellCmd(), Args: []string{"-NoProfile", "-NonInteractive", "-Command", script}, UseShell: false, Verbose: verbose, IsAsync: true})
}

func selectDialog(title string, isFile bool) models.Response[string] {
	response := models.NewResponse[string]()
	typeScriptData := logic.Ternary(isFile, "Write-Output $dialog.FileName", "Write-Output $dialog.SelectedPath")
	script := fmt.Sprintf(scriptSelectDialogTemplate, title, system.HomeDir(), typeScriptData)
	cmd := models.Command{Cmd: shell.GetPowershellCmd(), Args: []string{"-NoProfile", "-NonInteractive", "-Command", script}, UseShell: false, Verbose: verbose, IsAsync: false}
	response.Data, response.Error = exe.Exec(cmd)
	return *response
}

func dialogBox(title string, message string, dialogType dialogType) error {
	var icon string
	switch dialogType {
	case typeWarning:
		icon = "Warning"
	case typeError:
		icon = "Error"
	default:
		icon = "Information"
	}
	script := fmt.Sprintf(scriptDialogBoxTemplate, message, title, icon)
	return exe.ExecRealTime(models.Command{Cmd: shell.GetPowershellCmd(), Args: []string{"-NoProfile", "-NonInteractive", "-Command", script}, UseShell: false, Verbose: verbose, IsAsync: true})
}
