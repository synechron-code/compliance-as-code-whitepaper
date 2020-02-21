package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"citihub.com/compliance-as-code/internal/azureutil"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/s3"
)

const encryptionAtRestRule = "s3-bucket-server-side-encryption-enabled"

const sleepTime = 30 * time.Second

const maxRetry = 10

// EncryptionAtRestAWS azure implementation of the encryption in flight for Object Storage feature
type EncryptionAtRestAWS struct {
	ctx              context.Context
	session          *session.Session
	evalResults      []*configservice.EvaluationResult
	svc              *s3.S3
	bucketName       string
	runningErr       error
	setEncryptionErr error
}

func (state *EncryptionAtRestAWS) setup() {
	log.Println("Setting up \"EncryptionAtRestAWS\"")
	state.ctx = context.Background()

	// Create Session
	var err error
	state.session, err = session.NewSession()
	state.svc = s3.New(state.session)
	if err != nil {
		log.Fatalf("Unable create session to AWS due to %v", err)
	}
}

func (state *EncryptionAtRestAWS) teardown() {
	log.Println("Teardown completed")
}

func (state *EncryptionAtRestAWS) policyOrRuleAvailable() error {
	// It is available
	log.Printf("Checking AWS Config Rule: %s", encryptionAtRestRule)
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
		log.Printf("AWS Config Rule: \"%v\" evaluation results count: %v", encryptionAtRestRule, resultCount)
		return nil
	}
	return fmt.Errorf("no evaluation result on AWS Config Rule:\"%v\". [Step Failed]", encryptionAtRestRule)
}

func (state *EncryptionAtRestAWS) prepareToCreateContainer() error {
	state.bucketName = "test" + azureutil.RandStringBytesMaskImprSrcUnsafe(5) + "bucket"
	return nil
}

func (state *EncryptionAtRestAWS) createContainerWithEncryptionOption(encryptionOption string) error {
	if encryptionOption == "true" {
		state.bucketName = state.bucketName + "enc"
	} else {
		state.bucketName = state.bucketName + "unenc"
	}

	resp, err := state.svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(state.bucketName),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String("ap-southeast-1"),
		},
	})

	if err == nil {
		log.Printf("Created Bucket: %v", resp)
		if encryptionOption == "true" {
			_, state.setEncryptionErr = state.svc.PutBucketEncryption(&s3.PutBucketEncryptionInput{
				Bucket: aws.String(state.bucketName),
				ServerSideEncryptionConfiguration: &s3.ServerSideEncryptionConfiguration{
					Rules: []*s3.ServerSideEncryptionRule{
						{
							ApplyServerSideEncryptionByDefault: &s3.ServerSideEncryptionByDefault{
								SSEAlgorithm: aws.String("AES256"),
							},
						},
					},
				},
			})
			if state.setEncryptionErr != nil {
				log.Printf("Unable to set encryption on bucket %v due to %v", state.bucketName, state.setEncryptionErr)
				return state.setEncryptionErr
			}
		}
	}

	state.runningErr = err
	return nil
}

func (state *EncryptionAtRestAWS) createResult(result string) error {
	if result == "Success" {
		state.deleteCurrentTestBucket()
		return state.runningErr
	}
	// All other case
	// Fail
	if state.runningErr != nil {
		log.Printf("Bucket correctly fail in creation due to %v", state.runningErr)
		return nil
	}

	// Check if bucket has be remediated with server side encryption
	encrypted := state.checkBucketEncryption()
	count := 0
	for !encrypted {
		time.Sleep(sleepTime)
		encrypted = state.checkBucketEncryption()
		log.Printf("Try count: %d/%d Encryption options: %v", count, maxRetry, encrypted)
		count++
		if count >= maxRetry {
			break
		}
	}

	state.deleteCurrentTestBucket()

	if encrypted {
		return nil
	}

	return fmt.Errorf("bucket %v is not encrypted not self-healed", state.bucketName)
}

func (state *EncryptionAtRestAWS) checkBucketEncryption() bool {
	_, err := state.svc.GetBucketEncryption(&s3.GetBucketEncryptionInput{
		Bucket: aws.String(state.bucketName),
	})
	if err != nil {
		return false
	}

	return true
}

func (state *EncryptionAtRestAWS) deleteCurrentTestBucket() {
	_, err := state.svc.DeleteBucket(&s3.DeleteBucketInput{Bucket: aws.String(state.bucketName)})
	if err != nil {
		log.Printf("Error in deleting test bucket %v. Please manually clean up.", state.bucketName)
	} else {
		log.Printf("Bucket %v clean up successful.", state.bucketName)
	}
}
