package sgrep

import "path/filepath"
import "strings"
import "os"
import "bufio"

const COMMENT_STRING string = "#"

/**
Basic rule struct: when we encounter an sgrep file
*/

type Rule struct {
	// the fully-qualified path for the .sgrep file that this rule was
	// loaded from.
	containingFileAbsPath string
	// the raw text of the associated rule
	RawRuleText string
}

func ConstructRule(containingFileAbsPath, rawRuleText string) *Rule {
	r := Rule{}
	r.containingFileAbsPath = containingFileAbsPath
	r.RawRuleText = rawRuleText
	return &r
}

// returns true if this rule filters (ie., says not to look in) file
// named filename.
func (rule *Rule) FileFilterer(filename string) bool {

	didMatch, err := filepath.Match(rule.RawRuleText, filename)
	if err != nil {
		panic("Broken match operation in rule")
	}

	if didMatch {
		return true
	}
	return false
}

/**
@param sgrepAbsFilename --- The absolute path for the file that
contains the line that we're about to parse into a rule.

@param text --- The line that we're trying to turn into an sgrep rule

@returns nil if text does not produce a valid sgrep rule (eg., empty
line, commented-out rule, etc.)  Otherwise, returns rule.
*/
func ParseRule(sgrepAbsFilename, text string) *Rule {
	// ignore comments
	commentIndex := strings.Index(text, COMMENT_STRING)
	if commentIndex != -1 {
		text = text[0:commentIndex]
	}
	text = strings.TrimSpace(text)
	// ignore blank lines
	if text == "" {
		return nil
	}
	return ConstructRule(sgrepAbsFilename, text)
}

/**
@param absFilename --- Absolute filename to read sgrep rules from.
*/
func RuleSliceFromSgrepFile(absFilename string) []*Rule {
	var toReturn []*Rule
	fh, err := os.Open(absFilename)
	if err != nil {
		panic("Could not open " + absFilename +
			" for reading sgrep rules.")
	}
	// at end of function, close fh
	defer fh.Close()

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		newRule := ParseRule(absFilename, scanner.Text())
		if newRule != nil {
			toReturn = append(toReturn, newRule)
		}

	}
	err = scanner.Err()
	if err != nil {
		panic("Error scanning file " + absFilename +
			" while reading sgrep rules.")
	}

	return toReturn
}
