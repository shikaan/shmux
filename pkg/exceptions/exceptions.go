package exceptions

import (
	"fmt"
	"os"
	"os/exec"
)

func HandleException(e error) {
	if e != nil {
		os.Stderr.WriteString(fmt.Sprintf("shmux: %s\n", e.Error()))
		os.Exit(1)
	}
}

func HandleScriptError(scriptName string, e error) {
	if e != nil {
		status := 1
		if exitError, ok := e.(*exec.ExitError); ok {
			status = exitError.ExitCode()
		}

		os.Stderr.WriteString(fmt.Sprintf("shmux: script \"%s\" exited with code %d\n", scriptName, status))
		os.Exit(status)
	}
}
