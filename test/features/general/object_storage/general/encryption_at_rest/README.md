# Encryption at Rest

## AWS

### Implementation Details

In our AWS example we demonstrate a "*chaos*" test, by deploying an S3 bucket which violates our policy and attests that the expected behaviour occurs.

We attest that AWS SSM `AWS-EnableS3BucketEncryption` *auto-remediate* action is in place and that this action indeed does result in auto-remediation of the non-compliant S3 Bucket we deploy as part of the test. This self-healing occurs out-of-band with a few minutes of latency, therefore testing of the self-healing case requires retry and wait. (default to 10 retry in 30 second interval)

### Example Run

``` 
>go test
2020/03/11 11:59:00 [DEBUG] Setting up "EncryptionAtRestAWS"
Feature: Object Storage Encryption at Rest
  As a Cloud Security Architect
  I want to ensure that suitable security controls are applied to Object Storage
  So that my organisation is protected against data leakage due to misconfiguration

  Rule: CHC2-AGP140 - Ensure cryptographic controls are in place to protect the confidentiality and integrity of data in-transit, stored, generated and processed in the cloud

  Scenario Outline: Prevent Creation of Object Storage Without Encryption at Rest # features\encryption_at_rest.feature:18
    Given security controls that restrict data from being unencrypted at rest     # encryption_at_rest_test.go:19 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.securityControlsThatRestrictDataFromBeingUnencryptedAtRest-fm
    When we provision an Object Storage bucket                                    # encryption_at_rest_test.go:20 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.weProvisionAnObjectStorageBucket-fm
    And encryption at rest is "<Encryption Option>"                               # encryption_at_rest_test.go:21 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.encryptionAtRestIs-fm
    Then creation will "<Result>" with an error matching "<Error Description>"    # encryption_at_rest_test.go:22 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.creationWillWithAnErrorMatching-fm

    Examples:
      | Encryption Option | Result  | Error Description                                                      |
      | enabled           | Fail    | Storage Buckets must not be created without encryption as rest enabled |
        AWS do not have preventive measure but instead reliant on detective measure
      | disabled          | Succeed |                                                                        |
        AWS do not have preventive measure but instead reliant on detective measure
2020/03/11 11:59:00 [DEBUG] Checking AWS Config Rule: s3-bucket-server-side-encryption-enabled

  Scenario: Detect creation of Object Storage Without Encryption at Rest                                 # features\encryption_at_rest.feature:30
    Given there is a detective capability for creation of Object Storage without encryption at rest      # encryption_at_rest_test.go:23 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.policyOrRuleAvailable-fm
    And the capability for detecting the creation of Object Storage without encryption at rest is active # encryption_at_rest_test.go:24 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.checkPolicyOrRuleAssignment-fm
2020/03/11 11:59:03 [DEBUG] Created Bucket: http://testvduisunencbucket.s3.amazonaws.com/
    When Object Storage is created with without encryption at rest                                       # encryption_at_rest_test.go:27 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.createContainerWithoutEncryption-fm
2020/03/11 11:59:03 [DEBUG] Waiting for bucket to be detected by Config Rule...
2020/03/11 11:59:03 [DEBUG] nextToken: <nil>
2020/03/11 11:59:03 [DEBUG] AWS Config Rule: "s3-bucket-server-side-encryption-enabled" evaluation results count: 19
2020/03/11 11:59:03 [DEBUG] Config Rule not pick up bucket 'testvduisunencbucket' yet wait for 30 s, retry 0/10
2020/03/11 11:59:33 [DEBUG] nextToken: <nil>
2020/03/11 11:59:33 [DEBUG] AWS Config Rule: "s3-bucket-server-side-encryption-enabled" evaluation results count: 19
2020/03/11 11:59:33 [DEBUG] Config Rule not pick up bucket 'testvduisunencbucket' yet wait for 30 s, retry 1/10
2020/03/11 12:00:03 [DEBUG] nextToken: <nil>
2020/03/11 12:00:03 [DEBUG] AWS Config Rule: "s3-bucket-server-side-encryption-enabled" evaluation results count: 19
2020/03/11 12:00:03 [DEBUG] Config Rule not pick up bucket 'testvduisunencbucket' yet wait for 30 s, retry 2/10
2020/03/11 12:00:33 [DEBUG] nextToken: <nil>
2020/03/11 12:00:33 [DEBUG] AWS Config Rule: "s3-bucket-server-side-encryption-enabled" evaluation results count: 19
2020/03/11 12:00:33 [DEBUG] Config Rule not pick up bucket 'testvduisunencbucket' yet wait for 30 s, retry 3/10
2020/03/11 12:01:04 [DEBUG] nextToken: <nil>
2020/03/11 12:01:04 [DEBUG] AWS Config Rule: "s3-bucket-server-side-encryption-enabled" evaluation results count: 19
2020/03/11 12:01:04 [DEBUG] Config Rule not pick up bucket 'testvduisunencbucket' yet wait for 30 s, retry 4/10
2020/03/11 12:01:34 [DEBUG] nextToken: <nil>
2020/03/11 12:01:34 [DEBUG] AWS Config Rule: "s3-bucket-server-side-encryption-enabled" evaluation results count: 19
2020/03/11 12:01:34 [DEBUG] Config Rule not pick up bucket 'testvduisunencbucket' yet wait for 30 s, retry 5/10
2020/03/11 12:02:04 [DEBUG] nextToken: <nil>
2020/03/11 12:02:04 [DEBUG] AWS Config Rule: "s3-bucket-server-side-encryption-enabled" evaluation results count: 20
2020/03/11 12:02:04 [DEBUG] Bucket 'testvduisunencbucket' is 'NON_COMPLIANT'
    Then the detective capability detects the creation of Object Storage without encryption at rest      # encryption_at_rest_test.go:28 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.detectiveDetectsNonCompliant-fm
2020/03/11 12:02:04 [DEBUG] Checking bucket policy for SSE setting
2020/03/11 12:02:04 [DEBUG] Bucket policy still unencrypted wait for 30 s, retry 0/10
2020/03/11 12:02:34 [DEBUG] Checking bucket policy for SSE setting
2020/03/11 12:02:34 [DEBUG] Bucket policy still unencrypted wait for 30 s, retry 1/10
2020/03/11 12:03:04 [DEBUG] Checking bucket policy for SSE setting
    And the detective capability enforces encryption at rest on the Object Storage Bucket                # encryption_at_rest_test.go:29 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.containerIsRemediated-fm
2020/03/11 12:03:05 [DEBUG] Bucket testvduisunencbucket clean up successful.
2020/03/11 12:03:05 [DEBUG] Teardown completed

--- Failed steps:

  Scenario Outline: Prevent Creation of Object Storage Without Encryption at Rest # features\encryption_at_rest.feature:18
    Given security controls that restrict data from being unencrypted at rest # features\encryption_at_rest.feature:19
      Error: AWS do not have preventive measure but instead reliant on detective measure

  Scenario Outline: Prevent Creation of Object Storage Without Encryption at Rest # features\encryption_at_rest.feature:18
    Given security controls that restrict data from being unencrypted at rest # features\encryption_at_rest.feature:19
      Error: AWS do not have preventive measure but instead reliant on detective measure


3 scenarios (1 passed, 2 failed)
13 steps (5 passed, 2 failed, 6 skipped)
4m5.1350883s
testing: warning: no tests to run
PASS
exit status 1
FAIL    citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest  245.771s
```

## Azure

### Implementation Details

In Azure, encryption-at-rest is always-on by default and cannot be turned off.  We decided just to pass the test by default.

### Example Run
```
>go test
Feature: Object Storage Encryption at Rest
  As a Cloud Security Architect
  I want to ensure that suitable security controls are applied to Object Storage
  So that my organisation is protected against data leakage due to misconfiguration

  Rule: CHC2-AGP140 - Ensure cryptographic controls are in place to protect the confidentiality and integrity of data in-transit, stored, generated and processed in the cloud
2020/03/11 12:05:19 [DEBUG] Azure Storage account is encrypted by default and cannot be turned off. No test to run. Checking Azure Policy. (Unless customise this test to check for specific key usage.

  Scenario Outline: Prevent Creation of Object Storage Without Encryption at Rest # features\encryption_at_rest.feature:18
    Given security controls that restrict data from being unencrypted at rest     # encryption_at_rest_test.go:19 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.securityControlsThatRestrictDataFromBeingUnencryptedAtRest-fm
    When we provision an Object Storage bucket                                    # encryption_at_rest_test.go:20 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.weProvisionAnObjectStorageBucket-fm
    And encryption at rest is "<Encryption Option>"                               # encryption_at_rest_test.go:21 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.encryptionAtRestIs-fm
    Then creation will "<Result>" with an error matching "<Error Description>"    # encryption_at_rest_test.go:22 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.creationWillWithAnErrorMatching-fm

    Examples:
      | Encryption Option | Result  | Error Description                                                      |
      | enabled           | Fail    | Storage Buckets must not be created without encryption as rest enabled |
2020/03/11 12:05:19 [DEBUG] Azure Storage account is encrypted by default and cannot be turned off. No test to run. Checking Azure Policy. (Unless customise this test to check for specific key usage.
      | disabled          | Succeed |                                                                        |
2020/03/11 12:05:19 [DEBUG] Azure Storage account is encrypted by default and cannot be turned off. No test to run. Checking Azure Policy. (Unless customise this test to check for specific key usage.

  Scenario: Detect creation of Object Storage Without Encryption at Rest                                 # features\encryption_at_rest.feature:30
    Given there is a detective capability for creation of Object Storage without encryption at rest      # encryption_at_rest_test.go:23 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.policyOrRuleAvailable-fm
    And the capability for detecting the creation of Object Storage without encryption at rest is active # encryption_at_rest_test.go:24 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.checkPolicyOrRuleAssignment-fm
    When Object Storage is created with without encryption at rest                                       # encryption_at_rest_test.go:27 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.createContainerWithoutEncryption-fm
    Then the detective capability detects the creation of Object Storage without encryption at rest      # encryption_at_rest_test.go:28 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.detectiveDetectsNonCompliant-fm
    And the detective capability enforces encryption at rest on the Object Storage Bucket                # encryption_at_rest_test.go:29 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest.EncryptionAtRest.containerIsRemediated-fm

3 scenarios (3 passed)
13 steps (13 passed)
53.9682ms
testing: warning: no tests to run
PASS
ok      citihub.com/compliance-as-code/test/features/general/object_storage/general/encryption_at_rest  0.621s
```