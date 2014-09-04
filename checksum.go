package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
)

// Checksum file
func Checksum(file string) string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", md5.Sum(data))
}
