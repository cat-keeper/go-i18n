package i18n

import "testing"

func TestReplaceArgs(t *testing.T) {
	template := "Hello, {0}, it's {1}"
	data := []interface{}{"John", "a beautiful day"}
	expected := "Hello, John, it's a beautiful day"
	actual, err := ReplaceArgs(template, data...)
	if actual != expected || err != nil {
		t.Errorf("ReplaceArgs(%q) returned %q; expected %q", template, actual, expected)
	}
}
