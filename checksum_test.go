package main

import (
	"testing"
)

func TestChecksum(t *testing.T) {
	_, err := Checksum("main.go")
	if err != nil {
		t.Fatal(err)
	}
}
