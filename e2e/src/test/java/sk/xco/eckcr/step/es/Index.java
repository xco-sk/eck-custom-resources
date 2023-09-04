package sk.xco.eckcr.step.es;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.fail;
import static sk.xco.eckcr.util.ESClient.getIndexState;

import co.elastic.clients.elasticsearch._types.ElasticsearchException;
import io.cucumber.java.en.Then;
import java.util.concurrent.TimeUnit;
import lombok.extern.slf4j.Slf4j;
import org.awaitility.Awaitility;

@Slf4j
public class Index {

  @Then(
      "the Index with name {string} with {int} replica shards is present in {string} Elasticsearch")
  public void indexWithNameAndReplicasPresent(
      String indexName, int replicas, String elasticsearchName) {
    Awaitility.await()
        .atMost(10, TimeUnit.SECONDS)
        .untilAsserted(
            () -> {
              try {
                var indexState = getIndexState(indexName);
                assertThat(indexState.settings()).isNotNull();
                assertThat(Integer.valueOf(indexState.settings().index().numberOfReplicas()))
                    .isEqualTo(replicas);
              } catch (ElasticsearchException e) {
                fail("Failed to get Index", e);
              }
            });
  }

  @Then("the Index with name {string} is not present in {string} Elasticsearch")
  public void indexNotPresent(String indexName, String elasticsearchName) {
    Awaitility.await()
        .atMost(10, TimeUnit.SECONDS)
        .untilAsserted(
            () -> {
              try {
                var indexState = getIndexState(indexName);
                fail("Index %s present in Elasticsearch".formatted(indexName));
              } catch (ElasticsearchException e) {
                assertThat(e.status()).isEqualTo(404);
              }
            });
  }

  public static void waitForIndex(String indexName) {
    Awaitility.await()
        .atMost(10, TimeUnit.SECONDS)
        .until(
            () -> {
              try {
                getIndexState(indexName);
                return true;
              } catch (ElasticsearchException e) {
                return false;
              }
            });
  }
}
