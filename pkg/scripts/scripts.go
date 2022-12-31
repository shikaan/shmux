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
type Script struct {
	Name        string
	Lines       []string
	Interpreter string
	Options     []string
}

// Executes a script in a given shell
// Arguments are positional arguments ($1, $2...) which can be used in the script as replacement
func RunScript(script *Script, arguments []string) (string, error) {
	path := getTempScriptPath()
	os.RemoveAll(path)

	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileContent := strings.Join(script.Lines, "\n")
	file.WriteString(replaceArguments(fileContent, arguments, script.Name))

	out, err := exec.Command(script.Interpreter, append(script.Options, path)...).Output()
	if err != nil {
		exceptions.HandleScriptError(script.Name, err, string(out))
		return "", err
	}

	return string(out), nil
}

// Parses the provided file, retuning the Script whose name is the provided one
func ReadScript(scriptName string, shell string, file io.Reader) (*Script, error) {
	script := &Script{Name: scriptName, Interpreter: shell}
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

			if firstLine {
				// Identifies the intendeation for this script by looking
				// at the whitespace prepending the first line of the script
				tabLength = len(line) - len(strings.TrimSpace(line))
				firstLine = false

				// Overrides the interpreter to be used to execute the script,
				// if it's provided as hashbang
				isShellLine, interpreter, options := readShebang(line)

				if isShellLine {
					script.Interpreter = interpreter
					script.Options = options
					continue
				}
			}

			script.Lines = append(script.Lines, line[tabLength:])
		}
	}

	// If shouldCollect was never toggled, the script was not found
	if !shouldCollect {
		return nil, fmt.Errorf("could not find \"%s\". Available scripts: %s", scriptName, strings.Join(availableScripts, ", "))
	}

	return script, nil
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
	scriptName := get(submatch, 1)

	return scriptName != "", scriptName
}

const SHEBANG_REGEXP = "#!\\s?(\\S+)\\s?(.*)"

// Extracts the interpreter and the options from a shebang line
func readShebang(line string) (isShebangLine bool, interpreter string, options []string) {
	r, _ := regexp.Compile(SHEBANG_REGEXP)
	submatch := r.FindStringSubmatch(line)

	interpreter = get(submatch, 1)
	stringOptions := strings.TrimSpace(get(submatch, 2))

	if stringOptions != "" {
		options = strings.Split(stringOptions, " ")
	}

	return interpreter != "", interpreter, options
}

// Caps a slice to a certain length
func cap(slice []string, n int) []string {
	if len(slice) < n {
		return slice
	}

	return slice[:n]
}

func get(slice []string, index int) string {
	if len(slice) > index {
		return slice[index]
	}

	return ""
}
