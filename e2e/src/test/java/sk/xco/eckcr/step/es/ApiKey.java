package sk.xco.eckcr.step.es;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.fail;
import static sk.xco.eckcr.util.ESClient.getApiKey;

import co.elastic.clients.elasticsearch._types.ElasticsearchException;
import io.cucumber.java.en.Then;
import sk.xco.eckcr.step.Common;
import sk.xco.eckcr.util.Await;

public class ApiKey {

  @Then(
      "the API Key with name {string} suffixed with {string} is in state {string} in {string} Elasticsearch")
  public void apiKeyPresent(String apiKeyName, String variableName, String state, String esName) {
    Await.untilAsserted(
        () -> {
          try {
            var resourceName = apiKeyName + Common.VARIABLES.get(variableName);

            var apiKey = getApiKey(resourceName);
            assertThat(apiKey).isNotNull();
            assertThat(apiKey.invalidated()).isEqualTo("invalidated".equals(state));
          } catch (ElasticsearchException e) {
            fail("Failed to get resource", e);
          }
        });
  }
}
