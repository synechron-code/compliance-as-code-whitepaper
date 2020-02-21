@intrusive_test
@service.object_storage
@service.object_storage.customer_managed_keys
@CHC2-SVD001
@CHC2-AGP135
Feature: Customer Managed Keys for Object Storage

  As a Cloud Security Architect
  I want to ensure that customer managed keys are used for Object Storage
  So that my organisation is protected against data leakage and third party subpoenas

  Rule: CHC2-AGP135 - Ensure encryption keys are owned and managed by the FI following industry best practice for key management

  @preventative
  @csp.azure
  Scenario Outline: Prevent Creation of Object Storage Without Customer Managed Keys
    Given security controls that restrict data from being encrypted with non-customer managed keys
    When we provision an Object Storage account
    And the encryption key is managed by "<Encryption Key Owner>"
    Then creation will "<Result>" with an error matching "<Error Description>"

    Examples:
      | Encryption Key Owner | Result  | Error Description                                             |
      | non-customer         | Fail    | Storage accounts must not be encrypted with non-customer keys |
      | customer             | Succeed |                                                               |

  @detective
  @csp.azure
  @csp.aws
  Scenario: Ensure Detective Checks for Customer Managed Encryption Keys for Object Storage are Enabled, When Supported
    Given the CSP provides a detective capability for non-customer managed encryption keys for Object Storage
    When we examine the detective measure
    Then the detective measure is enabled
