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
type EncryptionInFlight interface {
	setup()
	securityControlsThatRestrictDataFromBeingUnencryptedInFlight() error
	weProvisionAnObjectStorageBucket() error
	httpAccessIs(arg1 string) error
	httpsAccessIs(arg1 string) error
	creationWillWithAnErrorMatching(result, errDescription string) error
	cSPProvideDetectiveMeasureForNonComplianceSecureTransferOnObjectStorage() error
	weExamineTheDetectiveMeasure() error
	theDetectiveMeasureIsEnabled() error
	teardown()
}

var opt = godog.Options{Output: colors.Colored(os.Stdout)}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opt.Paths = flag.Args()

	status := godog.RunWithOptions("encryption_in_flight", func(s *godog.Suite) {
		FeatureContext(s)
	}, opt)

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func FeatureContext(s *godog.Suite) {
	var state EncryptionInFlight

	cspEnv := os.Getenv(csp)
	if strings.EqualFold(cspEnv, "azure") {
		state = &EncryptionInFlightAzure{}
	} else if strings.EqualFold(cspEnv, "aws") {
		state = &EncryptionInFlightAWS{}
	} else {
		log.Panicf("Environment variable %s is defined as \"%s\"", csp, cspEnv)
	}

	s.BeforeSuite(state.setup)

	s.Step(`^security controls that restrict data from being unencrypted in flight$`, state.securityControlsThatRestrictDataFromBeingUnencryptedInFlight)
	s.Step(`^we provision an Object Storage bucket$`, state.weProvisionAnObjectStorageBucket)
	s.Step(`^http access is "([^"]*)"$`, state.httpAccessIs)
	s.Step(`^https access is "([^"]*)"$`, state.httpsAccessIs)
	s.Step(`^creation will "([^"]*)" with an error matching "([^"]*)"$`, state.creationWillWithAnErrorMatching)

	s.Step(`^the CSP provides a detective capability for unencrypted data transfer to Object Storage$`, state.cSPProvideDetectiveMeasureForNonComplianceSecureTransferOnObjectStorage)
	s.Step(`^we examine the detective measure$`, state.weExamineTheDetectiveMeasure)
	s.Step(`^the detective measure is enabled$`, state.theDetectiveMeasureIsEnabled)

	s.AfterSuite(state.teardown)
}
