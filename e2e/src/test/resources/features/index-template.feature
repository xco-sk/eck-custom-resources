@Elasticsearch
Feature: Kubernetes Resource Creation

  Background:
    Given Kubernetes cluster is available
    And ECK-CR operator is installed
    And Elasticsearch "quickstart" is available

  Scenario: Create an Index Template
    When the "index-template" is applied with "replicas" set to "3"
    And the "index-template_index" is applied
    Then the Index Template with name "test" is present in "quickstart" Elasticsearch
    Then the Index with name "index-tpl-1" with 3 replica shards is present in "quickstart" Elasticsearch

  Scenario: Update an Index Template
    Given the Index Template "index-template" is present with "replicas" set to "0"
    When the "index-template" is applied with "replicas" set to "2"
    And the "index-template_index" is applied
    Then the Index with name "index-tpl-1" with 2 replica shards is present in "quickstart" Elasticsearch

  Scenario: Delete an Index Template
    Given the Index Template "index-template" is present with "replicas" set to "0"
    When the "index-template" is deleted
    Then the Index Template with name "test" is not present in "quickstart" Elasticsearch
