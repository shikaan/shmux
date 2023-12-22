package main

import (
	"fmt"
	"io"
	"os"

	"github.com/shikaan/shmux/pkg/arguments"
	"github.com/shikaan/shmux/pkg/exceptions"
	"github.com/shikaan/shmux/pkg/scripts"
)

const MAX_FILE_SIZE = 1<<20; // 1MB

func main() {
	shell, config, scriptName, args, err := arguments.Parse()
	exceptions.HandleException(err, 1)

	file, err := os.Open(config)
	exceptions.HandleException(err, 1)
	defer file.Close()

	if scriptName == arguments.HELP_SCRIPT {
		fmt.Print(scripts.MakeHelp(file))
		return
	}

	limitedReader := io.LimitReader(file, MAX_FILE_SIZE)
	content, err := io.ReadAll(limitedReader)
	exceptions.HandleException(err, 1)

	script, err := scripts.ReadScript(scriptName, shell, content, 0)
	exceptions.HandleException(err, 1)

	status, err := scripts.RunScript(script, args)
	exceptions.HandleException(err, status)
}
