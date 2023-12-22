package exceptions

import (
	"fmt"
	"os"
)

func HandleException(e error, status int) {
	if e != nil {
		os.Stderr.WriteString(fmt.Sprintf("shmux: %s\n", e.Error()))
		os.Exit(status)
	}
}
