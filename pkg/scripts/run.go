package scripts

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// The file we use to store the script we want to execute
var TEMP_SCRIPT_FILE = ""
func init() {
	TEMP_SCRIPT_FILE = fmt.Sprintf("shmux_%d", time.Now().UnixNano())
}

// Creates an exec.Cmd that executes a script in a given shell
// Arguments are positional arguments ($1, $2...) which can be used in the script as replacement
func RunScript(script *Script, arguments []string) (status int, err error) {
	for _, dependency := range script.Dependencies {
		status, err = RunScript(&dependency, arguments)
		if err != nil { return }
	}
	path := getTempScriptPath(script.Name)

	file, err := os.Create(path)
	if err != nil {
		return 1, err
	}
	defer file.Close()
	defer os.RemoveAll(path)

	fileContent := strings.Join(script.Lines, "\n")
	file.WriteString(replaceArguments(fileContent, arguments, script.Name))

	cmd := exec.Command(script.Interpreter, append(script.Options, path)...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		status = 1
		if exitError, ok := err.(*exec.ExitError); ok {
			status = exitError.ExitCode()
			err = fmt.Errorf("script \"%s\" exited with code %d", script.Name, status)
		}
	}

	return
}

// Returns the path of the temporary script we use for execution
func getTempScriptPath(scriptName string) string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("%s_%s", TEMP_SCRIPT_FILE, scriptName))
}

// Replaces positional arguments in the arguments slice ($1..$9) and
// $@ with script name in 'content
func replaceArguments(content string, arguments []string, scriptName string) string {
	result := content

	for i, arg := range cap(arguments, 9) {
		placeholder := fmt.Sprintf("$%d", i+1)
		result = strings.ReplaceAll(result, placeholder, arg)
	}

	return strings.ReplaceAll(result, "$@", scriptName)
}