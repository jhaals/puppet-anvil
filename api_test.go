package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/benschw/dns-clb-go/dns"
	"github.com/benschw/opin-go/rando"
	"github.com/benschw/puppet-anvil/api"

	. "gopkg.in/check.v1"
)

func dl(path string, outPath string) {
	out, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	resp, err := http.Get("http://example.com/")
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

var _ = Suite(&TestSuite{})

type TestSuite struct {
	svc    *AnvilService
	client *api.AnvilClient
	path   string
}

func (s *TestSuite) SetUpSuite(c *C) {
	s.path = "/tmp/anvil-modules"
	os.MkdirAll(s.path, 0755)

	addr := dns.Address{Address: "localhost", Port: uint16(rando.Port())}
	s.svc = NewAnvilService(string(addr.Port), s.path)

	go s.svc.Run()
}
func (s *TestSuite) TearDownSuite(c *C) {
	//s.svc.Stop()
}
func (s *TestSuite) SetUpTest(c *C) {
	log.Println("foo")
	dl("https://forgeapi.puppetlabs.com/v3/files/puppetlabs-apache-1.5.0.tar.gz", s.path+"/puppetlabs-apache-1.5.0.tar.gz")
	dl("https://forgeapi.puppetlabs.com/v3/files/puppetlabs-concat-1.2.3.tar.gz", s.path+"/puppetlabs-concat-1.2.3.tar.gz")
	dl("https://forgeapi.puppetlabs.com/v3/files/puppetlabs-stdlib-4.6.0.tar.gz", s.path+"/puppetlabs-stdlib-4.6.0.tar.gz")
}
func (s *TestSuite) TearDownTest(c *C) {
	//os.RemoveAll(s.path)
}

func (s *TestSuite) TestAdd(c *C) {
	// given

	// when
	loc, err := s.client.AddModule(s.path + "/puppetlabs-apache-1.5.0.tar.gz")

	// then
	c.Assert(err, Equals, nil)
	c.Assert(loc, Not(Equals), "ss")
}
