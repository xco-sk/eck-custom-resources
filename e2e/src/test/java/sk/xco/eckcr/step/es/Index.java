package sk.xco.eckcr.step.es;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.fail;
import static sk.xco.eckcr.util.ESClient.getIndexState;

import co.elastic.clients.elasticsearch._types.ElasticsearchException;
import io.cucumber.java.en.Then;
import lombok.extern.slf4j.Slf4j;
import sk.xco.eckcr.util.Await;
import sk.xco.eckcr.util.ESClient;

@Slf4j
public class Index {

  @Then(
      "the Index with name {string} with {int} replica shards is present in {string} Elasticsearch")
  public void indexWithNameAndReplicasPresent(
      String indexName, int replicas, String elasticsearchName) {
    Await.untilAsserted(
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
    ESClient.awaitResourceNotPresent(indexName, ESClient::getIndexState);
  }
}
