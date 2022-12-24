package arguments

import (
	"flag"
	"fmt"
	"os"
)

const DEFAULT_CONFIGURATION = "shmux.sh"
const CONFIGURATION_ENVIRONMENT = "SHMUX_CONFIG"
const DEFAULT_SHELL = "/bin/sh"
const SHELL_ENVIRONMENT = "SHMUX_SHELL"
const ARGUMENT_SEPARATOR = "--"
const HELP_SCRIPT = "$$$$___HELP___$$$$"
const HELP_TEXT = `usage: shmux [-config <path>] [-shell <path>] <script> -- [arguments ...]

shmux is a utility to run multiple scripts from one file. Scripts can be written in (almost) any language and they don't need to be in the same language.

The scripts are defined in a configuration file called shmuxfile. Similarly to a Makefile for GNU Make, this file serves as a manifest for the scripts to be run.

More information is available at https://github.com/shikaan/shmux.

Available flags:

`

func Parse() (shell string, config string, scriptName string, arguments []string) {
	flag.Usage = func() {
		fmt.Print(HELP_TEXT)
		flag.PrintDefaults()
	}

	configFlag := flag.String("config", "", fmt.Sprintf("Configuration file path. It falls back to the SHMUX_SHELL environment variable. (default %s)", DEFAULT_CONFIGURATION))
	shellFlag := flag.String("shell", "", fmt.Sprintf("Shell to be used to run the scripts. It falls back to the SHMUX_SHELL environment variable. (default %s)", DEFAULT_SHELL))

	flag.Parse()

	shell = oneOf(*shellFlag, os.Getenv(SHELL_ENVIRONMENT), DEFAULT_SHELL)
	config = oneOf(*configFlag, os.Getenv(CONFIGURATION_ENVIRONMENT), DEFAULT_CONFIGURATION)
	scriptName = oneOf(flag.Arg(0), HELP_SCRIPT)

	doubleDashIndex := findIndex(flag.Args(), ARGUMENT_SEPARATOR)

	if doubleDashIndex == -1 {
		arguments = []string{}
	} else {
		arguments = flag.Args()[:doubleDashIndex]
	}

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
