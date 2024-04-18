package jnoronha_golangutils

import (
	"bufio"
	"fmt"
	"jnoronha_golangutils/entities"
	"log"
	"os"
	"os/exec"
	"strings"
)

func addShellCommand(commandInfo entities.CommandInfo) entities.CommandInfo {
	ValidateSystem()
	name := commandInfo.Cmd
	if IsWindows() {
		if !commandInfo.UsePowerShell {
			commandInfo.Cmd = "cmd.exe"
			commandInfo.Args = append([]string{"/c",  name}, commandInfo.Args...)
		} else {
			commandInfo.Cmd = "powershell.exe"
			commandInfo.Args = append([]string{name}, commandInfo.Args...)
		}
		
	} /* else if IsUnix() || IsDarwin() || IsLinux() {
		commandInfo.Cmd = "/bin/sh"
		commandInfo.Args = append([]string{"-c", name}, commandInfo.Args...)
	}*/
	return commandInfo
}

func Exec(commandInfo entities.CommandInfo) entities.Response[string] {
	var output []byte
	var cmd *exec.Cmd
	var err error
	var commandStr string = fmt.Sprintf("%s %s", commandInfo.Cmd, strings.Join(commandInfo.Args, " "))
	command := addShellCommand(commandInfo)
	cmd = exec.Command(command.Cmd, command.Args...)
	cmd.Dir = commandInfo.Cwd
	if commandInfo.Verbose {
		PromptLog(commandStr)
	}
	output, err = cmd.CombinedOutput()
	outputStr := string(output[:])
	if commandInfo.IsThrow && err != nil {
		log.Fatal(err)
	}
	if commandInfo.Verbose {
		fmt.Println(outputStr)
	}
	return entities.Response[string]{Data: outputStr, Error: err}
}

func ExecRealTime(commandInfo entities.CommandInfo) {
	var cmd *exec.Cmd
	var commandStr string = fmt.Sprintf("%s %s", commandInfo.Cmd, strings.Join(commandInfo.Args, " "))
	command := addShellCommand(commandInfo)
	cmd = exec.Command(command.Cmd, command.Args...)
	cmd.Dir = command.Cwd
	cmd.Env = commandInfo.EnvVars
	if commandInfo.Verbose {
		PromptLog(commandStr)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Start the command
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error starting command:", err)
		return
	}
}

func Which(cmd string) []string {
	commandInfo := entities.CommandInfo{Verbose: false, IsThrow: false}
	if IsWindows() {
		commandInfo.Cmd = "Get-Command " + cmd + " | Select-Object -ExpandProperty Definition"
		commandInfo.UsePowerShell = true
	} else {
		commandInfo.Cmd = "which " + cmd
	}
	response := Exec(commandInfo)
	return strings.Split(response.Data, SystemInfo().Eol);
}

func Confirm(message string, isNoDefault bool) bool {
	yesNoMsg := "[y/N]"
	if !isNoDefault {
		yesNoMsg = "[Y/n]"
	}
	fmt.Printf("%s %s?: ", message, yesNoMsg)
	var response string
	fmt.Scanln(&response)
	if response == "Y" || response == "y" {
		return true
	}
	return false
}

func ConfirmOrExit(message string, isNoDefault bool) bool {
	yesNoMsg := "[y/N/0(Exit)]"
	if !isNoDefault {
		yesNoMsg = "[Y/n/0(Exit)]"
	}
	fmt.Printf("%s %s?: ", message, yesNoMsg)
	var response string
	fmt.Scanln(&response)
	if response == "0" {
		os.Exit(0)
	}
	if response == "Y" || response == "y" {
		return true
	}
	return false
}

func WaitForAnyKeyPressed(message string) {
	LogLog(message, true)
	bufio.NewReader(os.Stdin).ReadBytes('\n') 
}
