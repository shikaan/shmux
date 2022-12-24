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

func Parse() (shell string, config string, scriptName string, arguments []string) {
	shellFlag := flag.String("s", "", fmt.Sprintf("Shell to be used to run the scripts (default %s)", DEFAULT_SHELL))
	configFlag := flag.String("c", "", fmt.Sprintf("Configuration file path (default %s)", DEFAULT_CONFIGURATION))

	flag.Parse()

	shell = oneOf(*shellFlag, os.Getenv(SHELL_ENVIRONMENT), DEFAULT_SHELL)
	config = oneOf(*configFlag, os.Getenv(CONFIGURATION_ENVIRONMENT), DEFAULT_CONFIGURATION)
	scriptName = flag.Arg(0)

	doubleDashIndex := index(flag.Args(), ARGUMENT_SEPARATOR)

	if doubleDashIndex == -1 {
		arguments = []string{}
	} else {
		arguments = flag.Args()[:doubleDashIndex]
	}

	return
}

func index(slice []string, element string) int {
	for k, v := range slice {
		if element == v {
			return k
		}
	}

	return -1
}

func oneOf(items ...string) string {
	for _, i := range items {
		if i != "" {
			return i
		}
	}

	return ""
}
