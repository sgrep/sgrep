package sgrep

import "path/filepath"

/**
Basic rule struct: when we encounter an sgrep file
*/

type Rule struct {
    // the fully-qualified path for the .sgrep file that this rule was
    // loaded from.
    containingFileAbsPath string
    // the raw text of the associated rule
    rawRuleText string
}


func ConstructRule(containingFileAbsPath, rawRuleText string) *Rule {
    r :=  Rule {}
    r.containingFileAbsPath = containingFileAbsPath
    r.rawRuleText = rawRuleText
    return &r
}


// returns true if this rule filters (ie., says not to look in) file
// named filename.
func (rule* Rule) FileFilterer(filename string) bool {

    didMatch, err :=  filepath.Match(rule.rawRuleText,filename)
    if err != nil {
        panic("Broken match operation in rule")
    }

    if didMatch {
        return true
    }
    return false
}
