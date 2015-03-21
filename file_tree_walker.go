package sgrep

import "io/ioutil"
import "path"
import "fmt"
import "os"

const SGREP_FILENAME string = ".sgrep"

type Directory struct {
	// The name of this directory relative to its parent directory.
	// Eg., for structure a/b/c/, would have a directory with name
	// a, name b, and name c.
	Name        string
	Files       []string
	Directories []*Directory
	Rules       []*Rule
}

func (dir *Directory) PrettyPrint() {
	dir.prettyPrintHelper(0)
}

/**
@returns true if dir's rule slice requires us to filter file. false otherwise.
*/
func (dir *Directory) shouldFilterFile(filename string) bool {
	for _, rule := range dir.Rules {
		if rule.FileFilterer(filename) {
			return true
		}
	}
	return false
}

func (dir *Directory) ListNonRuleFilteredFiles() []string {
	toReturn := make([]string, 0)

	// check files directly in dir (not including subdirectories)
	for _, filename := range dir.Files {
		absFilename := path.Join(dir.Name, filename)
		if !dir.shouldFilterFile(absFilename) {
			toReturn = append(toReturn, absFilename)
		}
	}

	// check subdirectories
	for _, subdir := range dir.Directories {
		subdirFileSlice := subdir.ListNonRuleFilteredFiles()

		for _, filename := range subdirFileSlice {
			absFilename := path.Join(dir.Name, filename)
			if !dir.shouldFilterFile(absFilename) {
				toReturn = append(toReturn, absFilename)
			}
		}
	}
	return toReturn
}

/**
  Returns all files as fully-qualified filename from directory dir.
*/
func (dir *Directory) ListFiles() []string {
	toReturn := make([]string, 0)

	for _, filename := range dir.Files {
		absFilename := path.Join(dir.Name, filename)
		toReturn = append(toReturn, absFilename)
	}

	for _, subdir := range dir.Directories {
		subdirFileSlice := subdir.ListFiles()

		for _, filename := range subdirFileSlice {
			absFilename := path.Join(dir.Name, filename)
			toReturn = append(toReturn, absFilename)
		}
	}
	return toReturn
}

func (dir *Directory) prettyPrintHelper(indentationLevel uint32) {

	indentStr := ""
	for i := uint32(0); i < indentationLevel; i++ {
		indentStr += "\t"
	}
	fmt.Println(indentStr + dir.Name + "/")

	for _, filename := range dir.Files {
		fmt.Println(indentStr + "\t" + filename)
	}

	for _, subDirectory := range dir.Directories {
		subDirectory.prettyPrintHelper(indentationLevel + 1)
	}
}

/**
  @param dirToWalk The directory that we want to list all
  subdirectories of.

  @returns A Directory struct containing all subfiles and folders.
*/
func WalkFolder(dirToWalk string) *Directory {
	dirContentsList, err := ioutil.ReadDir(dirToWalk)
	if err != nil {
		panic("IOError when reading directory " + dirToWalk)
	}

	root := createDirectoryFromSgrep(dirToWalk)
	for _, fileInfo := range dirContentsList {
		absPath := path.Join(dirToWalk, fileInfo.Name())
		if err != nil {
			panic(
				"Could not stat file or folder named " +
					absPath)
		}

		if fileInfo.IsDir() {
			subDirectory := WalkFolder(absPath)
			root.Directories = append(root.Directories, subDirectory)
		} else {
			root.Files = append(root.Files, fileInfo.Name())
		}
	}
	return root
}

/**
From a file system directory, check for .sgrep file and produce
rules based on its contents.
 */
func createDirectoryFromSgrep(directory string) *Directory {
	root := new(Directory)
	root.Name = path.Base(directory)

	potentialSgrepFilename := path.Join(directory, SGREP_FILENAME)
	// true if .sgrep file exists
	if _, err := os.Stat(potentialSgrepFilename); err == nil {
		root.Rules = RuleSliceFromSgrepFile(potentialSgrepFilename)
	} else {
		// no rules to apply
		root.Rules = make([]*Rule, 0)
	}
	return root
}
