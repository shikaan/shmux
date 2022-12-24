package main

import (
	"fmt"
	"os"

	"github.com/shikaan/shmux/pkg/arguments"
	"github.com/shikaan/shmux/pkg/exceptions"
	"github.com/shikaan/shmux/pkg/scripts"
)

func main() {
	shell, config, scriptName, arguments := arguments.Parse()

	file, err := os.Open(config)
	exceptions.HandleException(err)
	defer file.Close()

	script, err := scripts.ReadScript(scriptName, file)
	exceptions.HandleException(err)

	output, err := scripts.RunScript(script, scriptName, shell, arguments)
	exceptions.HandleException(err)

	fmt.Print(output)
}
