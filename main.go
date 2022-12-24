package main

import (
	"fmt"
	"os"

	"github.com/shikaan/shmux/pkg/arguments"
	"github.com/shikaan/shmux/pkg/exceptions"
	"github.com/shikaan/shmux/pkg/scripts"
)

func main() {
	shell, config, scriptName, args := arguments.Parse()

	file, err := os.Open(config)
	exceptions.HandleException(err)
	defer file.Close()

	if scriptName == arguments.HELP_SCRIPT {
		fmt.Print(scripts.MakeHelp(file))
		return
	}

	script, err := scripts.ReadScript(scriptName, file)
	exceptions.HandleException(err)

	output, err := scripts.RunScript(script, scriptName, shell, args)
	exceptions.HandleException(err)

	fmt.Print(output)
}
