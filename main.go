package main

import "github.com/sgrep/sgrep/sgreplib"
import "flag"
import "os"
import "os/exec"
import "log"
import "io"
import "path/filepath"

// libraries used to check for coloring
import "syscall"
import "golang.org/x/crypto/ssh/terminal"


// FIXME: probably more generic ways to do this (eg., for windows)
const GREP_BIN_PATH string = "grep"

// Arguments passed for whether or not to colorize output
const COLORIZE_AUTO = "auto"
const COLORIZE_ALWAYS = "always"
const COLORIZE_NEVER = "never"

/**
Basic structure:

  1) Check arguments and based on arguments, create directory objects

  2) Directory objects read through all .sgrep files and produce a
  list of files that we should grep over.

  3) Exec a grep command with the list of files generated from
  directory objects in #2, then redirect grep's stdout and stderr this
  process's stdout and stderr.
*/
func main() {
	args := parseArgs()

	// determine whether or not to put ansi codes to stdout.
	colorArg := "--color=never"
	if args.shouldColorize {
		colorArg = "--color=always"
	}

	// a list of files to actually exec over
	var filesToGrepOver []string

	// List of files passed in at command line by user to sgrep
	// over.  For instance, issuing "sgrep something *txt", in a
	// directory with files a.txt and b.txt in it would produce
	// filesToCheckWhetherToGrepOver with a.txt and b.txt.  We
	// later produce directory objects that determine whether to
	// filter the files when passing them into grep exec.
	var filesToCheckWhetherToGrepOver []string

	// for all files and folders passed into command line:
	//   1) build directory objects from folders
	//   2) copy all passed in files into filesToCheckWhetherToGrepOver
	for _, toGrepOver := range args.whereToGrep {
		file, err := os.Open(toGrepOver)
		if err != nil {
			log.Fatal("Could not open file " + toGrepOver)
		}
		defer file.Close()
		fi, err := file.Stat()
		if err != nil {
			log.Fatal("Could not stat file " + toGrepOver)
		}

		if fi.IsDir() {
			// if it's a directory, then read through all
			// folders and check for subdirectories to
			// recursively grep over.
			if !args.recursive {
				log.Fatal(
					"Specifying directory named " +
						toGrepOver +
						" without recursive flag, -r.")
			}
			dir := sgreplib.GenerateSgrepDirectories(toGrepOver)
			dirFiles := dir.ListNonRuleFilteredFiles()
			filesToGrepOver = append(filesToGrepOver, dirFiles...)
		} else {
			// later will run .sgrep filters over these
			// lists via call to
			// sgreplib.NonFilteredFilesToGrepOver
			filesToCheckWhetherToGrepOver = append(
				filesToCheckWhetherToGrepOver, toGrepOver)
		}
	}

	// run .sgrep filter over all command line argument filenames
	// pass in
	nonFilteredFiles := sgreplib.NonFilteredFilesToGrepOver(
		filesToCheckWhetherToGrepOver)
	filesToGrepOver = append(filesToGrepOver, nonFilteredFiles...)

	// Constuct args to grep call
	var argArray []string
	argArray = append(argArray, colorArg, args.whatToGrepFor)
	argArray = append(argArray, filesToGrepOver...)

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
	colorizeArgPtr := flag.String(
		"color", COLORIZE_AUTO,
		"Should add colors to stdout.  Options: "+
			"'"+COLORIZE_AUTO+",' "+
			"'"+COLORIZE_ALWAYS+",' or "+
			"'"+COLORIZE_NEVER)
	flag.Parse()

	// recursive arg
	toReturn.recursive = *recursiveArgPtr

	// check if should colorize output
	colorizeArg := *colorizeArgPtr
	if colorizeArg == COLORIZE_AUTO {

		if terminal.IsTerminal(syscall.Stdout) {
			// grep explicitly checks whether the term env
			// variable is set to dumb, and, if it is not,
			// we do not colorize
			termEnv := os.Getenv("TERM")
			if (termEnv == "dumb") {
				toReturn.shouldColorize = false
			} else {
				toReturn.shouldColorize = true
			}
		} else {
			toReturn.shouldColorize = false
		}
	} else if colorizeArg == COLORIZE_NEVER {
		toReturn.shouldColorize = false
	} else if colorizeArg == COLORIZE_ALWAYS {
		toReturn.shouldColorize = true
	} else {
		log.Fatal("Unknown arg to color: " + colorizeArg)
	}

	flagArgs := flag.Args()
	if len(flagArgs) == 0 {
		log.Fatal(
			"Currently, sgrep requires a single argument: " +
				"what to grep for.")
	}

	toReturn.whatToGrepFor = flagArgs[0]

	// construct absolute pathnames for what to grep over
	for _, whereToGrep := range flagArgs[1:] {
		absPath, err := filepath.Abs(whereToGrep)
		if err != nil {
			log.Fatal("Error creating abs path for file")
		}
		toReturn.whereToGrep = append(toReturn.whereToGrep, absPath)
	}
	return toReturn
}

type SgrepArgs struct {
	recursive      bool
	whatToGrepFor  string
	whereToGrep    []string
	shouldColorize bool
}
