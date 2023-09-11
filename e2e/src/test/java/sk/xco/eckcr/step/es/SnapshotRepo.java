package sk.xco.eckcr.step.es;

import static java.util.Objects.nonNull;
import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.fail;
import static sk.xco.eckcr.util.ESClient.getSnapshotRepo;

import co.elastic.clients.elasticsearch._types.ElasticsearchException;
import io.cucumber.java.en.Then;
import java.util.concurrent.TimeUnit;
import org.awaitility.Awaitility;

public class SnapshotRepo {
  public static void waitForSnapshotRepo(String repoName) {
    Awaitility.await()
        .atMost(10, TimeUnit.SECONDS)
        .until(
            () -> {
              try {
                getSnapshotRepo(repoName);
                return true;
              } catch (ElasticsearchException e) {
                return false;
              }
            });
  }

  @Then(
      "the Snapshot Repository with name {string} is present in {string} Elasticsearch with {string} set to {string}")
  public void repoPresent(String repoName, String esName, String attrKey, String attrValue) {
    Awaitility.await()
        .atMost(10, TimeUnit.SECONDS)
        .untilAsserted(
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
    Awaitility.await()
        .atMost(5, TimeUnit.SECONDS)
        .untilAsserted(
            () -> {
              try {
                var repo = getSnapshotRepo(repoName);
                if (nonNull(repo)) {
                  fail("Snapshot repo %s present in Elasticsearch: %s".formatted(repoName, repo));
                }
              } catch (ElasticsearchException e) {
                assertThat(e.status()).isEqualTo(404);
              }
            });
  }
}
