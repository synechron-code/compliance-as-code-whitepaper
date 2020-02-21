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

// AccessWhitelistingAWS Azure implementation of the encryption in flight for Object Storage feature
type AccessWhitelistingAWS struct {
	ctx        context.Context
	tags       map[string]*string
	svc        *s3.S3
	session    *session.Session
	bucketName string
}

func (state *AccessWhitelistingAWS) setup() {
	log.Println("Setting up \"AccessWhitelistingAWS\"")
	state.ctx = context.Background()

	// Create Session
	var err error
	state.session, err = session.NewSession()
	state.svc = s3.New(state.session)
	if err != nil {
		log.Fatalf("Unable create session to AWS due to %v", err)
	}
}

func (state *AccessWhitelistingAWS) teardown() {
	log.Println("Teardown completed")
}

func (state *AccessWhitelistingAWS) checkPolicyAssigned() error {
	return fmt.Errorf("AWS do not support preventative controls for access whitelisting on S3")
}

func (state *AccessWhitelistingAWS) prepareToCreateStorageContainer() error {
	// Not supported
	return nil
}

func (state *AccessWhitelistingAWS) createWithWhiteList(arg1 string) error {
	// Not supported
	return nil
}

func (state *AccessWhitelistingAWS) creationWill(arg1 string) error {
	// Not supported
	return nil
}

func (state *AccessWhitelistingAWS) isCspCapable() error {
	return nil
}

func (state *AccessWhitelistingAWS) examineStorageContainer(containerNameEnvVar string) error {
	containerName := os.Getenv(containerNameEnvVar)
	if containerName == "" {
		return fmt.Errorf("environment variable \"%s\" is not defined test can't run", containerNameEnvVar)
	}

	state.bucketName = containerName

	_, err := state.svc.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(state.bucketName),
	})

	if err != nil {
		return err
	}
	return nil
}

func (state *AccessWhitelistingAWS) whitelistingIsConfigured() error {
	result, err := state.svc.GetBucketPolicy(&s3.GetBucketPolicyInput{
		Bucket: aws.String(state.bucketName),
	})
	if err != nil {
		return err
	}

	var policyDoc citihubAws.PolicyDocument
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
					log.Printf("%v: %v", awsSourceVpce, conditionKey[awsSourceVpce])
					return nil
				}
				if conditionKey[awsSourceVpc] != nil {
					log.Printf("%v: %v", awsSourceVpc, conditionKey[awsSourceVpc])
					return nil
				}
			}
			if conditionMap[notIPAddress] != nil {
				log.Printf("%v: %v", notIPAddress, conditionMap[notIPAddress])
				return nil
			}
		} else if *stmt.Effect == "Allow" {
			conditionMap := *stmt.Condition
			if conditionMap[ipAddress] != nil {
				log.Printf("%v: %v", ipAddress, conditionMap[ipAddress])
				return nil
			}
		}
	}

	return fmt.Errorf("no Deny IP address in bucket policy: %v", result)
}
