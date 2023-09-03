Feature: Kubernetes Resource Creation

  Background:
    Given Kubernetes cluster is available
    And ECK-CR operator is installed

  Scenario: Create a Pod
    When I create a Pod with name "my-pod" and image "nginx"
    Then the Pod "my-pod" should be running

  Scenario: Create a Service
    When I create a Service with name "my-service" and type "ClusterIP"
    Then the Service "my-service" should be available
