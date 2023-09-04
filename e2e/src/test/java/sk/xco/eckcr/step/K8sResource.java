package sk.xco.eckcr.step;

import static sk.xco.eckcr.util.K8sClient.withK8sClient;

import io.cucumber.java.After;
import io.cucumber.java.ParameterType;
import io.cucumber.java.en.Given;
import io.cucumber.java.en.When;
import io.fabric8.kubernetes.client.KubernetesClientException;
import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;
import lombok.extern.slf4j.Slf4j;
import org.yaml.snakeyaml.Yaml;
import sk.xco.eckcr.step.es.Index;
import sk.xco.eckcr.util.ApiType;

@Slf4j
public class K8sResource {
  private final Set<String> toCleanup = new HashSet<>();

  @When("the {string} is applied")
  public void applyResource(String resourceName) throws IOException {

    var resource =
        new String(
            K8sResource.class.getResourceAsStream(getResourcePath(resourceName)).readAllBytes());

    apply(resourceName, resource);
  }

  @When("the {string} is applied with {string} set to {string}")
  public void applyResourceWithReplacement(
      String resourceName, String replaceKey, String replaceValue) {
    String modified = getModifiedResource(resourceName, replaceKey, replaceValue);

    apply(resourceName, modified);
  }

  @Given("the {ApiType} {string} is present")
  public void givenResource(ApiType apiType, String resourceName) throws IOException {
    var resource =
        new String(
            K8sResource.class.getResourceAsStream(getResourcePath(resourceName)).readAllBytes());

    apply(resourceName, resource);

    waitForResource(apiType, resourceName);
  }

  @Given("the {ApiType} {string} is present with {string} set to {string}")
  public void givenResourceWithReplacement(
      ApiType apiType, String resourceName, String replaceKey, String replaceValue) {
    String modified = getModifiedResource(resourceName, replaceKey, replaceValue);

    apply(resourceName, modified);
  }

  private void waitForResource(ApiType apiType, String resourceName) {
    switch (apiType) {
      case Index -> Index.waitForIndex(resourceName);
      default -> throw new UnsupportedOperationException("Api type not supported");
    }
  }

  private void apply(String resourceName, String resource) {
    withK8sClient()
        .run(
            client -> {
              client
                  .load(new ByteArrayInputStream(resource.getBytes()))
                  .inNamespace("default")
                  .serverSideApply();
              toCleanup.add(resourceName);
              return null;
            });
  }

  private String getModifiedResource(String resourceName, String replaceKey, String replaceValue) {
    Yaml yaml = new Yaml();
    Map<String, Object> input =
        yaml.load(K8sResource.class.getResourceAsStream(getResourcePath(resourceName)));
    ((Map<String, Object>) input.get("spec"))
        .put(
            "body",
            ((String) ((Map<String, Object>) input.get("spec")).get("body"))
                .replace("$" + replaceKey, replaceValue));

    var modified = yaml.dump(input);
    return modified;
  }

  @Given("the {string} is deleted")
  public void deleteResource(String resourceName) {
    withK8sClient()
        .run(
            client -> {
              client
                  .load(K8sResource.class.getResourceAsStream(getResourcePath(resourceName)))
                  .inNamespace("default")
                  .delete();
              toCleanup.remove(resourceName);
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

  @ParameterType("Index")
  public ApiType ApiType(String stringifiedApiType) {
    return ApiType.valueOf(stringifiedApiType);
  }

  private String getResourcePath(String resourceName) {
    return "/resources/%s.yaml".formatted(resourceName);
  }
}
