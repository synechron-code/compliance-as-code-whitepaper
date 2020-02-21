package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/colors"
)

const csp = "CSP"

// EncryptionInFlight is an interface. For each CSP specific implementation
type AccessWhitelisting interface {
	setup()
	isCspCapable() error
	examineStorageContainer(containerName string) error
	whitelistingIsConfigured() error
	checkPolicyAssigned() error
	prepareToCreateStorageContainer() error
	createWithWhiteList(ipPrefix string) error
	creationWill(result string) error
	teardown()
}

var opt = godog.Options{Output: colors.Colored(os.Stdout)}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opt.Paths = flag.Args()

	status := godog.RunWithOptions("access_whitelisting_test", func(s *godog.Suite) {
		FeatureContext(s)
	}, opt)

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func FeatureContext(s *godog.Suite) {
	var state AccessWhitelisting

	cspEnv := os.Getenv(csp)
	if strings.EqualFold(cspEnv, "azure") {
		state = &AccessWhitelistingAzure{}
	} else if strings.EqualFold(cspEnv, "aws") {
		state = &AccessWhitelistingAWS{}
	} else {
		log.Panicf("Environment variable %s is defined as \"%s\"", csp, cspEnv)
	}

	s.BeforeSuite(state.setup)

	s.Step(`^the CSP provides a whitelisting capability for Object Storage containers$`, state.isCspCapable)
	s.Step(`^we examine the Object Storage container in environment variable "([^"]*)"$`, state.examineStorageContainer)
	s.Step(`^whitelisting is configured with the given IP address range or an endpoint$`, state.whitelistingIsConfigured)
	s.Step(`^security controls that Prevent Object Storage from being created without network source address whitelisting are applied$`, state.checkPolicyAssigned)
	s.Step(`^we provision an Object Storage container$`, state.prepareToCreateStorageContainer)
	s.Step(`^it is created with whitelisting entry "([^"]*)"$`, state.createWithWhiteList)
	s.Step(`^creation will "([^"]*)"$`, state.creationWill)

	s.AfterSuite(state.teardown)
}
