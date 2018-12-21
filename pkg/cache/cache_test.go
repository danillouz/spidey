package cache

import (
	"testing"
)

func TestCache(t *testing.T) {
	c := New()
	l := "https://github.com"

	c.Add(l)

	got := c.Check(l)
	want := true

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
