@Elasticsearch
Feature: Snapshot repositories

  Background:
    Given Kubernetes cluster is available
    And ECK-CR operator is installed
    And Elasticsearch "quickstart" is available

  Scenario: Create a Snapshot repository
    When the "snapshot-repo.yaml" is applied with "compress" set to "false"
    Then the Snapshot Repository with name "test" is present in "quickstart" Elasticsearch with "compress" set to "false"

  Scenario: Update a Snapshot repository
    Given the Snapshot Repository "test" defined in "snapshot-repo.yaml" is present with "compress" set to "false"
    When the "snapshot-repo.yaml" is applied with "compress" set to "true"
    Then the Snapshot Repository with name "test" is present in "quickstart" Elasticsearch with "compress" set to "true"

  Scenario: Delete a Snapshot repository
    Given the Snapshot Repository "test" defined in "snapshot-repo.yaml" is present with "compress" set to "false"
    When the resource defined in "snapshot-repo.yaml" is deleted
    Then the Snapshot Repository with name "test" is not present in "quickstart" Elasticsearch
