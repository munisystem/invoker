package credential

import "testing"

func TestGeneratePassword(t *testing.T) {
	expected := 12
	actual := GeneratePassword(expected)

	if expected != len(actual) {
		t.Fatalf("didn't match str length: expected %d, actual %d", expected, len(actual))
	}
}
