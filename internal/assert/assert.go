package assert

import (
	"testing"
)

func Equal[T comparable](t *testing.T, got T, want T) {
	t.Helper()

	if got != want {
		t.Errorf("got '%v'; want '%v'", got, want)
	}
}

func EqualSlicesUnordered[T comparable](t *testing.T, got []T, want []T) {
	t.Helper()

	if len(got) != len(want) {
		t.Errorf("length mismatch; got %v; want %v", got, want)
		return
	}

	counts := map[T]int{}

	for _, el := range want {
		counts[el]++
	}

	for _, el := range got {
		if counts[el] < 1 {
			t.Errorf("unexpected element %v; got %v; want %v", el, got, want)
		}
		counts[el]--
	}
}
