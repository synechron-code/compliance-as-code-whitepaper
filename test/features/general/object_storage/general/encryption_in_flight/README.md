# Encryption in Flight

## AWS

### Implementation Details
An S3 Bucket can be created with a bucket policy that rejects insecure (HTTP) traffic. 

In our AWS example, we attest that the AWS ConfigRule `s3-bucket-ssl-requests-only` is in place and correctly detects that S3 Bucket we deploy as part of the test is non-compliant (AWS Config evaluates this post-creation, out-of-band).

We also configure the ConfigRule `s3-bucket-ssl-requests-only` to trigger auto remediation through SSM (AWS System Manager) and we attest that this has also happened once the detection event has occurred.

### Example Run:
```
>go test
2020/03/11 11:12:28 [DEBUG] Setting up "EncryptionInFlightAWS"
Feature: Object Storage Encryption in Flight
  As a Cloud Security Architect
  I want to ensure that suitable security controls are applied to Object Storage
  So that my organisation is not vulnerable to interception of data in transit

  Rule: CHC2-AGP140 - Ensure cryptographic controls are in place to protect the confidentiality and integrity of data in-transit, stored, generated and processed in the cloud

  Scenario Outline: Prevent Creation of Object Storage Without Encryption in Flight # features\encryption_in_flight.feature:17
    Given security controls that restrict data from being unencrypted in flight     # encryption_in_flight_test.go:18 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.securityControlsThatRestrictDataFromBeingUnencryptedInFlight-fm
    When we provision an Object Storage bucket                                      # encryption_in_flight_test.go:19 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.weProvisionAnObjectStorageBucket-fm
    And http access is "<HTTP Option>"                                              # encryption_in_flight_test.go:20 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.httpAccessIs-fm
    And https access is "<HTTPS Option>"                                            # encryption_in_flight_test.go:21 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.httpsAccessIs-fm
    Then creation will "<Result>" with an error matching "<Error Description>"      # encryption_in_flight_test.go:22 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.creationWillWithAnErrorMatching-fm

    Examples:
      | HTTP Option | HTTPS Option | Result  | Error Description                                     |
      | enabled     | disabled     | Fail    | Storage Buckets must not be accessible via plain HTTP |
        AWS do not support preventative controls for secure transfer on S3
      | enabled     | enabled      | Fail    | Storage Buckets must not be accessible via plain HTTP |
        AWS do not support preventative controls for secure transfer on S3
      | disabled    | enabled      | Succeed |                                                       |
        AWS do not support preventative controls for secure transfer on S3

  Scenario: Remediate Object Storage if Creation of Object Storage Without Encryption in Flight is Detected          # features\encryption_in_flight.feature:31
    Given there is a detective capability for creation of Object Storage with unencrypted data transfer enabled      # encryption_in_flight_test.go:24 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.detectObjectStorageUnencryptedTransferAvailable-fm
    And the capability for detecting the creation of Object Storage with unencrypted data transfer enabled is active # encryption_in_flight_test.go:25 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.detectObjectStorageUnencryptedTransferEnabled-fm
2020/03/11 11:12:30 [DEBUG] Created Bucket: {
  Location: "http://testdzbotunencbucket.s3.amazonaws.com/"
}
    When Object Storage is created with unencrypted data transfer enabled                                            # encryption_in_flight_test.go:26 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.createUnencryptedTransferObjectStorage-fm
2020/03/11 11:12:30 [DEBUG] Waiting for bucket to be detected by Config Rule...
2020/03/11 11:12:30 [DEBUG] nextToken: <nil>
2020/03/11 11:12:30 [DEBUG] AWS Config Rule: "s3-bucket-ssl-requests-only" evaluation results count: 19
2020/03/11 11:12:30 [DEBUG] Config Rule not pick up bucket 'testdzbotunencbucket' yet wait for 30 s, retry 0/10
2020/03/11 11:13:00 [DEBUG] nextToken: <nil>
2020/03/11 11:13:00 [DEBUG] AWS Config Rule: "s3-bucket-ssl-requests-only" evaluation results count: 19
2020/03/11 11:13:00 [DEBUG] Config Rule not pick up bucket 'testdzbotunencbucket' yet wait for 30 s, retry 1/10
2020/03/11 11:13:31 [DEBUG] nextToken: <nil>
2020/03/11 11:13:31 [DEBUG] AWS Config Rule: "s3-bucket-ssl-requests-only" evaluation results count: 19
2020/03/11 11:13:31 [DEBUG] Config Rule not pick up bucket 'testdzbotunencbucket' yet wait for 30 s, retry 2/10
2020/03/11 11:14:01 [DEBUG] nextToken: <nil>
2020/03/11 11:14:01 [DEBUG] AWS Config Rule: "s3-bucket-ssl-requests-only" evaluation results count: 19
2020/03/11 11:14:01 [DEBUG] Config Rule not pick up bucket 'testdzbotunencbucket' yet wait for 30 s, retry 3/10
2020/03/11 11:14:31 [DEBUG] nextToken: <nil>
2020/03/11 11:14:31 [DEBUG] AWS Config Rule: "s3-bucket-ssl-requests-only" evaluation results count: 19
2020/03/11 11:14:31 [DEBUG] Config Rule not pick up bucket 'testdzbotunencbucket' yet wait for 30 s, retry 4/10
2020/03/11 11:15:01 [DEBUG] nextToken: <nil>
2020/03/11 11:15:01 [DEBUG] AWS Config Rule: "s3-bucket-ssl-requests-only" evaluation results count: 19
2020/03/11 11:15:01 [DEBUG] Config Rule not pick up bucket 'testdzbotunencbucket' yet wait for 30 s, retry 5/10
2020/03/11 11:15:31 [DEBUG] nextToken: <nil>
2020/03/11 11:15:31 [DEBUG] AWS Config Rule: "s3-bucket-ssl-requests-only" evaluation results count: 20
2020/03/11 11:15:31 [DEBUG] Bucket 'testdzbotunencbucket' is 'NON_COMPLIANT'
    Then the detective capability detects the creation of Object Storage with unencrypted data transfer enabled      # encryption_in_flight_test.go:27 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.detectsTheObjectStorage-fm
2020/03/11 11:15:31 [DEBUG] Checking bucket policy for secure transport setting...
2020/03/11 11:15:32 [DEBUG] Bucket policy still insecure wait for 30 s, retry 0/10
2020/03/11 11:16:02 [DEBUG] Checking bucket policy for secure transport setting...
2020/03/11 11:16:02 [DEBUG] Bucket policy still insecure wait for 30 s, retry 1/10
2020/03/11 11:16:32 [DEBUG] Checking bucket policy for secure transport setting...
2020/03/11 11:16:32 [DEBUG] aws:SecureTransport: false
    And the detective capability enforces encrypted data transfer on the Object Storage Bucket                       # encryption_in_flight_test.go:28 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.encryptedDataTrafficIsEnforced-fm
2020/03/11 11:16:33 [DEBUG] Bucket testdzbotunencbucket clean up successful.
2020/03/11 11:16:33 [DEBUG] Teardown completed

--- Failed steps:

  Scenario Outline: Prevent Creation of Object Storage Without Encryption in Flight # features\encryption_in_flight.feature:17
    Given security controls that restrict data from being unencrypted in flight # features\encryption_in_flight.feature:18
      Error: AWS do not support preventative controls for secure transfer on S3

  Scenario Outline: Prevent Creation of Object Storage Without Encryption in Flight # features\encryption_in_flight.feature:17
    Given security controls that restrict data from being unencrypted in flight # features\encryption_in_flight.feature:18
      Error: AWS do not support preventative controls for secure transfer on S3

  Scenario Outline: Prevent Creation of Object Storage Without Encryption in Flight # features\encryption_in_flight.feature:17
    Given security controls that restrict data from being unencrypted in flight # features\encryption_in_flight.feature:18
      Error: AWS do not support preventative controls for secure transfer on S3


4 scenarios (1 passed, 3 failed)
20 steps (5 passed, 3 failed, 12 skipped)
4m5.0008202s
testing: warning: no tests to run
PASS
exit status 1
FAIL    citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight        245.722s
```

## Azure

### Implementation Details
In our Azure example, we attest that Azure Policy is in place which out-right prevents the creation of a Storage Account that does not have the secure transfer switch turned on, as an in-band evaluation.

In our `Given` clause we expect that the Azure Built-in policy `Secure transfer to storage accounts should be enabled` has been assigned on the subscription with `Deny` Effect (see the [terraform example](../../../../../../terraform/modules/policies/deny_http_storage)). When we attempt to create a storage account, it should prevent the creation request if the `supportHttpsTrafficOnly` field is false.

### Example Run
```
>go test
2020/03/11 11:06:16 [DEBUG] Setting up "EncryptionInFlightAzure"
2020/03/11 11:06:16 [DEBUG] creating Resource Group 'testbzynnqresourceGP' on location: 'eastasia'
2020/03/11 11:06:17 [DEBUG] Created Resource Group: 'testbzynnqresourceGP'
Feature: Object Storage Encryption in Flight
  As a Cloud Security Architect
  I want to ensure that suitable security controls are applied to Object Storage
  So that my organisation is not vulnerable to interception of data in transit

  Rule: CHC2-AGP140 - Ensure cryptographic controls are in place to protect the confidentiality and integrity of data in-transit, stored, generated and processed in the cloud
2020/03/11 11:06:17 [DEBUG] Getting Policy Assignment with scope: /providers/Microsoft.Management/managementGroups/boxbank-root
2020/03/11 11:06:17 [DEBUG] Policy assignment check: deny_http_storage [Step PASSED]

  Scenario Outline: Prevent Creation of Object Storage Without Encryption in Flight # features\encryption_in_flight.feature:17
2020/03/11 11:06:17 [DEBUG] Creating Storage Account with HTTPS: false
2020/03/11 11:06:18 [DEBUG] Detailed Error: Code="RequestDisallowedByPolicy" Message="Resource 'lpljxstorageac' was disallowed by policy. Policy identifiers: '[{\"policyAssignment\":{\"name\":\"Secure transfer to storage accounts should be enabled\",\"id\":\"/providers/Microsoft.Management/managementgroups/boxbank-root/providers/Microsoft.Authorization/policyAssignments/deny_http_storage\"},\"policyDefinition\":{\"name\":\"Secure transfer to storage accounts should be enabled\",\"id\":\"/providers/Microsoft.Authorization/policyDefinitions/404c3081-a854-4457-ae30-26a93ef643f9\"}}]'." Target="lpljxstorageac" AdditionalInfo=[{"info":{"evaluationDetails":{"evaluatedExpressions":[{"expression":"type","expressionKind":"Field","expressionValue":"Microsoft.Storage/storageAccounts","operator":"Equals","path":"type","result":"True","targetValue":"Microsoft.Storage/storageAccounts"},{"expression":"Microsoft.Storage/storageAccounts/supportsHttpsTrafficOnly","expressionKind":"Field","expressionValue":false,"operator":"Equals","path":"properties.supportsHttpsTrafficOnly","result":"False","targetValue":"True"}]},"policyAssignmentDisplayName":"Secure transfer to storage accounts should be enabled","policyAssignmentId":"/providers/Microsoft.Management/managementgroups/boxbank-root/providers/Microsoft.Authorization/policyAssignments/deny_http_storage","policyAssignmentName":"deny_http_storage","policyAssignmentParameters":{"effect":{"value":"Deny"}},"policyAssignmentScope":"/providers/Microsoft.Management/managementgroups/boxbank-root","policyDefinitionDisplayName":"Secure transfer to storage accounts should be enabled","policyDefinitionEffect":"Deny","policyDefinitionId":"/providers/Microsoft.Authorization/policyDefinitions/404c3081-a854-4457-ae30-26a93ef643f9","policyDefinitionName":"404c3081-a854-4457-ae30-26a93ef643f9"},"type":"PolicyViolation"}]
2020/03/11 11:06:18 [DEBUG] Request was Disallowed By Policy: deny_http_storage [Step PASSED]
    Given security controls that restrict data from being unencrypted in flight     # encryption_in_flight_test.go:18 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.securityControlsThatRestrictDataFromBeingUnencryptedInFlight-fm
    When we provision an Object Storage bucket                                      # encryption_in_flight_test.go:19 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.weProvisionAnObjectStorageBucket-fm
    And http access is "<HTTP Option>"                                              # encryption_in_flight_test.go:20 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.httpAccessIs-fm
    And https access is "<HTTPS Option>"                                            # encryption_in_flight_test.go:21 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.httpsAccessIs-fm
    Then creation will "<Result>" with an error matching "<Error Description>"      # encryption_in_flight_test.go:22 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.creationWillWithAnErrorMatching-fm

    Examples:
      | HTTP Option | HTTPS Option | Result  | Error Description                                     |
      | enabled     | disabled     | Fail    | Storage Buckets must not be accessible via plain HTTP |
2020/03/11 11:06:18 [DEBUG] Getting Policy Assignment with scope: /providers/Microsoft.Management/managementGroups/boxbank-root
2020/03/11 11:06:19 [DEBUG] Policy assignment check: deny_http_storage [Step PASSED]
2020/03/11 11:06:19 [DEBUG] Creating Storage Account with HTTPS: false
2020/03/11 11:06:20 [DEBUG] Detailed Error: Code="RequestDisallowedByPolicy" Message="Resource 'mzbdxstorageac' was disallowed by policy. Policy identifiers: '[{\"policyAssignment\":{\"name\":\"Secure transfer to storage accounts should be enabled\",\"id\":\"/providers/Microsoft.Management/managementgroups/boxbank-root/providers/Microsoft.Authorization/policyAssignments/deny_http_storage\"},\"policyDefinition\":{\"name\":\"Secure transfer to storage accounts should be enabled\",\"id\":\"/providers/Microsoft.Authorization/policyDefinitions/404c3081-a854-4457-ae30-26a93ef643f9\"}}]'." Target="mzbdxstorageac" AdditionalInfo=[{"info":{"evaluationDetails":{"evaluatedExpressions":[{"expression":"type","expressionKind":"Field","expressionValue":"Microsoft.Storage/storageAccounts","operator":"Equals","path":"type","result":"True","targetValue":"Microsoft.Storage/storageAccounts"},{"expression":"Microsoft.Storage/storageAccounts/supportsHttpsTrafficOnly","expressionKind":"Field","expressionValue":false,"operator":"Equals","path":"properties.supportsHttpsTrafficOnly","result":"False","targetValue":"True"}]},"policyAssignmentDisplayName":"Secure transfer to storage accounts should be enabled","policyAssignmentId":"/providers/Microsoft.Management/managementgroups/boxbank-root/providers/Microsoft.Authorization/policyAssignments/deny_http_storage","policyAssignmentName":"deny_http_storage","policyAssignmentParameters":{"effect":{"value":"Deny"}},"policyAssignmentScope":"/providers/Microsoft.Management/managementgroups/boxbank-root","policyDefinitionDisplayName":"Secure transfer to storage accounts should be enabled","policyDefinitionEffect":"Deny","policyDefinitionId":"/providers/Microsoft.Authorization/policyDefinitions/404c3081-a854-4457-ae30-26a93ef643f9","policyDefinitionName":"404c3081-a854-4457-ae30-26a93ef643f9"},"type":"PolicyViolation"}]
2020/03/11 11:06:20 [DEBUG] Request was Disallowed By Policy: deny_http_storage [Step PASSED]
      | enabled     | enabled      | Fail    | Storage Buckets must not be accessible via plain HTTP |
2020/03/11 11:06:20 [DEBUG] Getting Policy Assignment with scope: /providers/Microsoft.Management/managementGroups/boxbank-root
2020/03/11 11:06:21 [DEBUG] Policy assignment check: deny_http_storage [Step PASSED]
2020/03/11 11:06:21 [DEBUG] Creating Storage Account with HTTPS: true
      | disabled    | enabled      | Succeed |                                                       |

  Scenario: Remediate Object Storage if Creation of Object Storage Without Encryption in Flight is Detected          # features\encryption_in_flight.feature:31
    Given there is a detective capability for creation of Object Storage with unencrypted data transfer enabled      # encryption_in_flight_test.go:24 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.detectObjectStorageUnencryptedTransferAvailable-fm
    azure policy prevent creation of object storage with in-secure transport making this test irrelevant
    And the capability for detecting the creation of Object Storage with unencrypted data transfer enabled is active # encryption_in_flight_test.go:25 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.detectObjectStorageUnencryptedTransferEnabled-fm
    When Object Storage is created with unencrypted data transfer enabled                                            # encryption_in_flight_test.go:26 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.createUnencryptedTransferObjectStorage-fm
    Then the detective capability detects the creation of Object Storage with unencrypted data transfer enabled      # encryption_in_flight_test.go:27 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.detectsTheObjectStorage-fm
    And the detective capability enforces encrypted data transfer on the Object Storage Bucket                       # encryption_in_flight_test.go:28 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight.EncryptionInFlight.encryptedDataTrafficIsEnforced-fm
2020/03/11 11:06:40 [DEBUG] Deleting resources
2020/03/11 11:06:41 [DEBUG] Teardown completed

--- Failed steps:

  Scenario: Remediate Object Storage if Creation of Object Storage Without Encryption in Flight is Detected # features\encryption_in_flight.feature:31
    Given there is a detective capability for creation of Object Storage with unencrypted data transfer enabled # features\encryption_in_flight.feature:32
      Error: azure policy prevent creation of object storage with in-secure transport making this test irrelevant


4 scenarios (3 passed, 1 failed)
20 steps (15 passed, 1 failed, 4 skipped)
24.6034735s
testing: warning: no tests to run
PASS
exit status 1
FAIL    citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_in_flight        25.310s
```