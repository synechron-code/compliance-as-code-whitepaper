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

// EncryptionAtRest is an interface. For each CSP specific implementation
type EncryptionAtRest interface {
	setup()
	policyOrRuleAvailable() error
	checkPolicyOrRuleAssignment() error
	policyOrRuleAssigned() error
	prepareToCreateContainer() error
	createContainerWithEncryptionOption(encryptionOption string) error
	createResult(result string) error
	teardown()
}

var opt = godog.Options{Output: colors.Colored(os.Stdout)}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opt.Paths = flag.Args()

	status := godog.RunWithOptions("encryption_at_rest", func(s *godog.Suite) {
		FeatureContext(s)
	}, opt)

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func FeatureContext(s *godog.Suite) {
	var state EncryptionAtRest

	cspEnv := os.Getenv(csp)
	if strings.EqualFold(cspEnv, "azure") {
		state = &EncryptionAtRestAzure{}
	} else if strings.EqualFold(cspEnv, "aws") {
		state = &EncryptionAtRestAWS{}
	} else {
		log.Panicf("Environment variable %s is defined as \"%s\"", csp, cspEnv)
	}

	s.BeforeSuite(state.setup)

	s.Step(`^the CSP provides a detective capability for unencrypted Object Storage containers$`, state.policyOrRuleAvailable)
	s.Step(`^we examine the detective measure$`, state.checkPolicyOrRuleAssignment)
	s.Step(`^the detective measure is enabled$`, state.policyOrRuleAssigned)
	s.Step(`^security controls that enforce data at rest encryption for Object Storage are applied$`, state.policyOrRuleAvailable)
	s.Step(`^we provision an Object Storage container$`, state.prepareToCreateContainer)
	s.Step(`^it is created with encryption option "([^"]*)"$`, state.createContainerWithEncryptionOption)
	s.Step(`^creation will "([^"]*)"$`, state.createResult)

	s.AfterSuite(state.teardown)
}
