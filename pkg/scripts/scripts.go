package scripts

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const TEMP_SCRIPT_FILE = "shmux"
const SCRIPT_IDENTIFIER_REGEXP = "^(.+):$"

// A script is a slice of lines representing each a LOC
type Script = []string

// Executes a script in a given shell
// Arguments are positional arguments ($1, $2...) which can be used in the script as replacement
func RunScript(script Script, shell string, arguments []string) (string, error) {
	path := getTempScriptPath()
	os.RemoveAll(path)

	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	file.WriteString(strings.Join(script, "\n"))

	out, err := exec.Command(shell, path).Output()
	if err != nil {
		println(shell, path, err.Error())
		return "", err
	}

	return string(out), nil
}

func ReadScript(scriptName string, file *os.File) (Script, error) {
	lines := []string{}
	scanner := bufio.NewScanner(file)

	shouldCollect := false
	for scanner.Scan() {
		line := scanner.Text()
		isScriptLine, match := readScript(line)

		if isScriptLine {
			// Switch on collection, if we match the given script name.
			// This allows collecting lines from next line on.
			if !shouldCollect && match == scriptName {
				shouldCollect = true
				continue
			}

			// If collection is on and we find another script, we stop
			// collection. We don't check for EOF, as in that case the
			// loop would stop anyway.
			if shouldCollect && isScriptLine {
				break
			}
		}

		if shouldCollect {
			lines = append(lines, line)
		}
	}

	// If shouldCollect was never toggled, the script was not found
	if !shouldCollect {
		return nil, fmt.Errorf("could not find script \"%s\"", scriptName)
	}

	return lines, nil
}

// Returns the path of the temporary script we use for execution
func getTempScriptPath() string {
	return filepath.Join(os.TempDir(), TEMP_SCRIPT_FILE)
}

// Tries identifying lines including scripts, returning if it's a match and what that is
func readScript(line string) (isScriptLine bool, match string) {
	r, _ := regexp.Compile(SCRIPT_IDENTIFIER_REGEXP)
	submatch := r.FindStringSubmatch(line)

	if len(submatch) > 1 {
		return true, submatch[1]
	}

	return false, ""
}
