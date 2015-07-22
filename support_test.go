package main

import (
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

func Port() int {
	l, _ := net.Listen("tcp", ":0")
	defer l.Close()
	addrParts := strings.Split(l.Addr().String(), ":")
	port, _ := strconv.Atoi(addrParts[len(addrParts)-1])
	return port
}
func dl(path string, outPath string) {
	out, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	resp, err := http.Get(path)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}
