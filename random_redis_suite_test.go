// Test suite setup for the random-redis package
package randomredis

import (
	// Standard lib
	"io/ioutil"
	"testing"

	// Third-party
	log "github.com/Sirupsen/logrus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Tests the random-redis package
func TestConfig(t *testing.T) {
	// Register gomega fail handler
	RegisterFailHandler(Fail)

	// Have go's testing package run package specs
	RunSpecs(t, "Random Redis Suite")
}

func init() {
	// Set logger output so as not to log during tests
	log.SetOutput(ioutil.Discard)
}
