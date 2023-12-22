package scripts

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// To prevent infinite recursion, we limit the nesting level of script dependencies
const MAXIMUM_NESTING_LEVEL = 2

// A script is a slice of lines representing each a LOC
type Script struct {
	Name        string
	Lines       []string
	Interpreter string
	Options     []string
	Dependencies []Script
}

// Parses the provided file, retuning the Script whose name is the provided one
func ReadScript(scriptName string, shell string, fileContent []byte, currentNestingLevel uint) (script *Script, err error) {
	script = &Script{Name: scriptName, Interpreter: shell, Dependencies: []Script{}}
	lines := strings.Split(string(fileContent), "\n")

	availableScripts := []string{}
	shouldCollect := false
	firstLine := true
	tabLength := 0

	for _, line := range lines {
		isScriptLine, match, dependencies := readScript(line)

		if isScriptLine {
			availableScripts = append(availableScripts, match)

			// Switch on collection, if we match the given script name.
			// This allows collecting lines from next line on.
			if !shouldCollect && match == scriptName {
				shouldCollect = true

				script.Dependencies, err = getDependencyScripts(dependencies, match, currentNestingLevel, shell, fileContent)
				if err != nil { 
					return nil, err
				}

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

	return
}

func getDependencyScripts(dependencies []string, match string, currentNestingLevel uint, shell string, fileContent []byte) (deps []Script, err error) {
	for _, dependency := range dependencies {

		if dependency == match {
			return nil, fmt.Errorf("script \"%s\" has a circular dependency on itself", match)
		}

		if currentNestingLevel >= MAXIMUM_NESTING_LEVEL {
			return nil, fmt.Errorf("script \"%s\" has a dependency on \"%s\" which exceeds the maximum nesting level of %d", match, dependency, MAXIMUM_NESTING_LEVEL)
		}

		var dependencyScript *Script; 
		dependencyScript, err = ReadScript(dependency, shell, fileContent, currentNestingLevel+1)
		if err != nil { return }

		deps = append(deps, *dependencyScript)
	}

	return
}

const SCRIPT_IDENTIFIER_REGEXP = `^(\S+):\s*(.*)$`

// Tries identifying lines including scripts, returning if it's a match and what that is
func readScript(line string) (isScriptLine bool, script string, dependencies []string) {
	r, _ := regexp.Compile(SCRIPT_IDENTIFIER_REGEXP)
	submatch := r.FindStringSubmatch(line)
	
	if submatch == nil { return }

	script = submatch[1]
	isScriptLine = script != "" 
	
	if isScriptLine {
		dependencies = slices.Compact(strings.Fields(submatch[2]))
	}

	return
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
