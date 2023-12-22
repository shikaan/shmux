package scripts

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func MakeHelp(file io.Reader) string {
	availableScripts := []string{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		isScriptLine, match, _ := readScript(line)

		if isScriptLine {
			availableScripts = append(availableScripts, match)
		}
	}

	return fmt.Sprintf(`usage: shmux [-config <path>] [-shell <path>] <script> -- [arguments ...]

Available scripts: %s
Run 'shmux -h' for details.
`, strings.Join(availableScripts, ", "))
}