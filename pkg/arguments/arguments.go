package arguments

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const ARGUMENT_SEPARATOR = "--"
const CONFIGURATION_GLOB = "shmuxfile.*"
const ENVIRONMENT_CONFIGURATION = "SHMUX_CONFIG"
const ENVIRONMENT_SHELL = "SHMUX_SHELL"
const HELP_SCRIPT = "$$$$___HELP___$$$$"
const HELP_TEXT = `usage: shmux [-config <path>] [-shell <path>] <script> -- [arguments ...]

shmux is a utility to run multiple scripts from one file. Scripts can be written in (almost) any language and they don't need to be in the same language.

The scripts are defined in a configuration file called 'shmuxfile'. Similarly to a Makefile for GNU Make, this file serves as a manifest for the scripts to be run.

More information is available at https://github.com/shikaan/shmux.

Available flags:

`

func Parse() (shell string, config string, scriptName string, arguments []string, err error) {
	flag.Usage = func() {
		fmt.Print(HELP_TEXT)
		flag.PrintDefaults()
	}

	configFlag := flag.String("config", "", "Configuration file path. It falls back to the closest 'shmuxfile.*' available.")
	shellFlag := flag.String("shell", "", "Interpreter used to run the scripts. It defaults to the current $SHELL.")

	flag.Parse()

	shell = oneOf(*shellFlag, os.Getenv(ENVIRONMENT_SHELL), os.Getenv("SHELL"))

	config, err = getConfigurationLocation(oneOf(*configFlag, os.Getenv(ENVIRONMENT_CONFIGURATION)))
	if err != nil {
		return
	}

	err = validateShell(shell)
	if err != nil {
		return
	}

	scriptName = oneOf(flag.Arg(0), HELP_SCRIPT)
	arguments = getAdditionalArguments(flag.Args())

	return
}

// Returns the index of the given element in the slice or -1 if missing
func findIndex(slice []string, element string) int {
	for k, v := range slice {
		if element == v {
			return k
		}
	}

	return -1
}

// Returns the first non-zero value item in the list
func oneOf(items ...string) string {
	for _, i := range items {
		if i != "" {
			return i
		}
	}

	return ""
}

func validateShell(shell string) error {
	info, err := os.Stat(shell)
	if err != nil {
		return err
	}

	if !canOwnerExec(info.Mode()) {
		return fmt.Errorf("cannot execute \"%s\"", shell)
	}

	return nil
}

func canOwnerExec(mode os.FileMode) bool {
	return mode&0100 != 0
}

// Looks up and returns the absolute path to the configuration file to be found
// in the current working directory or the parent folders
func getConfigurationLocation(configFileName string) (string, error) {
	// Utilise provided configuration if any, else do globmatching
	searchTerm := oneOf(configFileName, CONFIGURATION_GLOB)
	workingDirectory, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Ensure folders are all in the same shape before the comparison
	currentDirectory, err := filepath.Abs(workingDirectory)
	if err != nil {
		return "", err
	}

	for {
		// Error can only be bad patters, hence ignored
		matches, _ := filepath.Glob(filepath.Join(currentDirectory, searchTerm))

		if len(matches) > 0 {
			return matches[0], nil
		}

		newDirectory := filepath.Dir(currentDirectory)
		isRoot := currentDirectory == newDirectory

		if isRoot {
			return "", fmt.Errorf("cannot find \"%s\" here (%s) or in any parent folder", searchTerm, workingDirectory)
		}

		currentDirectory = newDirectory
	}
}

// Returns arguments provided after ARGUMENT_SEPARATOR
func getAdditionalArguments(args []string) []string {
	doubleDashIndex := findIndex(args, ARGUMENT_SEPARATOR)

	if doubleDashIndex == -1 {
		return []string{}
	}

	return args[(doubleDashIndex + 1):]
}
