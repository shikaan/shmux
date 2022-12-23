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
	exceptions.HandleError(err)
	defer file.Close()

	script, err := scripts.ReadScript(scriptName, file)
	exceptions.HandleError(err)

	output, err := scripts.RunScript(script, shell, arguments)
	exceptions.HandleError(err)

	fmt.Print(output)
}
