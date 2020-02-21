package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"
)

const sslRequestOnly = "s3-bucket-ssl-requests-only"

// EncryptionInFlightAWS azure implementation of the encryption in flight for Object Storage feature
type EncryptionInFlightAWS struct {
	ctx         context.Context
	tags        map[string]*string
	httpOption  bool
	httpsOption bool
	session     *session.Session
	evalResults []*configservice.EvaluationResult
}

func (state *EncryptionInFlightAWS) setup() {
	log.Println("Setting up \"EncryptionInFlightAWS\"")
	state.ctx = context.Background()

	// Create Session
	var err error
	state.session, err = session.NewSession()
	if err != nil {
		log.Fatalf("Unable create session to AWS due to %v", err)
	}
}

func (state *EncryptionInFlightAWS) teardown() {
	log.Println("Teardown completed")
}

func (state *EncryptionInFlightAWS) securityControlsThatRestrictDataFromBeingUnencryptedInFlight() error {
	return fmt.Errorf("AWS do not support preventative controls for secure transfer on S3")
}

func (state *EncryptionInFlightAWS) weProvisionAnObjectStorageBucket() error {
	// Not supported
	return nil
}

func (state *EncryptionInFlightAWS) httpAccessIs(arg1 string) error {
	// Not supported
	return nil
}

func (state *EncryptionInFlightAWS) httpsAccessIs(arg1 string) error {
	// Not supported
	return nil
}

func (state *EncryptionInFlightAWS) creationWillWithAnErrorMatching(result, errDescription string) error {
	// Not supported
	return nil
}

func (state *EncryptionInFlightAWS) cSPProvideDetectiveMeasureForNonComplianceSecureTransferOnObjectStorage() error {
	return nil
}

func (state *EncryptionInFlightAWS) weExamineTheDetectiveMeasure() error {
	svc := configservice.New(state.session)
	resp, err := svc.GetComplianceDetailsByConfigRule(&configservice.GetComplianceDetailsByConfigRuleInput{
		ConfigRuleName: aws.String(sslRequestOnly),
	})

	if err != nil { // resp is now filled
		return err
	}
	state.evalResults = resp.EvaluationResults
	return nil
}

func (state *EncryptionInFlightAWS) theDetectiveMeasureIsEnabled() error {
	resultCount := len(state.evalResults)
	if resultCount > 0 {
		log.Printf("AWS Config Rule: \"%v\" evaluation results count: %v", sslRequestOnly, resultCount)
		return nil
	}
	return fmt.Errorf("no evaluation result on AWS Config Rule:\"%v\". [Step Failed]", sslRequestOnly)
}
