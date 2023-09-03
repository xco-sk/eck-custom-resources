Feature: Kubernetes Resource Creation

  Background:
    Given Kubernetes cluster is available
    And ECK-CR operator is installed
    And Elasticsearch "quickstart" is available

  Scenario: Create an Index
    When the "index-test" is applied
    Then the Index with name "test" with "0" replica shards is present in "quickstart" Elasticsearch

  Scenario: Update an Index
    When the "index-modified" is applied
    Then the Index with name "test" with "1" replica shards is present in "quickstart" Elasticsearch

  Scenario: Delete an Index
    When the "index-test" is applied
    When the resource of type "index" with name "test" is deleted
    Then the Index with name "test" is not present in "quickstart" Elasticsearch
