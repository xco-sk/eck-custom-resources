package sk.xco.eckcr.step.es;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.fail;
import static sk.xco.eckcr.util.ESClient.getTemplate;

import co.elastic.clients.elasticsearch._types.ElasticsearchException;
import io.cucumber.java.en.Then;
import java.util.concurrent.TimeUnit;
import lombok.extern.slf4j.Slf4j;
import org.awaitility.Awaitility;
import sk.xco.eckcr.util.ESClient;

@Slf4j
public class IndexTemplate {

  @Then("the Index Template with name {string} is present in {string} Elasticsearch")
  public void indexTemplateWithNamePresent(String templateName, String elasticsearchName) {
    Awaitility.await()
        .atMost(10, TimeUnit.SECONDS)
        .untilAsserted(
            () -> {
              try {
                var template = getTemplate(templateName);
                assertThat(template).isNotNull();
              } catch (ElasticsearchException e) {
                fail("Failed to get Index Template", e);
              }
            });
  }

  @Then("the Index Template with name {string} is not present in {string} Elasticsearch")
  public void indexTemplateNotPresent(String templateName, String elasticsearchName) {
    ESClient.awaitResourceNotPresent(templateName, ESClient::getTemplate);
  }
}
