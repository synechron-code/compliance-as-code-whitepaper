package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	citihubAws "citihub.com/compliance-as-code/internal/aws"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	notIPAddress  = "NotIpAddress"
	ipAddress     = "IpAddress"
	awsSourceVpc  = "aws:sourceVpc"
	awsSourceVpce = "aws:sourceVpce"
)

type accessWhitelistingAWS struct {
	ctx        context.Context
	tags       map[string]*string
	svc        *s3.S3
	session    *session.Session
	bucketName string
}

func (state *accessWhitelistingAWS) setup() {
	log.Println("[DEBUG] Setting up 'accessWhitelistingAWS'")
	state.ctx = context.Background()

	var err error
	state.session, err = session.NewSession()
	state.svc = s3.New(state.session)
	if err != nil {
		log.Fatalf("Unable create session to AWS due to %v", err)
	}
}

func (state *accessWhitelistingAWS) teardown() {
	log.Println("[DEBUG] Teardown completed")
}

func (state *accessWhitelistingAWS) checkPolicyAssigned() error {
	return fmt.Errorf("AWS does not support preventative controls for access whitelisting on S3")
}

func (state *accessWhitelistingAWS) provisionStorageContainer() error {
	// Not supported
	return nil
}

func (state *accessWhitelistingAWS) createWithWhitelist(arg1 string) error {
	// Not supported
	return nil
}

func (state *accessWhitelistingAWS) creationWill(arg1 string) error {
	// Not supported
	return nil
}

func (state *accessWhitelistingAWS) cspSupportsWhitelisting() error {
	return nil
}

func (state *accessWhitelistingAWS) examineStorageContainer(containerNameEnvVar string) error {

	name, b := os.LookupEnv(containerNameEnvVar)
	if !b {
		return fmt.Errorf("environment variable \"%s\" is not defined test can't run", containerNameEnvVar)
	}

	state.bucketName = name
	log.Printf("[DEBUG] Trying to access bucket: '%s'", state.bucketName)

	_, err := state.svc.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(state.bucketName),
	})

	if err != nil {
		return err
	}
	return nil
}

func (state *accessWhitelistingAWS) whitelistingIsConfigured() error {
	result, err := state.svc.GetBucketPolicy(&s3.GetBucketPolicyInput{
		Bucket: aws.String(state.bucketName),
	})
	if err != nil {
		return err
	}

	var policyDoc citihubAws.PolicyDocument
	log.Printf("[DEBUG] policy: %v", *result.Policy)
	err = json.Unmarshal([]byte(*result.Policy), &policyDoc)
	if err != nil {
		log.Panicf("%v", err)
		return err
	}
	for _, stmt := range *policyDoc.Statement {
		if *stmt.Effect == "Deny" {
			conditionMap := *stmt.Condition
			if conditionMap["StringNotEquals"] != nil {
				var conditionKey map[string]interface{}
				conditionKey = conditionMap["StringNotEquals"].(map[string]interface{})
				if conditionKey[awsSourceVpce] != nil {
					log.Printf("[DEBUG] %v: %v", awsSourceVpce, conditionKey[awsSourceVpce])
					return nil
				}
				if conditionKey[awsSourceVpc] != nil {
					log.Printf("[DEBUG] %v: %v", awsSourceVpc, conditionKey[awsSourceVpc])
					return nil
				}
			}
			if conditionMap[notIPAddress] != nil {
				log.Printf("[DEBUG] %v: %v", notIPAddress, conditionMap[notIPAddress])
				return nil
			}
		} else if *stmt.Effect == "Allow" {
			if stmt.Condition != nil {
				conditionMap := *stmt.Condition
				if conditionMap[ipAddress] != nil {
					log.Printf("[DEBUG] %v: %v", ipAddress, conditionMap[ipAddress])
					return nil
				}
			}
		}
	}

	return fmt.Errorf("no Deny IP address in bucket policy: %v", result)
}
