package sk.xco.eckcr.step;

import static sk.xco.eckcr.util.K8sClient.withK8sClient;

import io.cucumber.java.After;
import io.cucumber.java.ParameterType;
import io.cucumber.java.en.Given;
import io.cucumber.java.en.When;
import io.fabric8.kubernetes.client.KubernetesClientException;
import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.time.Duration;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;
import lombok.extern.slf4j.Slf4j;
import org.yaml.snakeyaml.Yaml;
import sk.xco.eckcr.util.ApiType;
import sk.xco.eckcr.util.ESClient;

@Slf4j
public class K8sResource {
  private final Set<String> toCleanup = new HashSet<>();

  @When("the {string} is applied")
  public void applyResource(String fileName) throws IOException, InterruptedException {

    var resource =
        new String(K8sResource.class.getResourceAsStream(getResourcePath(fileName)).readAllBytes());

    apply(fileName, resource);
    Thread.sleep(Duration.ofMillis(500).toMillis());
  }

  @When("the {string} is applied with {string} set to {string}")
  public void applyResourceWithReplacement(String fileName, String replaceKey, String replaceValue)
      throws InterruptedException {
    String modified = getModifiedResource(fileName, replaceKey, replaceValue);

    apply(fileName, modified);
    Thread.sleep(Duration.ofMillis(500).toMillis());
  }

  @Given("the {ApiType} {string} defined in {string} is present")
  public void givenResource(ApiType apiType, String resourceName, String fileName)
      throws IOException {
    var resource =
        new String(K8sResource.class.getResourceAsStream(getResourcePath(fileName)).readAllBytes());

    apply(resourceName, resource);

    waitForResource(apiType, resourceName);
  }

  @Given("the {ApiType} {string} defined in {string} is present with {string} set to {string}")
  public void givenResourceWithReplacement(
      ApiType apiType,
      String resourceName,
      String fileName,
      String replaceKey,
      String replaceValue) {
    String modified = getModifiedResource(fileName, replaceKey, replaceValue);

    apply(fileName, modified);

    waitForResource(apiType, resourceName);
  }

  @Given("the resource defined in {string} is deleted")
  public void deleteResource(String fileName) {
    withK8sClient()
        .run(
            client -> {
              client
                  .load(K8sResource.class.getResourceAsStream(getResourcePath(fileName)))
                  .inNamespace("default")
                  .delete();
              toCleanup.remove(fileName);
              return null;
            });
  }

  @After
  public void afterScenario() {
    withK8sClient()
        .run(
            client -> {
              toCleanup.forEach(
                  r -> {
                    log.info("Cleaning-up {}", r);
                    try {
                      client
                          .load(K8sResource.class.getResourceAsStream(getResourcePath(r)))
                          .inNamespace("default")
                          .delete();
                    } catch (KubernetesClientException e) {
                      log.warn("Failed to cleanup {}", r, e);
                    }
                  });
              return null;
            });
    toCleanup.clear();
  }

  @ParameterType(
      "Index|Index Template|Index Lifecycle Policy|Ingest Pipeline|Snapshot Repository|Snapshot Lifecycle Policy|User")
  public ApiType ApiType(String stringifiedApiType) {
    return switch (stringifiedApiType) {
      case "Index Template" -> ApiType.IndexTemplate;
      case "Index Lifecycle Policy" -> ApiType.IndexLifecyclePolicy;
      case "Ingest Pipeline" -> ApiType.IngestPipeline;
      case "Snapshot Repository" -> ApiType.SnapshotRepo;
      case "Snapshot Lifecycle Policy" -> ApiType.SnapshotLifecyclePolicy;
      default -> ApiType.valueOf(stringifiedApiType);
    };
  }

  private void waitForResource(ApiType apiType, String resourceName) {
    switch (apiType) {
      case Index -> ESClient.waitForResource(resourceName, ESClient::getIndexState);
      case IndexTemplate -> ESClient.waitForResource(resourceName, ESClient::getTemplate);
      case IndexLifecyclePolicy -> ESClient.waitForResource(resourceName, ESClient::getIlmPolicy);
      case IngestPipeline -> ESClient.waitForResource(resourceName, ESClient::getIngestPipeline);
      case SnapshotRepo -> ESClient.waitForResource(resourceName, ESClient::getSnapshotRepo);
      case SnapshotLifecyclePolicy -> ESClient.waitForResource(
          resourceName, ESClient::getSnapshotLifecyclePolicy);
      case User -> ESClient.waitForResource(resourceName, ESClient::getUser);
      default -> throw new UnsupportedOperationException("Api type not supported");
    }
  }

  private void apply(String fileName, String resource) {
    withK8sClient()
        .run(
            client -> {
              client
                  .load(new ByteArrayInputStream(resource.getBytes()))
                  .inNamespace("default")
                  .serverSideApply();
              toCleanup.add(fileName);
              return null;
            });
  }

  private String getModifiedResource(String fileName, String replaceKey, String replaceValue) {
    Yaml yaml = new Yaml();
    Map<String, Object> input =
        yaml.load(K8sResource.class.getResourceAsStream(getResourcePath(fileName)));
    ((Map<String, Object>) input.get("spec"))
        .put(
            "body",
            ((String) ((Map<String, Object>) input.get("spec")).get("body"))
                .replace("$" + replaceKey, replaceValue));

    var modified = yaml.dump(input);
    return modified;
  }

  private String getResourcePath(String resourceName) {
    return "/resources/%s".formatted(resourceName);
  }
}
