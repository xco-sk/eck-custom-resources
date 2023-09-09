@Elasticsearch
Feature: Index lifecycle policy

  Background:
    Given Kubernetes cluster is available
    And ECK-CR operator is installed
    And Elasticsearch "quickstart" is available

  Scenario: Create an Index Lifecycle Policy
    When the "index-lifecycle-policy.yaml" is applied with "minAge" set to "15d"
    Then the Index Lifecycle Policy with name "test" is present in "quickstart" Elasticsearch

  Scenario: Update an Index Lifecycle Policy
    Given the Index Lifecycle Policy "test" defined in "index-lifecycle-policy.yaml" is present with "minAge" set to "15d"
    When the "index-lifecycle-policy.yaml" is applied with "minAge" set to "30d"
    Then the Index Lifecycle Policy with name "test" is present in "quickstart" Elasticsearch
    And the Index Lifecycle Policy with name "test" got delete min age set to "30d"

  Scenario: Delete an Index Lifecycle Policy
    Given the Index Lifecycle Policy "test" defined in "index-lifecycle-policy.yaml" is present with "minAge" set to "15d"
    When the resource defined in "index-lifecycle-policy.yaml" is deleted
    Then the Index Lifecycle Policy with name "test" is not present in "quickstart" Elasticsearch
