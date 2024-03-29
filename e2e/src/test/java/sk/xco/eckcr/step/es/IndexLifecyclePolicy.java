package sk.xco.eckcr.step.es;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.fail;
import static sk.xco.eckcr.util.ESClient.getIlmPolicy;

import co.elastic.clients.elasticsearch._types.ElasticsearchException;
import io.cucumber.java.en.And;
import io.cucumber.java.en.Then;
import java.util.concurrent.TimeUnit;
import org.awaitility.Awaitility;
import sk.xco.eckcr.util.ESClient;

public class IndexLifecyclePolicy {
    
  @Then("the Index Lifecycle Policy with name {string} is present in {string} Elasticsearch")
  public void ilmPresent(String policyName, String esName) {
    Awaitility.await()
        .atMost(10, TimeUnit.SECONDS)
        .untilAsserted(
            () -> {
              try {
                var ilmPolicy = getIlmPolicy(policyName);
                assertThat(ilmPolicy).isNotNull();
              } catch (ElasticsearchException e) {
                fail("Failed to get resource", e);
              }
            });
  }

  @And("the Index Lifecycle Policy with name {string} got delete min age set to {string}")
  public void ilmGotDeleteSet(String policyName, String minAge) {
    Awaitility.await()
        .atMost(10, TimeUnit.SECONDS)
        .untilAsserted(
            () -> {
              try {
                var ilmPolicy = getIlmPolicy(policyName);
                assertThat(ilmPolicy.phases().delete().minAge().time()).isEqualTo(minAge);
              } catch (ElasticsearchException e) {
                fail("Failed to get resource", e);
              }
            });
  }

  @Then("the Index Lifecycle Policy with name {string} is not present in {string} Elasticsearch")
  public void ilmNotPresent(String policyName, String esName) {
    ESClient.awaitResourceNotPresent(policyName, ESClient::getIlmPolicy);
  }
}
