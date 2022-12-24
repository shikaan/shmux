package exceptions

import (
	"fmt"
	"os"
)

func HandleException(e error) {
	if e != nil {
		os.Stderr.WriteString(fmt.Sprintf("shmux: %s\n", e.Error()))
		os.Exit(1)
	}
}

func HandleScriptError(scriptName string, e error, output string) {
	if e != nil {
		os.Stderr.WriteString(fmt.Sprintf("shmux: script \"%s\" returned non-zero status code.\n", scriptName))
		os.Stderr.WriteString(fmt.Sprintf("Error: %s\n", e.Error()))
		os.Exit(1)
	}
}
