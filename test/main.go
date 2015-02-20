package main

import "sgrep/lib"
import "path/filepath"
import "fmt"

func main() {

    fmt.Print("py exclude .... ")
    if !test_py_exclude() {
        fmt.Println("FAILED")
    } else {
        fmt.Println("PASSED")
    }
}


func test_py_exclude() bool {
    py_exclude_rule := sgrep.ConstructRule(".sgrep","*py")

    if ! py_exclude_rule.FileFilterer(filepath.Join("a","b","c.py")) {
        return false
    }
    return true
}
