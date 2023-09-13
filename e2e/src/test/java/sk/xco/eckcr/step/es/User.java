package sk.xco.eckcr.step.es;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.fail;
import static sk.xco.eckcr.util.ESClient.getUser;

import co.elastic.clients.elasticsearch._types.ElasticsearchException;
import io.cucumber.java.en.Then;
import lombok.extern.slf4j.Slf4j;
import sk.xco.eckcr.util.Await;
import sk.xco.eckcr.util.ESClient;

@Slf4j
public class User {
  @Then(
      "the User with name {string} is present in {string} Elasticsearch with {string} set to {string}")
  public void userPresent(String userName, String esName, String attrKey, String attrValue) {
    Await.untilAsserted(
        () -> {
          try {
            var user = getUser(userName);
            assertThat(user).isNotNull();
            assertThat(user.fullName()).isEqualTo(attrValue);
          } catch (ElasticsearchException e) {
            fail("Failed to get resource", e);
          }
        });
  }

  @Then("the User with name {string} is not present in {string} Elasticsearch")
  public void userNotPresent(String userName, String esName) {
    ESClient.awaitResourceNotPresent(userName, ESClient::getUser);
  }
}
