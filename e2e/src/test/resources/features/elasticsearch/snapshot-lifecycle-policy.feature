@Elasticsearch
Feature: Snapshot Lifecycle Policies

  Background:
    Given Kubernetes cluster is available
    And ECK-CR operator is installed
    And Elasticsearch "quickstart" is available
    And the Snapshot Repository "test" defined in "snapshot-repo.yaml" is present with "compress" set to "false"

  Scenario: Create a Snapshot Lifecycle Policy
    When the "snapshot-lifecycle-policy.yaml" is applied with "ignoreUnavailable" set to "false"
    Then the Snapshot Lifecycle Policy with name "test" is present in "quickstart" Elasticsearch with "ignoreUnavailable" set to "false"

  Scenario: Update a Snapshot Lifecycle Policy
    Given the Snapshot Lifecycle Policy "test" defined in "snapshot-lifecycle-policy.yaml" is present with "ignoreUnavailable" set to "false"
    When the "snapshot-lifecycle-policy.yaml" is applied with "ignoreUnavailable" set to "true"
    Then the Snapshot Lifecycle Policy with name "test" is present in "quickstart" Elasticsearch with "ignoreUnavailable" set to "true"

  Scenario: Delete a Snapshot Lifecycle Policy
    Given the Snapshot Lifecycle Policy "test" defined in "snapshot-lifecycle-policy.yaml" is present with "ignoreUnavailable" set to "false"
    When the resource defined in "snapshot-lifecycle-policy.yaml" is deleted
    Then the Snapshot Lifecycle Policy with name "test" is not present in "quickstart" Elasticsearch
