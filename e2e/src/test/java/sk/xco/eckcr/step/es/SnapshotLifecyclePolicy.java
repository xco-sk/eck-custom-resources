package sk.xco.eckcr.step.es;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.fail;
import static sk.xco.eckcr.util.ESClient.getSnapshotLifecyclePolicy;

import co.elastic.clients.elasticsearch._types.ElasticsearchException;
import io.cucumber.java.en.Then;
import java.util.concurrent.TimeUnit;
import org.awaitility.Awaitility;
import sk.xco.eckcr.util.ESClient;

public class SnapshotLifecyclePolicy {
  @Then(
      "the Snapshot Lifecycle Policy with name {string} is present in {string} Elasticsearch with {string} set to {string}")
  public void policyPresent(String policyName, String esName, String attrKey, String attrValue) {
    Awaitility.await()
        .atMost(10, TimeUnit.SECONDS)
        .untilAsserted(
            () -> {
              try {
                var policy = getSnapshotLifecyclePolicy(policyName);
                assertThat(policy).isNotNull();
                assertThat(policy.policy().config().ignoreUnavailable())
                    .isEqualTo(Boolean.valueOf(attrValue));
              } catch (ElasticsearchException e) {
                fail("Failed to get resource", e);
              }
            });
  }

  @Then("the Snapshot Lifecycle Policy with name {string} is not present in {string} Elasticsearch")
  public void policyNotPresent(String policyName, String esName) {
    ESClient.awaitResourceNotPresent(policyName, ESClient::getSnapshotLifecyclePolicy);
  }
}
