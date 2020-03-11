# Restrict access to a known set of IP addresses

## AWS

### Implementation Details

Our AWS example demonstrates how we can use simple configuration inspection to attest whether the Bucket Policy of a specific S3 Bucket enforces the expected network access controls.

We test for the presence of
* Specific IP addresses
* VPC
* VPCe

### Example Run

``` 
>go test
2020/03/11 12:25:59 [DEBUG] Setting up 'accessWhitelistingAWS'
Feature: Object Storage Has Network Whitelisting Measures Enforced
  As a Cloud Security Architect
  I want to ensure that suitable security controls are applied to Object Storage
  So that my organisation's data can only be accessed from whitelisted IP addresses

  Rule: CHC2-SVD030 - protect cloud service network access by limiting access from the appropriate source network only

  Scenario: Check Object Storage is Configured With Network Source Address Whitelisting             # features\access_whitelisting.feature:16
    Given the CSP provides a whitelisting capability for Object Storage containers                  # access_whitelisting_test.go:18 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/access_whitelisting.accessWhitelisting.cspSupportsWhitelisting-fm
    When we examine the Object Storage container in environment variable "TARGET_STORAGE_CONTAINER" # access_whitelisting_test.go:19 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/access_whitelisting.accessWhitelisting.examineStorageContainer-fm
2020/03/11 12:26:00 [DEBUG] policy: {"Version":"2012-10-17","Id":"VPCe and SourceIP","Statement":[{"Sid":"allowmw","Effect":"Allow","Principal":{"AWS":"arn:aws:iam::857701739302:user/mark.wong"},"Action":"s3:*","Resource":"arn:aws:s3:::bddipaddronlytest"},{"Sid":"VPCe and SourceIP","Effect":"Deny","Principal":"*","Action":"s3:*Object*","Resource":["arn:aws:s3:::bddipaddronlytest","arn:aws:s3:::bddipaddronlytest/*"],"Condition":{"StringNotLike":{"aws:sourceVpce":["vpce-1111111","vpce-2222222"]},"NotIpAddress":{"aws:SourceIp":["11.11.11.11/32","22.22.22.22/32"]}}},{"Effect":"Deny","Principal":"*","Action":"*","Resource":"arn:aws:s3:::bddipaddronlytest/*","Condition":{"Bool":{"aws:SecureTransport":"false"}}}]}
2020/03/11 12:26:00 [DEBUG] NotIpAddress: map[aws:SourceIp:[11.11.11.11/32 22.22.22.22/32]]
    Then whitelisting is configured with the given IP address range or an endpoint                  # access_whitelisting_test.go:20 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/access_whitelisting.accessWhitelisting.whitelistingIsConfigured-fm

  Scenario Outline: Prevent Object Storage from Being Created Without Network Source Address Whitelisting                          # features\access_whitelisting.feature:22
    Given security controls that Prevent Object Storage from being created without network source address whitelisting are applied # access_whitelisting_test.go:21 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/access_whitelisting.accessWhitelisting.checkPolicyAssigned-fm
    When we provision an Object Storage container                                                                                  # access_whitelisting_test.go:22 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/access_whitelisting.accessWhitelisting.provisionStorageContainer-fm
    And it is created with whitelisting entry "<Whitelist Entry>"                                                                  # access_whitelisting_test.go:23 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/access_whitelisting.accessWhitelisting.createWithWhitelist-fm
    Then creation will "<Result>"                                                                                                  # access_whitelisting_test.go:24 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/access_whitelisting.accessWhitelisting.creationWill-fm

    Examples:
      | Whitelist Entry | Result  |
      | 219.79.19.0/24  | Success |
        AWS does not support preventative controls for access whitelisting on S3
      | 219.79.19.1     | Fail    |
        AWS does not support preventative controls for access whitelisting on S3
      | 219.108.32.1    | Fail    |
        AWS does not support preventative controls for access whitelisting on S3
      | 170.74.231.168  | Success |
        AWS does not support preventative controls for access whitelisting on S3
      | nil             | Fail    |
        AWS does not support preventative controls for access whitelisting on S3
2020/03/11 12:26:00 [DEBUG] Teardown completed

--- Failed steps:

  Scenario Outline: Prevent Object Storage from Being Created Without Network Source Address Whitelisting # features\access_whitelisting.feature:22
    Given security controls that Prevent Object Storage from being created without network source address whitelisting are applied # features\access_whitelisting.feature:23
      Error: AWS does not support preventative controls for access whitelisting on S3

  Scenario Outline: Prevent Object Storage from Being Created Without Network Source Address Whitelisting # features\access_whitelisting.feature:22
    Given security controls that Prevent Object Storage from being created without network source address whitelisting are applied # features\access_whitelisting.feature:23
      Error: AWS does not support preventative controls for access whitelisting on S3

  Scenario Outline: Prevent Object Storage from Being Created Without Network Source Address Whitelisting # features\access_whitelisting.feature:22
    Given security controls that Prevent Object Storage from being created without network source address whitelisting are applied # features\access_whitelisting.feature:23
      Error: AWS does not support preventative controls for access whitelisting on S3

  Scenario Outline: Prevent Object Storage from Being Created Without Network Source Address Whitelisting # features\access_whitelisting.feature:22
    Given security controls that Prevent Object Storage from being created without network source address whitelisting are applied # features\access_whitelisting.feature:23
      Error: AWS does not support preventative controls for access whitelisting on S3

  Scenario Outline: Prevent Object Storage from Being Created Without Network Source Address Whitelisting # features\access_whitelisting.feature:22
    Given security controls that Prevent Object Storage from being created without network source address whitelisting are applied # features\access_whitelisting.feature:23
      Error: AWS does not support preventative controls for access whitelisting on S3


6 scenarios (1 passed, 5 failed)
23 steps (3 passed, 5 failed, 15 skipped)
578.7679ms
testing: warning: no tests to run
PASS
exit status 1
FAIL    citihub.com/compliance-as-code/test/features/general/object_storage/general/access_whitelisting 1.384s
```

## Azure

### Implementation Details

In our Azure Policy example, we attest that Azure Policy is in place by attempting to create Storage Accounts that do not have specific IP Whitelisting rules in place, as an in-band evaluation.  

The exact values will be organisational-specific, so the execution can be modified to your needs simply by modifying the values in the Scenario Outline table.

### Example Run
```
>go test
2020/03/11 12:23:47 [DEBUG] Setting up 'accessWhitelistingAzure'
2020/03/11 12:23:47 [DEBUG] creating Resource Group 'testyrnawaresourceGP' on location: 'eastasia'
2020/03/11 12:23:48 [DEBUG] Created Resource Group: testyrnawaresourceGP
Feature: Object Storage Has Network Whitelisting Measures Enforced
  As a Cloud Security Architect
  I want to ensure that suitable security controls are applied to Object Storage
  So that my organisation's data can only be accessed from whitelisted IP addresses

  Rule: CHC2-SVD030 - protect cloud service network access by limiting access from the appropriate source network only

  Scenario: Check Object Storage is Configured With Network Source Address Whitelisting             # features\access_whitelisting.feature:16
    Given the CSP provides a whitelisting capability for Object Storage containers                  # access_whitelisting_test.go:18 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/access_whitelisting.accessWhitelisting.cspSupportsWhitelisting-fm
2020/03/11 12:23:49 [DEBUG] IP WhiteListing: 219.73.58.34, Allow
2020/03/11 12:23:49 [DEBUG] Whitelisting rule exists. [Step PASSED]
    When we examine the Object Storage container in environment variable "TARGET_STORAGE_CONTAINER" # access_whitelisting_test.go:19 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/access_whitelisting.accessWhitelisting.examineStorageContainer-fm
    Then whitelisting is configured with the given IP address range or an endpoint                  # access_whitelisting_test.go:20 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/access_whitelisting.accessWhitelisting.whitelistingIsConfigured-fm
2020/03/11 12:23:49 [DEBUG] Getting Policy Assignment with scope: /providers/Microsoft.Management/managementGroups/boxbank-root
2020/03/11 12:23:50 [DEBUG] Policy Assignment check: deny_storage_wo_net_acl [Step PASSED]

  Scenario Outline: Prevent Object Storage from Being Created Without Network Source Address Whitelisting                          # features\access_whitelisting.feature:22
    Given security controls that Prevent Object Storage from being created without network source address whitelisting are applied # access_whitelisting_test.go:21 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/access_whitelisting.accessWhitelisting.checkPolicyAssigned-fm
    When we provision an Object Storage container                                                                                  # access_whitelisting_test.go:22 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/access_whitelisting.accessWhitelisting.provisionStorageContainer-fm
    And it is created with whitelisting entry "<Whitelist Entry>"                                                                  # access_whitelisting_test.go:23 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/access_whitelisting.accessWhitelisting.createWithWhitelist-fm
    Then creation will "<Result>"                                                                                                  # access_whitelisting_test.go:24 -> citihub.com/compliance-as-code/test/features/general/object_storage/general/access_whitelisting.accessWhitelisting.creationWill-fm

    Examples:
      | Whitelist Entry | Result  |
      | 219.79.19.0/24  | Success |
2020/03/11 12:24:09 [DEBUG] Getting Policy Assignment with scope: /providers/Microsoft.Management/managementGroups/boxbank-root
2020/03/11 12:24:10 [DEBUG] Policy Assignment check: deny_storage_wo_net_acl [Step PASSED]
      | 219.79.19.1     | Fail    |
2020/03/11 12:24:11 [DEBUG] Getting Policy Assignment with scope: /providers/Microsoft.Management/managementGroups/boxbank-root
2020/03/11 12:24:12 [DEBUG] Policy Assignment check: deny_storage_wo_net_acl [Step PASSED]
      | 219.108.32.1    | Fail    |
2020/03/11 12:24:13 [DEBUG] Getting Policy Assignment with scope: /providers/Microsoft.Management/managementGroups/boxbank-root
2020/03/11 12:24:13 [DEBUG] Policy Assignment check: deny_storage_wo_net_acl [Step PASSED]
      | 170.74.231.168  | Success |
2020/03/11 12:24:33 [DEBUG] Getting Policy Assignment with scope: /providers/Microsoft.Management/managementGroups/boxbank-root
2020/03/11 12:24:34 [DEBUG] Policy Assignment check: deny_storage_wo_net_acl [Step PASSED]
      | nil             | Fail    |
2020/03/11 12:24:35 [DEBUG] Deleting resources
2020/03/11 12:24:35 [DEBUG] Teardown completed

6 scenarios (6 passed)
23 steps (23 passed)
48.3228598s
testing: warning: no tests to run
PASS
ok      citihub.com/compliance-as-code/test/features/general/object_storage/general/access_whitelisting 49.009s
```