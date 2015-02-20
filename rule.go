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


// returns true if this rule filters (ie., says not to look in) file named filename.
func (rule* Rule) file_filterer(filename string) bool {

    did_match, err :=  filepath.Match(rule.raw_rule_text,filename)
    if err != nil {
        panic("Broken match operation in rule")
    }

    if did_match {
        return true
    }
    return false
}
