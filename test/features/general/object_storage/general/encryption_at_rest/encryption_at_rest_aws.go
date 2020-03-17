package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"os"

	"citihub.com/compliance-as-code/internal/azureutil"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	encryptionAtRestRule = "s3-bucket-server-side-encryption-enabled"
	sleepTime            = 30 * time.Second
	maxRetry             = 10
)

// EncryptionAtRestAWS azure implementation of the encryption in flight for Object Storage feature
type EncryptionAtRestAWS struct {
	ctx              context.Context
	session          *session.Session
	evalResults      []*configservice.EvaluationResult
	s3Svc            *s3.S3
	configSvc        *configservice.ConfigService
	bucketName       string
	runningErr       error
	setEncryptionErr error
    region           string
}

func (state *EncryptionAtRestAWS) setup() {
	log.Println("[DEBUG] Setting up \"EncryptionAtRestAWS\"")
	state.ctx = context.Background()
    state.region = os.Getenv("AWS_REGION")

	// Create Session
	var err error
	state.session, err = session.NewSession()
	state.s3Svc = s3.New(state.session)
	state.configSvc = configservice.New(state.session)
	if err != nil {
		log.Fatalf("Unable create session to AWS due to %v", err)
	}
}

func (state *EncryptionAtRestAWS) teardown() {
	state.deleteCurrentTestBucket()
	log.Println("[DEBUG] Teardown completed")
}

func (state *EncryptionAtRestAWS) securityControlsThatRestrictDataFromBeingUnencryptedAtRest() error {
	return fmt.Errorf("AWS do not have preventive measure but instead reliant on detective measure")
}

func (state *EncryptionAtRestAWS) weProvisionAnObjectStorageBucket() error {
	return nil
}

func (state *EncryptionAtRestAWS) encryptionAtRestIs(encryptionOption string) error {
	return nil
}

func (state *EncryptionAtRestAWS) creationWillWithAnErrorMatching(result string) error {
	return nil
}

func (state *EncryptionAtRestAWS) policyOrRuleAvailable() error {
	// It is available
	log.Printf("[DEBUG] Checking AWS Config Rule: %s", encryptionAtRestRule)
	return nil
}

func (state *EncryptionAtRestAWS) checkPolicyOrRuleAssignment() error {
	svc := configservice.New(state.session)
	resp, err := svc.GetComplianceDetailsByConfigRule(&configservice.GetComplianceDetailsByConfigRuleInput{
		ConfigRuleName: aws.String(encryptionAtRestRule),
	})

	if err != nil { // resp is now filled
		return err
	}
	state.evalResults = resp.EvaluationResults
	return nil
}

func (state *EncryptionAtRestAWS) policyOrRuleAssigned() error {
	resultCount := len(state.evalResults)
	if resultCount > 0 {
		log.Printf("[DEBUG] AWS Config Rule: \"%v\" evaluation results count: %v", encryptionAtRestRule, resultCount)
		return nil
	}
	return fmt.Errorf("no evaluation result on AWS Config Rule:\"%v\". [Step Failed]", encryptionAtRestRule)
}

func (state *EncryptionAtRestAWS) prepareToCreateContainer() error {

	return nil
}

func (state *EncryptionAtRestAWS) createContainerWithoutEncryption() error {
	state.bucketName = fmt.Sprintf("test%sunencbucket", azureutil.RandString(5))
	resp, err := state.s3Svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(state.bucketName),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(state.region),
		},
	})
	log.Printf("[DEBUG] Created Bucket: %v", *resp.Location)
	return err
}

// Wait for Config rule to detect the bucket has been created
func (state *EncryptionAtRestAWS) detectiveDetectsNonCompliant() error {
	log.Printf("[DEBUG] Waiting for bucket to be detected by Config Rule...")
	for i := 0; i < maxRetry; i++ {
		resp, err := state.configSvc.GetComplianceDetailsByConfigRule(&configservice.GetComplianceDetailsByConfigRuleInput{
			ConfigRuleName: aws.String(encryptionAtRestRule),
			Limit:          aws.Int64(100),
		})
		if err != nil {
			return err
		}
		a := resp.EvaluationResults
		next := resp.NextToken
		log.Printf("[DEBUG] nextToken: %v", resp.NextToken)

		// This is to get all the compliance details results if there are over 100 and span over multiple pages.
		for {
			if next == nil {
				break
			}
			resp, err := state.configSvc.GetComplianceDetailsByConfigRule(&configservice.GetComplianceDetailsByConfigRuleInput{
				ConfigRuleName: aws.String(encryptionAtRestRule),
				Limit:          aws.Int64(100),
				NextToken:      next,
			})
			if err != nil {
				return err
			}
			a = append(a, resp.EvaluationResults...)
			next = resp.NextToken
		}

		resultCount := len(resp.EvaluationResults)
		if resultCount > 0 {
			log.Printf("[DEBUG] AWS Config Rule: \"%v\" evaluation results count: %v", encryptionAtRestRule, resultCount)
			for _, e := range resp.EvaluationResults {
				id := e.EvaluationResultIdentifier.EvaluationResultQualifier.ResourceId

				// Only interested in the bucket we created
				if *id == state.bucketName {
					log.Printf("[DEBUG] Bucket '%v' is '%v'", *id, *e.ComplianceType)
					return nil
				}
			}
		}
		log.Printf("[DEBUG] Config Rule not pick up bucket '%v' yet wait for %d s, retry %d/%d", state.bucketName, sleepTime/time.Second, i, maxRetry)
		time.Sleep(sleepTime)
	}
	return fmt.Errorf("failed to find bucket '%v' in evaluation result of AWS Config Rule:'%v' [Step Failed]", state.bucketName, encryptionAtRestRule)
}

func (state *EncryptionAtRestAWS) containerIsRemediated() error {
	for i := 0; i < maxRetry; i++ {
		log.Printf("[DEBUG] Checking bucket policy for SSE setting")
		encrypted := state.checkBucketEncryption()
		if encrypted { // Remediated
			return nil
		}
		log.Printf("[DEBUG] Bucket policy still unencrypted wait for %d s, retry %d/%d", sleepTime/time.Second, i, maxRetry)
		time.Sleep(sleepTime)
	}
	return fmt.Errorf("after 5 mins the bucket '%v' is still not remediated [Step Failed]", state.bucketName)
}

func (state *EncryptionAtRestAWS) checkBucketEncryption() bool {
	_, err := state.s3Svc.GetBucketEncryption(&s3.GetBucketEncryptionInput{
		Bucket: aws.String(state.bucketName),
	})
	if err != nil {
		return false
	}

	return true
}

func (state *EncryptionAtRestAWS) deleteCurrentTestBucket() {
	_, err := state.s3Svc.DeleteBucket(&s3.DeleteBucketInput{Bucket: aws.String(state.bucketName)})
	if err != nil {
		log.Printf("[ERROR] Error in deleting test bucket %v. Please manually clean up.", state.bucketName)
	} else {
		log.Printf("[DEBUG] Bucket %v clean up successful.", state.bucketName)
	}
}
