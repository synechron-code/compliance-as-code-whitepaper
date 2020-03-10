@intrusive_test
@service.object_storage
@encryption_at_rest
@CCO:CHC2-SVD001
@CCO:CHC2-AGP140
@CCO:CHC2-EUC001
@csp.aws
@csp.azure
Feature: Object Storage Encryption at Rest

  As a Cloud Security Architect
  I want to ensure that suitable security controls are applied to Object Storage
  So that my organisation is protected against data leakage due to misconfiguration

  Rule: CHC2-AGP140 - Ensure cryptographic controls are in place to protect the confidentiality and integrity of data in-transit, stored, generated and processed in the cloud

    @detective
    Scenario: Ensure Detective Checks for Object Storage Encryption at Rest are Enabled, When Supported
      Given the CSP provides a detective capability for unencrypted Object Storage containers
      When we examine the detective measure
      Then the detective measure is enabled

    @preventative
    Scenario Outline: Prevent creation of Object Storage Without Encryption at Rest
      Given security controls that enforce data at rest encryption for Object Storage are applied
      When we provision an Object Storage container
      And it is created with encryption option "<Encryption Option>"
      Then creation will "<Result>"

      Examples:
        | Encryption Option | Result  |
        | true              | Success |
        | false             | Fail    |
