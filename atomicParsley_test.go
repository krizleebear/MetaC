package main

import "testing"

func Test_removeDuplicateWhitespace(t *testing.T) {
	assertEquals(t, "a b", removeDuplicateWhitespace("a   b"))
}
