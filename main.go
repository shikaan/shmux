package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const FILE_NAME = "shmux.conf"
var TEMP = filepath.Join(os.TempDir(), "shmux")

func readCommand(line string) (bool, string) {
  r, _ := regexp.Compile("(.+):")
  match := r.FindStringSubmatch(line)

  if len(match) > 1 {
    return true, match[1]
  }

  return false, ""
}

func runCommandLines(lines []string) {
    os.RemoveAll(TEMP)
    file, _ := os.Create(TEMP)
    defer file.Close()
    file.WriteString("#!/bin/sh\n")
    file.WriteString(strings.Join(lines, "\n"))
    _ = os.Chdir(TEMP)
    out, _ := exec.Command("sh", TEMP).Output()
    println(string(out))
}

func main() { 
  file, err := os.Open(FILE_NAME);
  command := os.Args[1]
  
  println("command: ", command)

  if err != nil {
    log.Fatal(err)
    os.Exit(1)
  }

  defer file.Close()
  
  lines := []string{}
  scanner := bufio.NewScanner(file)

  isCollecting := false
  for scanner.Scan() {
    line := scanner.Text()
    isCommand, match := readCommand(line)

    // finds the first occurrence of the command
    if isCommand && !isCollecting && match == command {
      isCollecting = true
      continue
    }

    // after finding the command stop on the next
    if isCommand && isCollecting {
      break
    }

    if isCollecting {
      lines = append(lines, line)
    }
  }

  runCommandLines(lines)
}
