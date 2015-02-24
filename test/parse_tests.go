package main

import "sgrep/lib"
import "io/ioutil"

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


/**
Creates a temporary sgrep file and tries to read its contents to
ensure we read all the rules we expect in it.
*/
func testSgrepFileRead() bool {
    tmpFile, err := ioutil.TempFile("","empty")

    if err != nil {
        panic("Cannot create temporary file for sgrep read test")
    }

    sgrepFileContents := "a\n"
    sgrepFileContents += "# m\n"
    sgrepFileContents += "fefe   m\n"
    sgrepFileContents += "fefes #  m\n"
    sgrepFileContents += "\n\n"
    
    tmpFile.WriteString(sgrepFileContents)
    tmpFile.Sync()

    // value has no meaning.  only using bool value because can't find
    // golang-native way to use a set.
    expectedRuleContents := map[string]bool {
        "a": true,
        "fefe": true,
        "fefes": true,
    }

    parsedRules := sgrep.RuleSliceFromSgrepFile(tmpFile.Name())

    if len(parsedRules) != len(expectedRuleContents) {
        return false
    }

    for _, parsedRule := range parsedRules {
        
        _, exists := expectedRuleContents[parsedRule.RawRuleText]
        if ! exists {
            return false
        }
    }
    return true
}

