@intrusive_test
@csp.aws
@csp.azure
@service.object_storage
@CHC2-SVD030
Feature: Object Storage Network Access Control
  Ensure that 'Public access level' is set to Private for blob containers

  As a Cloud Security Architect
  I want to ensure that suitable security controls are applied to Object Storage
  So that my organisation is protected against data leakage via misconfiguration

  Scenario Outline: Only Allow the Creation of Compliant Storage Buckets
    Given security controls that restrict data from being publicly accessible
    When we provision an Object Storage bucket
    And the bucket has the IP Whitelisting option set to "<IP Whitelisting Option>"
    And the bucket contains the whitelisted IPs "<Whitelisted IPs>"
    And the bucket has the public option set to "<Public Option>"
    Then creation will "<Result>" with an error matching "<Error Description>"

    Examples:
      | IP Whitelisting Option | Whitelisted IPs | Public Option | Result  | Error Description                                                 |
      | Enabled                | Non-Firm IPs    | Enabled       | Fail    | Storage Buckets must not be public                                |
      | Disabled               | n/a             | Disabled      | Fail    | Storage Buckets without IP whitelisting are disallowed            |
      | Disabled               | n/a             | Enabled       | Fail    | Storage Buckets with IP whitelisting must use a Firm's IP address |
      | Enabled                | Firm IPs Only   | Enabled       | Fail    | Storage Buckets must not be public                                |
      | Enabled                | Firm IPs Only   | Disabled      | Succeed |                                                                   |
