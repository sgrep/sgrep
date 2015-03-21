package sgreplib

import "io/ioutil"
import "testing"

type ruleTextAndExpected struct {
	ruleText          string
	shouldProduceRule bool
	// empty string if should not return a rule from parsing ruleText
	expected string
}

/**
@returns true if parsing textAndExpected.ruleText produces expectedResult
*/
func (textAndExpected *ruleTextAndExpected) correctlyParses() bool {
	rule := parseRule("dummy", textAndExpected.ruleText)
	if rule == nil {
		if !textAndExpected.shouldProduceRule {
			return true
		}
		return false
	}
	return rule.rawRuleText == textAndExpected.expected
}

/**
@param expected --- can be nil if should not return a rule from
parsing ruleText
*/
func constructRuleTextAndExpected(
	ruleText string, shouldProduceRule bool, expected string) *ruleTextAndExpected {

	toReturn := ruleTextAndExpected{}
	toReturn.ruleText = ruleText
	toReturn.shouldProduceRule = shouldProduceRule
	toReturn.expected = expected
	return &toReturn
}

func TestCommentedRules(t *testing.T) {
	var testSlice []*ruleTextAndExpected

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
		if !toTest.correctlyParses() {
			t.Errorf("Commented rule test failed: ", toTest)
		}
	}
}

/**
Creates a temporary sgrep file and tries to read its contents to
ensure we read all the rules we expect in it.
*/
func TestSgrepFileRead(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "empty")

	if err != nil {
		panic("Cannot create temporary file for sgrep read test")
	}

	sgrepFileContents := "a\n"
	sgrepFileContents += "# m\n"
	sgrepFileContents += "fefe   m\n"
	sgrepFileContents += "fefes #  m\n"
	sgrepFileContents += "\n\n"
	sgrepFileContents += "a/b/c*\n"
	sgrepFileContents += "*py\n"
	sgrepFileContents += "a/*/something.txt"

	tmpFile.WriteString(sgrepFileContents)
	tmpFile.Sync()

	// value has no meaning.  only using bool value because can't find
	// golang-native way to use a set.
	expectedRuleContents := map[string]bool{
		"a":                 true,
		"fefe   m":          true,
		"fefes":             true,
		"a/b/c*":            true,
		"*py":               true,
		"a/*/something.txt": true,
	}

	parsedRules := ruleSliceFromSgrepFile(tmpFile.Name())

	if len(parsedRules) != len(expectedRuleContents) {
		t.Fatalf("Parsing didn't return expected number of rules")
	}

	for _, parsedRule := range parsedRules {
		_, exists := expectedRuleContents[parsedRule.rawRuleText]
		if !exists {
			t.Errorf("Unexpected rule returned by parser: ", parsedRule.rawRuleText)
		}
	}
}
