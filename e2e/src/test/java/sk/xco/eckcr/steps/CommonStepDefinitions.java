package sk.xco.eckcr.steps;

import io.cucumber.java.en.And;
import io.cucumber.java.en.Given;
import io.cucumber.java.en.Then;
import io.cucumber.java.en.When;
import lombok.extern.slf4j.Slf4j;


import static org.junit.jupiter.api.Assertions.fail;
import static sk.xco.eckcr.util.K8sClient.withK8sClient;

@Slf4j
public class CommonStepDefinitions {

    private static final String ECK_CR_POD_NAME = "eck-custom-resources-operator";
    private static final String ES_POD_NAME_PATTERN = "%s-es";

    @Given("Kubernetes cluster is available")
    public void kubernetesClusterAvailable() {
        withK8sClient().run(client -> {
            client.pods().inNamespace("default").list();
            log.info("K8s available");
            return null;
        });
    }

    @And("ECK-CR operator is installed")
    public void eckCRInstalled() {
        withK8sClient().run(client -> {
            var pods = client.pods().inNamespace("default").list().getItems();
            if (pods.stream().noneMatch(pod -> pod.getMetadata().getName().contains(ECK_CR_POD_NAME))) {
                fail("ECK-CR not installed");
            }
            log.info("ECK-CR present");
            return null;
        });
    }

    @And("Elasticsearch {string} is available")
    public void elasticsearchAvailable(String esName) {
        withK8sClient().run(client -> {
            var pods = client.pods().inNamespace("default").list().getItems();
            var esPodName = ES_POD_NAME_PATTERN.formatted(esName);
            if (pods.stream().noneMatch(pod -> pod.getMetadata().getName().contains(esPodName))) {
                fail("ES %s not installed".formatted(esName));
            }
            log.info("ES {} present", esName);
            return null;
        });
    }

    @When("I create a Pod with name {string} and image {string}")
    public void iCreateAPodWithNameAndImage(String podName, String imageName) {
        // Use Kubernetes client to create a Pod here
    }

    @When("I create a Service with name {string} and type {string}")
    public void i_create_a_service_with_name_and_type(String string, String string2) {
        // Write code here that turns the phrase above into concrete actions
    }

    @Then("the Pod {string} should be running")
    public void thePodShouldBeRunning(String podName) {
        // Use Kubernetes client to check if the Pod is running here
        // Write assertions to verify the Pod's status
    }

    @Then("the Service {string} should be available")
    public void the_service_should_be_available(String string) {
        // Write code here that turns the phrase above into concrete actions
    }

}
