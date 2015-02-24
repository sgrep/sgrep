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



type ruleTextAndExpected struct {
    ruleText string
    shouldProduceRule bool
    // empty string if should not return a rule from parsing ruleText
    expected string
}

/**
@returns true if parsing textAndExpected.ruleText produces expectedResult
*/
func (textAndExpected *ruleTextAndExpected) correctlyParses () bool {
    rule := sgrep.ParseRule("dummy", textAndExpected.ruleText)
    if rule == nil {
        if ! textAndExpected.shouldProduceRule {
            return true
        }
        return false
    }
    return rule.RawRuleText == textAndExpected.expected
}

/**
@param expected --- can be nil if should not return a rule from
parsing ruleText
*/
func constructRuleTextAndExpected (
    ruleText string, shouldProduceRule bool, expected string) *ruleTextAndExpected {
        
        toReturn := ruleTextAndExpected {}
        toReturn.ruleText = ruleText
        toReturn.shouldProduceRule = shouldProduceRule
        toReturn.expected = expected
        return &toReturn
}

func testCommentedRules() bool {
    var testSlice [] *ruleTextAndExpected
    
    // should not produce a rule when call ParseRule on it.
    testSlice = append(
        testSlice,
        constructRuleTextAndExpected("#something", false, ""))
    testSlice = append(
        testSlice,
        constructRuleTextAndExpected("#", false, ""))
    testSlice = append(
        testSlice,
        constructRuleTextAndExpected("   #", false, ""))
    testSlice = append(
        testSlice,
        constructRuleTextAndExpected("# other   #", false, ""))

    // should produce rules when call parse on it
    testSlice = append(
        testSlice,
        constructRuleTextAndExpected("a", true, "a"))
    testSlice = append(
        testSlice,
        constructRuleTextAndExpected("a#", true, "a"))
    testSlice = append(
        testSlice,
        constructRuleTextAndExpected("a  #", true, "a"))

    for _, toTest := range testSlice {

        if ! toTest.correctlyParses() {
            return false
        }
    }
    return true
}
