package service

import (
	"testing"
)

func TestChecksum(t *testing.T) {
	_, err := Checksum("checksum.go")
	if err != nil {
		t.Fatal(err)
	}
}
