package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	citihubAws "citihub.com/compliance-as-code/internal/aws"
	"citihub.com/compliance-as-code/internal/azureutil"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	sslRequestOnly     = "s3-bucket-ssl-requests-only"
	awsSecureTransport = "aws:SecureTransport"
	maxRetry           = 10
	sleepTime          = 60 * time.Second
)

// EncryptionInFlightAWS stores the context used for the Encryption in Flight test on AWS.
type EncryptionInFlightAWS struct {
	ctx         context.Context
	tags        map[string]*string
	httpOption  bool
	httpsOption bool
	session     *session.Session
	s3Svc       *s3.S3
	configSvc   *configservice.ConfigService
	bucketName  string
	runningErr  error
	region      string
}

func (state *EncryptionInFlightAWS) setup() {
	log.Println("[DEBUG] Setting up \"EncryptionInFlightAWS\"")
	state.ctx = context.Background()
	state.region = os.Getenv("AWS_REGION")
	// Create Session
	var err error
	state.session, err = session.NewSession()
	state.s3Svc = s3.New(state.session)
	state.configSvc = configservice.New(state.session)
	if err != nil {
		log.Fatalf("unable create session to AWS due to %v", err)
	}
}

func (state *EncryptionInFlightAWS) teardown() {
	_, err := state.s3Svc.DeleteBucket(&s3.DeleteBucketInput{Bucket: aws.String(state.bucketName)})
	if err != nil {
		log.Printf("[ERROR] error in deleting test bucket %v. Please manually clean up", state.bucketName)
	} else {
		log.Printf("[DEBUG] Bucket %v clean up successful.", state.bucketName)
	}
	log.Println("[DEBUG] Teardown completed")
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

func (state *EncryptionInFlightAWS) detectObjectStorageUnencryptedTransferAvailable() error {
	return nil
}

func (state *EncryptionInFlightAWS) detectObjectStorageUnencryptedTransferEnabled() error {
	_, err := state.configSvc.GetComplianceDetailsByConfigRule(&configservice.GetComplianceDetailsByConfigRuleInput{
		ConfigRuleName: aws.String(sslRequestOnly),
	})
	return err
}

func (state *EncryptionInFlightAWS) createUnencryptedTransferObjectStorage() error {
	state.bucketName = fmt.Sprintf("test%sunencbucket", azureutil.RandString(5))
	resp, err := state.s3Svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(state.bucketName),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(state.region),
		},
	})

	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Created Bucket: %v", resp)
	return nil
}

// Wait for Config rule to detect the bucket has been created
func (state *EncryptionInFlightAWS) detectsTheObjectStorage() error {
	log.Printf("[DEBUG] Waiting for bucket to be detected by Config Rule...")
	for i := 0; i < maxRetry; i++ {
		resp, err := state.configSvc.GetComplianceDetailsByConfigRule(&configservice.GetComplianceDetailsByConfigRuleInput{
			ConfigRuleName: aws.String(sslRequestOnly),
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
				ConfigRuleName: aws.String(sslRequestOnly),
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
			log.Printf("[DEBUG] AWS Config Rule: \"%v\" evaluation results count: %v", sslRequestOnly, resultCount)
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
	return fmt.Errorf("failed to find bucket '%v' in evaluation result of AWS Config Rule:'%v' [Step Failed]", state.bucketName, sslRequestOnly)
}

// Checking with a sleep and retry mechanism on the bucket being remediated to secure transport enabled
func (state *EncryptionInFlightAWS) encryptedDataTrafficIsEnforced() error {
	for i := 0; i < maxRetry; i++ {
		log.Printf("[DEBUG] Checking bucket policy for secure transport setting...")
		err := state.checkIsSSLRequestOnly()
		if err == nil { // Deny unsecured transport
			return err
		}
		log.Printf("[DEBUG] Bucket policy still insecure wait for %d s, retry %d/%d", sleepTime/time.Second, i, maxRetry)
		time.Sleep(sleepTime)
	}
	return fmt.Errorf("after 5 mins the bucket '%v' is still not remediated [Step Failed]", state.bucketName)
}

// This is just to check if it there is a bucket policy that's configured with SSL
// return nil when found the right bucket policy statement on secure transport
func (state *EncryptionInFlightAWS) checkIsSSLRequestOnly() error {
	result, err := state.s3Svc.GetBucketPolicy(&s3.GetBucketPolicyInput{
		Bucket: aws.String(state.bucketName),
	})
	if err != nil {
		return err
	}

	var policyDoc citihubAws.PolicyDocument
	err = json.Unmarshal([]byte(*result.Policy), &policyDoc)
	if err != nil {
		return err
	}
	for _, stmt := range *policyDoc.Statement {
		if *stmt.Effect == "Deny" {
			conditionMap := *stmt.Condition

			// Only start checking of a "Bool" condition
			if conditionMap["Bool"] != nil {
				conditionKey := conditionMap["Bool"].(map[string]interface{})
				v := conditionKey[awsSecureTransport]
				if v != nil {
					log.Printf("[DEBUG] %v: %v", awsSecureTransport, v)

					// Only return nil positive when found the right bucket policy statement
					if v == "false" {
						return nil
					}
				}
			}
		}
	}
	return fmt.Errorf("incorrect bucket policy setting on '%v': %v", state.bucketName, result)
}
