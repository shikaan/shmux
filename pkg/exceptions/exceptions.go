package exceptions

import (
	"fmt"
	"os"
)

func HandleException(e error) {
	if e != nil {
		os.Stderr.WriteString(fmt.Sprintf("shmux: %s\n", e.Error()))
		os.Exit(1)
	}
}

func HandleError(e error) {
	if e != nil {
		log.Fatalf("shmux: %s\n", e.Error())
	}
}
