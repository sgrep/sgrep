package main

import "github.com/sgrep/sgrep/sgreplib"
import "fmt"
import "os"

import "os/exec"
import "log"
import "io"

// FIXME: probably more generic ways to do this (eg., for windows)
const GREP_BIN_PATH string = "grep"

func main() {
	if len(os.Args) != 2 {
		fmt.Println(
			"Currently, sgrep requires a single argument: " +
				"what to grep for.")
		return
	}
	whatToGrepFor := os.Args[1]
	
	currWorkingDir, err := os.Getwd()
	if err != nil {
		panic("Could not find current working directory")
	}
	dir := sgreplib.GenerateSgrepDirectories (currWorkingDir)
	filesToGrepOver := dir.ListNonRuleFilteredFiles()
	
	var argArray [] string
	argArray = append(argArray, whatToGrepFor)
	argArray = append(argArray, filesToGrepOver... )
	
	// execute grep command
	cmd := exec.Command(GREP_BIN_PATH, argArray...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
	cmd.Wait()
}
