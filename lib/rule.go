package sgrep

import "path/filepath"

/**
Basic rule struct: when we encounter an sgrep file
*/

type Rule struct {
    // the fully-qualified path for the .sgrep file that this rule was loaded from.
    containing_file_abs_path string
    // the raw text of the associated rule
    raw_rule_text string
}


func ConstructRule(containing_file_abs_path, raw_rule_text string) *Rule {
    r1 :=  Rule {}
    r1.containing_file_abs_path = containing_file_abs_path
    r1.raw_rule_text = raw_rule_text
    return &r1
}


// returns true if this rule filters (ie., says not to look in) file named filename.
func (rule* Rule) FileFilterer(filename string) bool {

    did_match, err :=  filepath.Match(rule.raw_rule_text,filename)
    if err != nil {
        panic("Broken match operation in rule")
    }

    if did_match {
        return true
    }
    return false
}
