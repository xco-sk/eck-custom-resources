package sk.xco.eckcr.step.es;

import static java.util.Objects.nonNull;
import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.fail;
import static sk.xco.eckcr.util.ESClient.getIngestPipeline;

import co.elastic.clients.elasticsearch._types.ElasticsearchException;
import io.cucumber.java.en.Then;
import java.util.concurrent.TimeUnit;
import org.awaitility.Awaitility;

public class IngestPipeline {

  @Then("the Ingest Pipeline with name {string} is present in {string} Elasticsearch")
  public void pipelinePresent(String resourceName, String esName) {
    Awaitility.await()
        .atMost(10, TimeUnit.SECONDS)
        .untilAsserted(
            () -> {
              try {
                var pipeline = getIngestPipeline(resourceName);
                assertThat(pipeline).isNotNull();
              } catch (ElasticsearchException e) {
                fail("Failed to get resource", e);
              }
            });
  }

  @Then(
      "the Ingest Pipeline with name {string} and value {string} is present in {string} Elasticsearch")
  public void pipelinePresentWithValue(String resourceName, String value, String esName) {
    Awaitility.await()
        .atMost(10, TimeUnit.SECONDS)
        .untilAsserted(
            () -> {
              try {
                var pipeline = getIngestPipeline(resourceName);
                assertThat(pipeline).isNotNull();
                assertThat(pipeline.processors().get(0).set().value().toJson().toString())
                    .isEqualTo(value);
              } catch (ElasticsearchException e) {
                fail("Failed to get resource", e);
              }
            });
  }

  @Then("the Ingest Pipeline with name {string} is not present in {string} Elasticsearch")
  public void pipelineNotPresent(String resourceName, String esName) {
    Awaitility.await()
        .atMost(5, TimeUnit.SECONDS)
        .untilAsserted(
            () -> {
              try {
                var p = getIngestPipeline(resourceName);
                if (nonNull(p)) {
                  fail(
                      "Ingest pipeline %s present in Elasticsearch: %s".formatted(resourceName, p));
                }
              } catch (ElasticsearchException e) {
                assertThat(e.status()).isEqualTo(404);
              }
            });
  }

  public static void waitForIngestPipeline(String resourceName) {
    Awaitility.await()
        .atMost(10, TimeUnit.SECONDS)
        .untilAsserted(
            () -> {
              try {
                var template = getIngestPipeline(resourceName);
                assertThat(template).isNotNull();
              } catch (ElasticsearchException e) {
                fail("Failed to get resource", e);
              }
            });
  }
}
