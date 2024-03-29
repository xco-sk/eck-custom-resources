package sk.xco.eckcr.step.es;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.fail;
import static sk.xco.eckcr.util.ESClient.getSnapshotRepo;

import co.elastic.clients.elasticsearch._types.ElasticsearchException;
import io.cucumber.java.en.Then;
import sk.xco.eckcr.util.Await;
import sk.xco.eckcr.util.ESClient;

public class SnapshotRepo {

  @Then(
      "the Snapshot Repository with name {string} is present in {string} Elasticsearch with {string} set to {string}")
  public void repoPresent(String repoName, String esName, String attrKey, String attrValue) {
    Await.untilAsserted(
        () -> {
          try {
            var repo = getSnapshotRepo(repoName);
            assertThat(repo).isNotNull();
            assertThat(repo.settings().compress()).isEqualTo(Boolean.valueOf(attrValue));
          } catch (ElasticsearchException e) {
            fail("Failed to get resource", e);
          }
        });
  }

  @Then("the Snapshot Repository with name {string} is not present in {string} Elasticsearch")
  public void repoNotPresent(String repoName, String esName) {
    ESClient.awaitResourceNotPresent(repoName, ESClient::getSnapshotRepo);
  }
}
