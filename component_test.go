package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/jhaals/puppet-anvil/api"
	"github.com/jhaals/puppet-anvil/service"

	. "gopkg.in/check.v1"
)

var _ = Suite(&ComponentTestSuite{})

type ComponentTestSuite struct {
	svc    *service.AnvilService
	client *api.AnvilClient
	path   string
	port   int
}

// create and start a service, create a client
// download fixture modules if they aren't present
func (s *ComponentTestSuite) SetUpSuite(c *C) {
	s.path = "./tmp-test"

	os.MkdirAll(s.path+"/modules", 0755)
	if _, err := os.Stat(s.path + "/puppetlabs-apache-1.5.0.tar.gz"); err != nil {
		dl("https://forgeapi.puppetlabs.com/v3/files/puppetlabs-apache-1.5.0.tar.gz", s.path+"/puppetlabs-apache-1.5.0.tar.gz")
		dl("https://forgeapi.puppetlabs.com/v3/files/puppetlabs-apache-1.10.0.tar.gz", s.path+"/puppetlabs-apache-1.10.0.tar.gz")
		dl("https://forgeapi.puppetlabs.com/v3/files/puppetlabs-concat-1.2.3.tar.gz", s.path+"/puppetlabs-concat-1.2.3.tar.gz")
		dl("https://forgeapi.puppetlabs.com/v3/files/puppetlabs-stdlib-4.6.0.tar.gz", s.path+"/puppetlabs-stdlib-4.6.0.tar.gz")
	}

	s.port = Port()
	s.svc = service.New(strconv.Itoa(s.port), s.path+"/modules")
	s.client = &api.AnvilClient{Address: fmt.Sprintf("localhost:%d", s.port)}

	go s.svc.Run()
}

// tear down installed files
// (leave fixture downloads for future tests)
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
func (s *ComponentTestSuite) TestGetReleaseByModule(c *C) {
	// given
	file := s.path + "/puppetlabs-apache-1.5.0.tar.gz"
	f, _ := os.Open(file)
	defer f.Close()
	loc, _ := s.client.PublishModule(f, "puppetlabs-apache-1.5.0.tar.gz")

	// when
	resp, err := s.client.GetReleaseByModule("puppetlabs", "apache")

	// then
	c.Assert(err, Equals, nil)
	c.Assert(resp.Results[0].FileUri, Equals, loc)
	c.Assert(resp.Results[0].Version, Equals, "1.5.0")
}

// client should return empty array when module not found
func (s *ComponentTestSuite) TestGetReleaseByModuleNotFound(c *C) {
	// when
	resp, err := s.client.GetReleaseByModule("puppetlabs", "apache")

	// then
	c.Assert(err, Equals, nil)
	c.Assert(len(resp.Results), Equals, 0)
}

// client should return error when searching with invalid module name
func (s *ComponentTestSuite) TestGetReleaseByModuleInvalidModule(c *C) {
	// when
	_, err := s.client.GetReleaseByModule("puppetlabs", "apache-foo")

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

// client should return releases for a module
func (s *ComponentTestSuite) TestGetModulesByUserModule(c *C) {
	// given
	file := s.path + "/puppetlabs-apache-1.5.0.tar.gz"
	f, _ := os.Open(file)
	defer f.Close()
	loc, _ := s.client.PublishModule(f, "puppetlabs-apache-1.5.0.tar.gz")

	// when
	resp, err := s.client.GetModulesByUserModule("puppetlabs", "apache")

	// then
	c.Assert(err, Equals, nil)
	c.Assert(resp.Releases[0].FileUri, Equals, loc)
	c.Assert(resp.Releases[0].Version, Equals, "1.5.0")
}

// client should return data for a particular release
func (s *ComponentTestSuite) TestGetReleaseByUserModuleVersion(c *C) {
	// given
	file := s.path + "/puppetlabs-apache-1.5.0.tar.gz"
	f, _ := os.Open(file)
	defer f.Close()
	loc, _ := s.client.PublishModule(f, "puppetlabs-apache-1.5.0.tar.gz")

	// when
	resp, err := s.client.GetReleaseByUserModuleVersion("puppetlabs", "apache", "1.5.0")

	// then
	c.Assert(err, Equals, nil)
	c.Assert(resp.FileUri, Equals, loc)
	c.Assert(resp.Version, Equals, "1.5.0")
	c.Assert(resp.Metadata.Name, Equals, "puppetlabs-apache")
}

// client should return most recent release
func (s *ComponentTestSuite) TestGetModulesByUserModuleLatest(c *C) {
	// given
	file1 := s.path + "/puppetlabs-apache-1.5.0.tar.gz"
	f1, _ := os.Open(file1)
	defer f1.Close()
	loc1, _ := s.client.PublishModule(f1, "puppetlabs-apache-1.5.0.tar.gz")

	file2 := s.path + "/puppetlabs-apache-1.10.0.tar.gz"
	f2, _ := os.Open(file2)
	defer f2.Close()
	loc2, _ := s.client.PublishModule(f2, "puppetlabs-apache-1.10.0.tar.gz")

	// when
	resp, err := s.client.GetModulesByUserModule("puppetlabs", "apache")

	// then
	c.Assert(err, Equals, nil)
	c.Assert(resp.Releases[0].FileUri, Equals, loc2)
	c.Assert(resp.Releases[0].Version, Equals, "1.10.0")
	c.Assert(resp.Releases[1].FileUri, Equals, loc1)
	c.Assert(resp.Releases[1].Version, Equals, "1.5.0")
}
