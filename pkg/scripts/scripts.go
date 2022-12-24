package scripts

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/shikaan/shmux/pkg/exceptions"
)

const TEMP_SCRIPT_FILE = "shmux"

// A script is a slice of lines representing each a LOC
type Script = []string

// Executes a script in a given shell
// Arguments are positional arguments ($1, $2...) which can be used in the script as replacement
func RunScript(script Script, scriptName string, shell string, arguments []string) (string, error) {
	path := getTempScriptPath()
	os.RemoveAll(path)

	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileContent := strings.Join(script, "\n")
	file.WriteString(replaceArguments(fileContent, arguments, scriptName))

	out, err := exec.Command(shell, path).Output()
	if err != nil {
		exceptions.HandleScriptError(scriptName, err, string(out))
		return "", err
	}

	return string(out), nil
}

// Parses the provided file, retuning the Script whose name is the provided one
func ReadScript(scriptName string, file io.Reader) (Script, error) {
	lines := []string{}
	availableScripts := []string{}
	scanner := bufio.NewScanner(file)

	shouldCollect := false
	firstLine := true
	tabLength := 0
	for scanner.Scan() {
		line := scanner.Text()
		isScriptLine, match := readScript(line)

		if isScriptLine {
			availableScripts = append(availableScripts, match)

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
			// Skip empty lines
			if line == "" {
				continue
			}
			// Identifies the intendeation for this script by looking
			// at the whitespace prepending the first line of the script
			if firstLine {
				tabLength = len(line) - len(strings.TrimSpace(line))
				firstLine = false
			}

			lines = append(lines, line[tabLength:])
		}
	}

	// If shouldCollect was never toggled, the script was not found
	if !shouldCollect {
		return nil, fmt.Errorf("could not find \"%s\". Available scripts: %s", scriptName, strings.Join(availableScripts, ", "))
	}

	return lines, nil
}

func MakeHelp(file io.Reader) string {
	availableScripts := []string{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		isScriptLine, match := readScript(line)

		if isScriptLine {
			availableScripts = append(availableScripts, match)
		}
	}

	return fmt.Sprintf(`usage: shmux [-config <path>] [-shell <path>] <script> -- [arguments ...]

Available scripts: %s
Run 'shmux -h' for details.
`, strings.Join(availableScripts, ", "))
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

// Returns the path of the temporary script we use for execution
func getTempScriptPath() string {
	return filepath.Join(os.TempDir(), TEMP_SCRIPT_FILE)
}

const SCRIPT_IDENTIFIER_REGEXP = "^(\\S+):(.*)"

// Tries identifying lines including scripts, returning if it's a match and what that is
func readScript(line string) (isScriptLine bool, match string) {
	r, _ := regexp.Compile(SCRIPT_IDENTIFIER_REGEXP)
	submatch := r.FindStringSubmatch(line)

	if len(submatch) > 1 {
		return true, submatch[1]
	}

	return false, ""
}

// Caps a slice to a certain length
func cap(slice []string, n int) []string {
	if len(slice) < n {
		return slice
	}

	return slice[:n]
}
