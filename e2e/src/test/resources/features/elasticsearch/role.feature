@Elasticsearch
Feature: User Roles

  Background:
    Given Kubernetes cluster is available
    And ECK-CR operator is installed
    And Elasticsearch "quickstart" is available

  Scenario: Create a Role
    When the "role.yaml" is applied with "indexName" set to "index1"
    Then the Role with name "test" is present in "quickstart" Elasticsearch with "indexName" set to "index1"

  Scenario: Update a Role
    Given the Role "test" defined in "role.yaml" is present with "indexName" set to "index1"
    When the "role.yaml" is applied with "indexName" set to "index2"
    Then the Role with name "test" is present in "quickstart" Elasticsearch with "indexName" set to "index2"

  Scenario: Delete a Role
    Given the Role "test" defined in "role.yaml" is present with "indexName" set to "index1"
    When the resource defined in "role.yaml" is deleted
    Then the Role with name "test" is not present in "quickstart" Elasticsearch
