@non_intrusive_test
@service.object_storage
@service.object_storage.whitelisting
@CHC2-SVD030
Feature: Object Storage Has Network Whitelisting Measures Enforced

  As a Cloud Security Architect
  I want to ensure that suitable security controls are applied to Object Storage
  So that my organisation's data can only be accessed from whitelisted IP addresses

  Rule: CHC2-SVD030 - protect cloud service network access by limiting access from the appropriate source network only

  @detective
  @csp.aws
  @csp.azure
  Scenario: Check Object Storage is Configured With Network Source Address Whitelisting
    Given the CSP provides a whitelisting capability for Object Storage containers
    When we examine the Object Storage container in environment variable "TARGET_STORAGE_CONTAINER"
    Then whitelisting is configured with the given IP address range or an endpoint

  @preventative
  @csp.azure
  Scenario Outline: Prevent Object Storage from Being Created Without Network Source Address Whitelisting
    Given security controls that Prevent Object Storage from being created without network source address whitelisting are applied
    When we provision an Object Storage container
    And it is created with whitelisting entry "<Whitelist Entry>"
    Then creation will "<Result>"

    Examples:
      | Whitelist Entry | Result  |
      | 10.0.0.0        | Success |
      | 10.0.0.1        | Success |
      | nil             | Fail    |
