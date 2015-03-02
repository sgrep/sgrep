package sgrep

import "path/filepath"
import "testing"

/**
 * Check that correctly excludes python files in subfolders when use *py globs.
 */
func TestPyExcludeSubfolders(t *testing.T) {
	py_exclude_rule := ConstructRule(".sgrep", "*py")

	if !py_exclude_rule.FileFilterer(filepath.Join("a", "b", "c.py")) {
		t.Error("Python files not correctly excluded by *py glob")
	}
}
