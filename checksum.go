package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
)

// Checksum file
func Checksum(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", md5.Sum(data)), nil
}
