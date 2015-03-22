package sgreplib

import "io/ioutil"
import "path"
import "path/filepath"
import "fmt"
import "os"
import "strings"

const SGREP_FILENAME string = ".sgrep"

type Directory struct {
	// The name of this directory relative to its parent directory.
	// Eg., for structure a/b/c/, would have a directory with name
	// a, name b, and name c.
	name        string
	files       []string
	directories []*Directory
	rules       []*Rule
}

/**
From current directory, look through all subfolders for sgrep files
and connect those to sgrep files in all base directories.

Returns directory representing base of file system.
*/
func GenerateSgrepDirectories(curDirStr string) *Directory {
	curDir := walkFolderForwards(curDirStr)

	// using filepath.Dir includes all folders up until the curDir
	rootDir := walkFolderBackwards(filepath.Dir(curDirStr))

	// join subdirectories to previous directories.
	parentDir := deepestDir(rootDir)
	parentDir.directories = append(parentDir.directories, curDir)
	return rootDir
}

/**
Runs .sgrep filter over all command line argument filenames passed in.
*/
func NonFilteredFilesToGrepOver(absFilenames []string) []string {
	// keys are folder names and values are lists of filenames
	folderNameToFilenameLists := make(map[string][]string)

	// we do not want to construct a directory object for every
	// single file.  Instead, creating a separate directory object
	// for every file system directory.
	for _, absFilename := range absFilenames {
		folderName := filepath.Dir(absFilename)
		folderNameToFilenameLists[folderName] =
			append(folderNameToFilenameLists[folderName], absFilename)
	}

	// actually generate directories and try filtering
	// list of filenames that we actually should grep over
	var toReturn []string

	for folderName := range folderNameToFilenameLists {
		folderDir := GenerateSgrepDirectories(folderName)
		filesInDir := folderNameToFilenameLists[folderName]

		for _, absFilename := range filesInDir {
			if !folderDir.RecursiveShouldFilterFile(absFilename) {
				toReturn = append(toReturn, absFilename)
			}
		}
	}
	return toReturn
}

/**
@param absFilename --- The fully-qualified filename for a file.
Should be a filename reachable from dir.

@returns --- true if one of the .sgrep files in directories screens
the file from being grepped over.
*/
func (dir *Directory) RecursiveShouldFilterFile(absFilename string) bool {
	dirToCheck := dir
	splitList := strings.Split(absFilename, string(filepath.Separator))
	// FIXME: non-portable, putting / at beginning of separated list
	splitList = splitList[1 : len(splitList)-1]

	for index, individualDir := range splitList {
		constructedFile := path.Join(splitList[index:]...)
		constructedFile =
			path.Join(string(filepath.Separator), constructedFile)

		if dirToCheck.shouldFilterFile(constructedFile) {
			return true
		}

		// find next directory to put in
		foundDir := false
		for _, subdir := range dirToCheck.directories {
			if subdir.name == individualDir {
				dirToCheck = subdir
				foundDir = true
				break
			}
		}
		if !foundDir {
			panic("Could not find associated directory")
		}
	}
	return false
}

func (dir *Directory) PrettyPrint() {
	dir.prettyPrintHelper(0)
}

/**
@returns true if dir's rule slice requires us to filter file. false otherwise.
*/
func (dir *Directory) shouldFilterFile(filename string) bool {
	for _, rule := range dir.rules {
		if rule.fileFilterer(filename) {
			return true
		}
	}
	return false
}

func (dir *Directory) ListNonRuleFilteredFiles() []string {
	toReturn := make([]string, 0)

	// check files directly in dir (not including subdirectories)
	for _, filename := range dir.files {
		if !dir.shouldFilterFile(filename) {
			absFilename := path.Join(dir.name, filename)
			toReturn = append(toReturn, absFilename)
		}
	}

	// check subdirectories
	for _, subdir := range dir.directories {
		subdirFileSlice := subdir.ListNonRuleFilteredFiles()

		for _, filename := range subdirFileSlice {
			absFilename := path.Join(dir.name, filename)
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
func (dir *Directory) listFiles() []string {
	toReturn := make([]string, 0)

	for _, filename := range dir.files {
		absFilename := path.Join(dir.name, filename)
		toReturn = append(toReturn, absFilename)
	}

	for _, subdir := range dir.directories {
		subdirFileSlice := subdir.listFiles()

		for _, filename := range subdirFileSlice {
			absFilename := path.Join(dir.name, filename)
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
	fmt.Println(indentStr + dir.name + "/")

	for _, filename := range dir.files {
		fmt.Println(indentStr + "\t" + filename)
	}

	for _, subDirectory := range dir.directories {
		subDirectory.prettyPrintHelper(indentationLevel + 1)
	}
}

/**
  @param dirToWalk The directory that we want to list all
  subdirectories of.

  @returns A Directory struct containing all subfiles and folders.
*/
func walkFolderForwards(dirToWalk string) *Directory {
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
			subDirectory := walkFolderForwards(absPath)
			root.directories = append(root.directories, subDirectory)
		} else {
			root.files = append(root.files, fileInfo.Name())
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
	root.name = path.Base(directory)

	potentialSgrepFilename := path.Join(directory, SGREP_FILENAME)
	// true if .sgrep file exists
	if _, err := os.Stat(potentialSgrepFilename); err == nil {
		root.rules = ruleSliceFromSgrepFile(potentialSgrepFilename)
	} else {
		// no rules to apply
		root.rules = make([]*Rule, 0)
	}
	return root
}

/**
For a file system directory, generate rules from sgrep file (if it
exists), and returns the shallowest directory (root of file system).
Do not descend into files and folders in that directory.
*/
func walkFolderBackwards(dirToWalkStr string) *Directory {
	dir := createDirectoryFromSgrep(dirToWalkStr)
	parentDirStr := filepath.Dir(dirToWalkStr)

	// Means that we got to base of file system and we can go no
	// farther.
	if parentDirStr == dirToWalkStr {
		return dir
	}

	// base directory of file sytem, not necessarily parent directory
	baseDir := walkFolderBackwards(parentDirStr)

	// append dir to end of directory chain.
	parentDir := deepestDir(baseDir)
	parentDir.directories = append(parentDir.directories, dir)

	return baseDir
}

/**
@param dir --- Must have zero or one subdirectories.

@return --- Returns the deepest subdirectory in tree.
*/
func deepestDir(dir *Directory) *Directory {
	for true {
		if len(dir.directories) == 0 {
			break
		}
		dir = dir.directories[0]
	}
	return dir
}
