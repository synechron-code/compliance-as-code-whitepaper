@internal
@json
Feature: Azure Policy JSON Documents are Legal JSON and Valid Against the Published Microsoft Schema

  Azure Policy documents that we commit to the repository must be both well-formed JSONs and also must comply with the appropriate schema.

  Scenario:

    Given a directory of Azure Policy files in JSON format
    Then the documents must be valid JSON
    And the JSON must be valid against the Microsoft schema
