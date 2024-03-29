@Elasticsearch
Feature: Indices

  Background:
    Given Kubernetes cluster is available
    And ECK-CR operator is installed
    And Elasticsearch "quickstart" is available

  Scenario: Create an Index
    When the "index.yaml" is applied with "replicas" set to "0"
    Then the Index with name "test" with 0 replica shards is present in "quickstart" Elasticsearch

  Scenario: Update an Index
    Given the Index "test" defined in "index.yaml" is present with "replicas" set to "0"
    When the "index.yaml" is applied with "replicas" set to "1"
    Then the Index with name "test" with 1 replica shards is present in "quickstart" Elasticsearch

  Scenario: Delete an Index
    Given the Index "test" defined in "index.yaml" is present with "replicas" set to "0"
    When the resource defined in "index.yaml" is deleted
    Then the Index with name "test" is not present in "quickstart" Elasticsearch
