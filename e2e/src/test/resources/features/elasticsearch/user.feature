@Elasticsearch
Feature: Users

  Background:
    Given Kubernetes cluster is available
    And ECK-CR operator is installed
    And Elasticsearch "quickstart" is available
    And the "user-secret.yaml" is applied

  Scenario: Create an User
    When the "user.yaml" is applied with "fullName" set to "John Doe"
    Then the User with name "test" is present in "quickstart" Elasticsearch with "fullName" set to "John Doe"

  Scenario: Update an User
    Given the User "test" defined in "user.yaml" is present with "fullName" set to "John Doe"
    When the "user.yaml" is applied with "fullName" set to "Jane Doe"
    Then the User with name "test" is present in "quickstart" Elasticsearch with "fullName" set to "Jane Doe"

  Scenario: Delete an User
    Given the User "test" defined in "user.yaml" is present with "fullName" set to "John Doe"
    When the resource defined in "user.yaml" is deleted
    Then the User with name "test" is not present in "quickstart" Elasticsearch
