package main

import "sgrep/lib"
import "path/filepath"
import "fmt"

func main() {

    // test 1
    fmt.Print("py exclude subfolders .... ")
    if !testPyExcludeSubfolders() {
        fmt.Println("FAILED")
    } else {
        fmt.Println("PASSED")
    }

    // test 2
    fmt.Print("parse rules with comments .... ")
    if !testCommentedRules() {
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


