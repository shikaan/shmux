package exceptions

import "log"

func HandleError(e error) {
	if e != nil {
		log.Fatalf("shmux: %s\n", e.Error())
	}
}
