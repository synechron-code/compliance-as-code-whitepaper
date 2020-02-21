# Object Storage Compliance as code

## Essential controls required

* Secure transfer
* Data at rest encryption
* Restrict access to IP range
* 
## Limitation

Given the way different CSP provide compliance feature. Not all provide the same level of compliance effect. Below is a table stating the test-able level of controls for each CSP.

| Control description | AWS | Azure|
|---|---|---|
|Secure transfer (HTTPS) | Detective only | preventative |
|Data at rest encryption | Self Healing | By Design |
|Restrict access to IP range | ?* | preventative* |

*BDD not yet complete 

## Secure Transfer (HTTPS)

### AWS

S3 Bucket can be apply with a bucket policy that reject non-secure traffic. Unfortunately due to the way AWS Config Rule is designed, it can only evaluate if a resource is compliant (to ConfigRule `s3-bucket-ssl-requests-only`) post creation in a non-real time manner. Thus it can only serve as a **detective** control. There are currently no mechanism to *preventative* or *self-heal* from mis-configuration.

Alternatively, ConfigRule `s3-bucket-ssl-requests-only` can trigger auto remediation through SSM (AWS System Manager) automation to delete the bucket for non-compliant but it's less than idea.  There are no way to *self-heal* by applying the corresponding bucket policy as SSM cannot apply a bucket policy because I am unable to pass in a dynamic variable into the policy json in the auto remediation document.

Example run:
```
Feature: Object Storage Encryption at rest
  As a Cloud Security Architect
  I want to ensure that suitable security controls are applied to Object Storage
  So that my organisation is protected against data leakage via misconfiguration
2019/11/29 21:43:59 Azure Storage account is encrypted by default and cannot be turned off. No test to run.Checking Azure Policy. (Unless customise this test to check for specific key usage.

  Scenario: Ensure detective checks for Object Storage encryption at rest are enabled, when supported         # features\encryption_at_rest.feature:15
    Given the CSP provides a detective capability for unencrypted Object Storage containers # encryption_at_rest_test.go:18 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.policyOrRuleAvailable-fm
    When we examine the detective measure                                        # encryption_at_rest_test.go:19 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.checkPolicyOrRuleAssignment-fm
    Then the detective measure is enabled                                        # encryption_at_rest_test.go:20 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.policyOrRuleAssigned-fm
2019/11/29 21:43:59 Azure Storage account is encrypted by default and cannot be turned off. No test to run.Checking Azure Policy. (Unless customise this test to check for specific key usage.

  Scenario Outline: Prevent creation of Object Storage without encryption at rest                 # features\encryption_at_rest.feature:21
    Given security controls that enforce data at rest encryption for Object Storage are applied # encryption_at_rest_test.go:18 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.policyOrRuleAvailable-fm
    When we provision an Object Storage container                                              # encryption_at_rest_test.go:21 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.prepareToCreateContainer-fm
    And it is created with encryption option "<Encryption Option>"                             # encryption_at_rest_test.go:22 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.createContainerWithEncryptionOption-fm
    Then creation will "<Result>"                                                              # encryption_at_rest_test.go:23 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.createResult-fm

    Examples:
      | Encryption Option | Result                                   |
      | true              | Success                                  |
2019/11/29 21:43:59 Azure Storage account is encrypted by default and cannot be turned off. No test to run.Checking Azure Policy. (Unless customise this test to check for specific key usage.
      | false             | Fail or Success with encrytion turned on |

3 scenarios (3 passed)
11 steps (11 passed)
45.9748ms
```

### Azure

Azure policy can prevent the creation of Storage account that do not have secure transfer switch on.

*Explore whether it is possible to self-heal with policy "append" action.

Example run:
```
C:\work\compliance-as-code\test\features\general\object_storage\general\encryption_in_flight>go test
2019/11/29 21:50:53 Setting up "EncryptionInFlightAzure"
2019/11/29 21:50:53 'AZURE_POLICY_ASSIGNMENT_MANAGEMENT_GROUP' environment variable is not defined. Policy assignment check against subscription
2019/11/29 21:50:53 creating resource group 'testqwkidlresourecGP' on location: eastasia
2019/11/29 21:50:54 Created Resource Group: testqwkidlresourecGP
Feature: Object Storage Encryption in flight
  As a Cloud Security Architect
  I want to ensure that suitable security controls are applied to Object Storage
  So that my organisation is protected against data leakage via misconfiguration
2019/11/29 21:50:54 Getting Policy Assignment with subscriptionID: /subscriptions/56a69246-a48b-42fa-9068-55aead9d042d
2019/11/29 21:50:55 Policy assignment check: deny_http_storage [Step PASSED]

  Scenario Outline: Prevent creation of Object Storage without encryption at flight                   # features\encryption_in_flight.feature:15
2019/11/29 21:50:56 Expected create storage ac error: failed to start creating storage account: storage.AccountsClient#Create: Failure sending request: StatusCode=403 -- Original Error: Code="RequestDisallowedByPolicy" Message="Resource 'imywwstorageac' was disallowed by policy. Policy identifiers: '[{\"policyAssignment\":{\"name\":\"Secure transfer to storage accounts should be enabled\",\"id\":\"/subscriptions/56a69246-a48b-42fa-9068-55aead9d042d/providers/Microsoft.Authorization/policyAssignments/deny_http_storage\"},\"policyDefinition\":{\"name\":\"Secure transfer to storage accounts should be enabled\",\"id\":\"/providers/Microsoft.Authorization/policyDefinitions/404c3081-a854-4457-ae30-26a93ef643f9\"}}]'." Target="imywwstorageac" AdditionalInfo=[{"info":{"evaluationDetails":{"evaluatedExpressions":[{"expression":"type","expressionKind":"Field","expressionValue":"Microsoft.Storage/storageAccounts","operator":"Equals","path":"type","result":"True","targetValue":"Microsoft.Storage/storageAccounts"},{"expression":"Microsoft.Storage/storageAccounts/supportsHttpsTrafficOnly","expressionKind":"Field","expressionValue":false,"operator":"Equals","path":"properties.supportsHttpsTrafficOnly","result":"False","targetValue":"True"}]},"policyAssignmentDisplayName":"Secure transfer to storage accounts should be enabled","policyAssignmentId":"/subscriptions/56a69246-a48b-42fa-9068-55aead9d042d/providers/Microsoft.Authorization/policyAssignments/deny_http_storage","policyAssignmentName":"deny_http_storage","policyAssignmentParameters":{"effect":{"value":"Deny"}},"policyAssignmentScope":"/subscriptions/56a69246-a48b-42fa-9068-55aead9d042d","policyDefinitionDisplayName":"Secure transfer to storage accounts should be enabled","policyDefinitionEffect":"Deny","policyDefinitionId":"/providers/Microsoft.Authorization/policyDefinitions/404c3081-a854-4457-ae30-26a93ef643f9","policyDefinitionName":"404c3081-a854-4457-ae30-26a93ef643f9"},"type":"PolicyViolation"}]
    Given security controls that restrict data from being unencrypted in flight # encryption_in_flight_test.go:18 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.securityControlsThatRestrictDataFromBeingUnencryptedInFlight-fm
    When we provision an Object Storage bucket                                  # encryption_in_flight_test.go:19 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.weProvisionAnObjectStorageBucket-fm
    And http access is "<HTTP Option>"                                          # encryption_in_flight_test.go:20 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.httpAccessIs-fm
    And https access is "<HTTPS Option>"                                        # encryption_in_flight_test.go:21 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.httpsAccessIs-fm
    Then creation will "<Result>" with an error matching "<Error Description>"  # encryption_in_flight_test.go:22 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.creationWillWithAnErrorMatching-fm

    Examples:
      | HTTP Option | HTTPS Option | Result  | Error Description                                     |
      | enabled     | disabled     | Fail    | Storage Buckets must not be accessible via plain HTTP |
2019/11/29 21:50:56 Getting Policy Assignment with subscriptionID: /subscriptions/56a69246-a48b-42fa-9068-55aead9d042d
2019/11/29 21:50:57 Policy assignment check: deny_http_storage [Step PASSED]
2019/11/29 21:50:58 Expected create storage ac error: failed to start creating storage account: storage.AccountsClient#Create: Failure sending request: StatusCode=403 -- Original Error: Code="RequestDisallowedByPolicy" Message="Resource 'gkwgestorageac' was disallowed by policy. Policy identifiers: '[{\"policyAssignment\":{\"name\":\"Secure transfer to storage accounts should be enabled\",\"id\":\"/subscriptions/56a69246-a48b-42fa-9068-55aead9d042d/providers/Microsoft.Authorization/policyAssignments/deny_http_storage\"},\"policyDefinition\":{\"name\":\"Secure transfer to storage accounts should be enabled\",\"id\":\"/providers/Microsoft.Authorization/policyDefinitions/404c3081-a854-4457-ae30-26a93ef643f9\"}}]'." Target="gkwgestorageac" AdditionalInfo=[{"info":{"evaluationDetails":{"evaluatedExpressions":[{"expression":"type","expressionKind":"Field","expressionValue":"Microsoft.Storage/storageAccounts","operator":"Equals","path":"type","result":"True","targetValue":"Microsoft.Storage/storageAccounts"},{"expression":"Microsoft.Storage/storageAccounts/supportsHttpsTrafficOnly","expressionKind":"Field","expressionValue":false,"operator":"Equals","path":"properties.supportsHttpsTrafficOnly","result":"False","targetValue":"True"}]},"policyAssignmentDisplayName":"Secure transfer to storage accounts should be enabled","policyAssignmentId":"/subscriptions/56a69246-a48b-42fa-9068-55aead9d042d/providers/Microsoft.Authorization/policyAssignments/deny_http_storage","policyAssignmentName":"deny_http_storage","policyAssignmentParameters":{"effect":{"value":"Deny"}},"policyAssignmentScope":"/subscriptions/56a69246-a48b-42fa-9068-55aead9d042d","policyDefinitionDisplayName":"Secure transfer to storage accounts should be enabled","policyDefinitionEffect":"Deny","policyDefinitionId":"/providers/Microsoft.Authorization/policyDefinitions/404c3081-a854-4457-ae30-26a93ef643f9","policyDefinitionName":"404c3081-a854-4457-ae30-26a93ef643f9"},"type":"PolicyViolation"}]
      | enabled     | enabled      | Fail    | Storage Buckets must not be accessible via plain HTTP |
2019/11/29 21:50:58 Getting Policy Assignment with subscriptionID: /subscriptions/56a69246-a48b-42fa-9068-55aead9d042d
2019/11/29 21:50:59 Policy assignment check: deny_http_storage [Step PASSED]
      | disabled    | enabled      | Succeed |                                                       |

  Scenario: Ensure detective checks for Object Storage encryption in flight are enabled, when supported                                   # features\encryption_in_flight.feature:29
    Given the CSP provides a detective capability for unencrypted data transfer to Object Storage # encryption_in_flight_test.go:23 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.cSPProvideDetectiveMeasureForNonComplianceSecureTransferOnObjectStorage-fm
    When we examine the detective measure                                                    # encryption_in_flight_test.go:24 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.weExamineTheDetectiveMeasure-fm
2019/11/29 21:51:19 Getting Policy Assignment with subscriptionID: /subscriptions/56a69246-a48b-42fa-9068-55aead9d042d
2019/11/29 21:51:19 Policy assignment check: deny_http_storage [Step PASSED]
    Then the detective measure is enabled                                                    # encryption_in_flight_test.go:25 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.theDetectiveMeasureIsEnabled-fm
2019/11/29 21:51:19 deleting resources
2019/11/29 21:51:20 Teardown completed

4 scenarios (4 passed)
18 steps (18 passed)
26.4663439s
```

## Data at rest encryption

### AWS

AWS Config Rule `s3-bucket-server-side-encryption-enabled` can detect non compliant bucket and *self-heal* with SSM `AWS-EnableS3BucketEncryption` remediate action. However, this self-healing is not real time and have up to a few minutes latency. This make testing of self-healing test case require retry and wait. (default to 10 retry in 30 second interval)

Example run:

``` 
C:\work\compliance-as-code\test\features\general\object_storage\general\encryption_at_rest>go test
2019/11/29 21:32:07 Setting up "EncryptionAtRestAWS"
Feature: Object Storage Encryption at rest
  As a Cloud Security Architect
  I want to ensure that suitable security controls are applied to Object Storage
  So that my organisation is protected against data leakage via misconfiguration
2019/11/29 21:32:07 Checking AWS Config Rule: s3-bucket-server-side-encryption-enabled

  Scenario: Ensure detective checks for Object Storage encryption at rest are enabled, when supported         # features\encryption_at_rest.feature:15
    Given the CSP provides a detective capability for unencrypted Object Storage containers # encryption_at_rest_test.go:18 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.policyOrRuleAvailable-fm
    When we examine the detective measure                                        # encryption_at_rest_test.go:19 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.checkPolicyOrRuleAssignment-fm
2019/11/29 21:32:07 AWS Config Rule: "s3-bucket-server-side-encryption-enabled" evaluation results count: 10
    Then the detective measure is enabled                                        # encryption_at_rest_test.go:20 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.policyOrRuleAssigned-fm
2019/11/29 21:32:07 Checking AWS Config Rule: s3-bucket-server-side-encryption-enabled

  Scenario Outline: Prevent creation of Object Storage without encryption at rest                 # features\encryption_at_rest.feature:21
2019/11/29 21:32:09 Created Bucket: {
  Location: "http://testsvlofbucketenc.s3.amazonaws.com/"
}
    Given security controls that enforce data at rest encryption for Object Storage are applied # encryption_at_rest_test.go:18 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.policyOrRuleAvailable-fm
    When we provision an Object Storage container                                              # encryption_at_rest_test.go:21 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.prepareToCreateContainer-fm
    And it is created with encryption option "<Encryption Option>"                             # encryption_at_rest_test.go:22 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.createContainerWithEncryptionOption-fm
    Then creation will "<Result>"                                                              # encryption_at_rest_test.go:23 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.createResult-fm

    Examples:
      | Encryption Option | Result                                   |
      | true              | Success                                  |
2019/11/29 21:32:09 Checking AWS Config Rule: s3-bucket-server-side-encryption-enabled
2019/11/29 21:32:11 Created Bucket: {
  Location: "http://testvulkdbucketunenc.s3.amazonaws.com/"
}
2019/11/29 21:32:41 Try count: 0/10 Encryption options: false
2019/11/29 21:33:12 Try count: 1/10 Encryption options: false
2019/11/29 21:33:42 Try count: 2/10 Encryption options: false
2019/11/29 21:34:12 Try count: 3/10 Encryption options: false
2019/11/29 21:34:42 Try count: 4/10 Encryption options: false
2019/11/29 21:35:12 Try count: 5/10 Encryption options: true
2019/11/29 21:35:14 Bucket testvulkdbucketunenc clean up successful.
      | false             | Fail or Success with encrytion turned on |
2019/11/29 21:35:14 Teardown completed

3 scenarios (3 passed)
11 steps (11 passed)
3m6.9074817s
```

### Azure

Encryption is by default and cannot be turned off.

Example run:
```
Feature: Object Storage Encryption at rest
  As a Cloud Security Architect
  I want to ensure that suitable security controls are applied to Object Storage
  So that my organisation is protected against data leakage via misconfiguration
2019/11/29 21:43:59 Azure Storage account is encrypted by default and cannot be turned off. No test to run.Checking Azure Policy. (Unless customise this test to check for specific key usage.

  Scenario: Ensure detective checks for Object Storage encryption at rest are enabled, when supported         # features\encryption_at_rest.feature:15
    Given the CSP provides a detective capability for unencrypted Object Storage containers # encryption_at_rest_test.go:18 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.policyOrRuleAvailable-fm
    When we examine the detective measure                                        # encryption_at_rest_test.go:19 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.checkPolicyOrRuleAssignment-fm
    Then the detective measure is enabled                                        # encryption_at_rest_test.go:20 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.policyOrRuleAssigned-fm
2019/11/29 21:43:59 Azure Storage account is encrypted by default and cannot be turned off. No test to run.Checking Azure Policy. (Unless customise this test to check for specific key usage.

  Scenario Outline: Prevent creation of Object Storage without encryption at rest                 # features\encryption_at_rest.feature:21
    Given security controls that enforce data at rest encryption for Object Storage are applied # encryption_at_rest_test.go:18 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.policyOrRuleAvailable-fm
    When we provision an Object Storage container                                              # encryption_at_rest_test.go:21 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.prepareToCreateContainer-fm
    And it is created with encryption option "<Encryption Option>"                             # encryption_at_rest_test.go:22 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.createContainerWithEncryptionOption-fm
    Then creation will "<Result>"                                                              # encryption_at_rest_test.go:23 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.createResult-fm

    Examples:
      | Encryption Option | Result                                   |
      | true              | Success                                  |
2019/11/29 21:43:59 Azure Storage account is encrypted by default and cannot be turned off. No test to run.Checking Azure Policy. (Unless customise this test to check for specific key usage.
      | false             | Fail or Success with encrytion turned on |

3 scenarios (3 passed)
11 steps (11 passed)
45.9748ms
testing: warning: no tests to run
PASS
ok      citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest  0.478s
```

## Restrict access to IP range

### AWS
**TODO**
AWS Config Rule `s3-bucket-policy-grantee-check` can detect non-compliant bucket policy but it show compliant with no bucket policy. Not suit for purpose.

