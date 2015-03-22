package main

import "github.com/sgrep/sgrep/sgreplib"
import "flag"
import "os"
import "os/exec"
import "log"
import "io"


// FIXME: probably more generic ways to do this (eg., for windows)
const GREP_BIN_PATH string = "grep"

func main() {
	args := parseArgs()
	
	currWorkingDir, err := os.Getwd()
	if err != nil {
		panic("Could not find current working directory")
	}
	dir := sgreplib.GenerateSgrepDirectories (currWorkingDir)
	filesToGrepOver := dir.ListNonRuleFilteredFiles()
	
	var argArray [] string
	argArray = append(argArray, args.whatToGrepFor)
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

/**
Read the command line args passed into this binary.
*/
func parseArgs() *SgrepArgs {
	toReturn := new(SgrepArgs)
	
	recursiveArgPtr := flag.Bool("r", false, "Recursive")
	flag.Parse()

	toReturn.recursive = *recursiveArgPtr
	flagArgs := flag.Args()
	if len(flagArgs) == 0 {
		log.Fatal(
			"Currently, sgrep requires a single argument: " +
				"what to grep for.")
	}

	toReturn.whatToGrepFor = flagArgs[0]
	toReturn.whereToGrep = flagArgs[1:]
	return toReturn
}

type SgrepArgs struct {
	recursive bool
	whatToGrepFor string
	whereToGrep []string
}
