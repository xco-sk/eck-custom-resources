@Elasticsearch
Feature: API Keys

  Background:
    Given Kubernetes cluster is available
    And ECK-CR operator is installed
    And Elasticsearch "quickstart" is available

  Scenario: Create an API Key
    Given the "suffix" is set to "random number"
    When the "apikey.yaml" is applied with "suffix"
    Then the API Key with name "test" suffixed with "suffix" is in state "valid" in "quickstart" Elasticsearch

  Scenario: Delete an API Key
    Given the "suffix" is set to "random number"
    Given the API Key "test" suffixed with "suffix" defined in "apikey.yaml" is present with "suffix"
    When the resource defined in "apikey.yaml" is deleted
    Then the API Key with name "test" suffixed with "suffix" is in state "invalidated" in "quickstart" Elasticsearch
