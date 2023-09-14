@Elasticsearch
Feature: Ingest Pipelines

  Background:
    Given Kubernetes cluster is available
    And ECK-CR operator is installed
    And Elasticsearch "quickstart" is available

  Scenario: Create an Ingest Pipeline
    When the "ingest-pipeline.yaml" is applied with "value" set to "true"
    Then the Ingest Pipeline with name "test" is present in "quickstart" Elasticsearch

  Scenario: Update an Ingest Pipeline
    Given the Ingest Pipeline "test" defined in "ingest-pipeline.yaml" is present with "value" set to "true"
    When the "ingest-pipeline.yaml" is applied with "value" set to "false"
    Then the Ingest Pipeline with name "test" and value "false" is present in "quickstart" Elasticsearch

  Scenario: Delete an Ingest Pipeline
    Given the Ingest Pipeline "test" defined in "ingest-pipeline.yaml" is present with "value" set to "true"
    When the resource defined in "ingest-pipeline.yaml" is deleted
    Then the Ingest Pipeline with name "test" is not present in "quickstart" Elasticsearch
