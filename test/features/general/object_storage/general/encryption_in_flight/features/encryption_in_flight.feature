@intrusive_test
@service.object_storage
@service.object_storage.storage-encryption-in-flight
@CHC2-SVD001
@CHC2-AGP140
Feature: Object Storage Encryption in Flight

  As a Cloud Security Architect
  I want to ensure that suitable security controls are applied to Object Storage
  So that my organisation is protected against data leakage via misconfiguration

  Rule: CHC2-AGP140 - Ensure cryptographic controls are in place to protect the confidentiality and integrity of data in-transit, stored, generated and processed in the cloud

  @preventative
  @csp.azure
  Scenario Outline: Prevent Creation of Object Storage Without Encryption in Flight
    Given security controls that restrict data from being unencrypted in flight
    When we provision an Object Storage bucket
    And http access is "<HTTP Option>"
    And https access is "<HTTPS Option>"
    Then creation will "<Result>" with an error matching "<Error Description>"

    Examples:
      | HTTP Option | HTTPS Option | Result  | Error Description                                     |
      | enabled     | disabled     | Fail    | Storage Buckets must not be accessible via plain HTTP |
      | enabled     | enabled      | Fail    | Storage Buckets must not be accessible via plain HTTP |
      | disabled    | enabled      | Succeed |                                                       |

  @detective
  @csp.azure
  @csp.aws
  Scenario: Ensure Detective Checks for Object Storage Encryption in Flight are Enabled, When Supported
    Given the CSP provides a detective capability for unencrypted data transfer to Object Storage
    When we examine the detective measure
    Then the detective measure is enabled
