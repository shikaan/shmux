package arguments

import (
	"flag"
)

const DEFAULT_FILE_NAME = "shmux.conf"
const DEFAULT_SHELL = "/bin/sh"

const ARGUMENT_SEPARATOR = "--"

func Parse() (shell string, config string, scriptName string, arguments []string) {
	shell = *flag.String("shell", DEFAULT_SHELL, "Shell to be used to run the scripts")
	config = *flag.String("config", DEFAULT_FILE_NAME, "Configuration file")

	flag.Parse()

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
