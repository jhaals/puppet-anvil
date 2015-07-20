package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/benschw/opin-go/rando"
	"github.com/benschw/puppet-anvil/api"
	"github.com/benschw/puppet-anvil/service"

	. "gopkg.in/check.v1"
)

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

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&ComponentTestSuite{})

type ComponentTestSuite struct {
	svc    *service.AnvilService
	client *api.AnvilClient
	path   string
	port   int
}

func (s *ComponentTestSuite) SetUpSuite(c *C) {
	s.path = "./tmp-test"

	os.MkdirAll(s.path+"/modules", 0755)
	if _, err := os.Stat(s.path + "/puppetlabs-apache-1.5.0.tar.gz"); err != nil {
		dl("https://forgeapi.puppetlabs.com/v3/files/puppetlabs-apache-1.5.0.tar.gz", s.path+"/puppetlabs-apache-1.5.0.tar.gz")
		dl("https://forgeapi.puppetlabs.com/v3/files/puppetlabs-concat-1.2.3.tar.gz", s.path+"/puppetlabs-concat-1.2.3.tar.gz")
		dl("https://forgeapi.puppetlabs.com/v3/files/puppetlabs-stdlib-4.6.0.tar.gz", s.path+"/puppetlabs-stdlib-4.6.0.tar.gz")
	}

	s.port = rando.Port()
	s.svc = service.New(strconv.Itoa(s.port), s.path+"/modules")
	s.client = &api.AnvilClient{Address: fmt.Sprintf("localhost:%d", s.port)}

	go s.svc.Run()
}
func (s *ComponentTestSuite) TearDownSuite(c *C) {
	s.svc.Stop()
}
func (s *ComponentTestSuite) TearDownTest(c *C) {
	os.RemoveAll(s.path + "/modules")
}

// client should publish module to filesystem
func (s *ComponentTestSuite) TestPublish(c *C) {
	// given
	file := s.path + "/puppetlabs-apache-1.5.0.tar.gz"
	f, _ := os.Open(file)
	defer f.Close()

	// when
	loc, err := s.client.PublishModule(f, "puppetlabs-apache-1.5.0.tar.gz")

	// then
	c.Assert(err, Equals, nil)
	c.Assert(loc, Equals, "/v3/files/puppetlabs/apache/puppetlabs-apache-1.5.0.tar.gz")
}

// client should error out with invlid module fileName
func (s *ComponentTestSuite) TestPublishError(c *C) {
	// given
	file := s.path + "/puppetlabs-apache-1.5.0.tar.gz"
	f, _ := os.Open(file)
	defer f.Close()

	// when
	_, err := s.client.PublishModule(f, "puppetlabs-apache.tar.gz")

	// then
	c.Assert(err, Not(Equals), nil)
}

// client should return list releases for a given module
func (s *ComponentTestSuite) TestGetRelease(c *C) {
	// given
	file := s.path + "/puppetlabs-apache-1.5.0.tar.gz"
	f, _ := os.Open(file)
	defer f.Close()
	loc, _ := s.client.PublishModule(f, "puppetlabs-apache-1.5.0.tar.gz")

	// when
	resp, err := s.client.GetRelease("puppetlabs", "apache")

	// then
	c.Assert(err, Equals, nil)
	c.Assert(resp.Results[0].FileUri, Equals, loc)
	c.Assert(resp.Results[0].Version, Equals, "1.5.0")
}

// client should return empty array when module not found
func (s *ComponentTestSuite) TestGetReleaseNotFound(c *C) {
	// when
	resp, err := s.client.GetRelease("puppetlabs", "apache")

	// then
	c.Assert(err, Equals, nil)
	c.Assert(len(resp.Results), Equals, 0)
}

// client should return error when searching with invalid module name
func (s *ComponentTestSuite) TestGetReleaseInvalidModule(c *C) {
	// when
	_, err := s.client.GetRelease("puppetlabs", "apache-foo")

	// then
	c.Assert(err, Not(Equals), nil)
}

// file should be accessible from location-header address
func (s *ComponentTestSuite) TestDownloadFile(c *C) {
	// given
	file := s.path + "/puppetlabs-apache-1.5.0.tar.gz"
	expected, _ := service.Checksum(file)
	f, _ := os.Open(file)
	defer f.Close()
	loc, _ := s.client.PublishModule(f, "puppetlabs-apache-1.5.0.tar.gz")

	outPath := s.path + "/out.tar.gz"
	// when
	output, _ := os.Create(outPath)
	defer output.Close()

	resp, _ := http.Get(fmt.Sprintf("http://localhost:%d/%s", s.port, loc))
	defer resp.Body.Close()

	io.Copy(output, resp.Body)
	// then

	found, err := service.Checksum(outPath)

	c.Assert(err, Equals, nil)
	c.Assert(found, Equals, expected)
}
