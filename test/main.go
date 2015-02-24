package main

import "sgrep/lib"
import "path/filepath"
import "fmt"

func main() {

    // test 1: check that we apply rules correctly
    fmt.Print("py exclude subfolders .... ")
    if !testPyExcludeSubfolders() {
        fmt.Println("FAILED")
    } else {
        fmt.Println("PASSED")
    }

    // test 2: test that we parse individual rules correctly
    fmt.Print("parse rules with comments .... ")
    if !testCommentedRules() {
        fmt.Println("FAILED")
    } else {
        fmt.Println("PASSED")
    }

    // test 3: test that we parse files of rules correctly
    fmt.Print("parse sgrep file .... ")
    if !testSgrepFileRead() {
        fmt.Println("FAILED")
    } else {
        fmt.Println("PASSED")
    }
    
}

/**
Check that correctly excludes python files in subfolders when use *py
globs.
 */
func testPyExcludeSubfolders() bool {
    py_exclude_rule := sgrep.ConstructRule(".sgrep","*py")

    if ! py_exclude_rule.FileFilterer(filepath.Join("a","b","c.py")) {
        return false
    }
    return true
}


