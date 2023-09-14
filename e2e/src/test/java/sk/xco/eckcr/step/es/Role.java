package sk.xco.eckcr.step.es;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.fail;
import static sk.xco.eckcr.util.ESClient.getRole;

import co.elastic.clients.elasticsearch._types.ElasticsearchException;
import io.cucumber.java.en.Then;
import sk.xco.eckcr.util.Await;
import sk.xco.eckcr.util.ESClient;

public class Role {
  @Then(
      "the Role with name {string} is present in {string} Elasticsearch with {string} set to {string}")
  public void rolePresent(String roleName, String esName, String attrKey, String attrValue) {
    Await.untilAsserted(
        () -> {
          try {
            var role = getRole(roleName);
            assertThat(role).isNotNull();
            assertThat(role.indices().get(0).names()).contains(attrValue);
          } catch (ElasticsearchException e) {
            fail("Failed to get resource", e);
          }
        });
  }

  @Then("the Role with name {string} is not present in {string} Elasticsearch")
  public void roleNotPresent(String roleName, String esName) {
    ESClient.awaitResourceNotPresent(roleName, ESClient::getRole);
  }
}
