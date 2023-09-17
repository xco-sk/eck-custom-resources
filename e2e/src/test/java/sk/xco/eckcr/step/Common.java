package sk.xco.eckcr.step;

import static org.junit.jupiter.api.Assertions.fail;
import static sk.xco.eckcr.util.K8sClient.withK8sClient;

import io.cucumber.java.After;
import io.cucumber.java.en.And;
import io.cucumber.java.en.Given;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ThreadLocalRandom;
import lombok.extern.slf4j.Slf4j;

@Slf4j
public class Common {

  private static final String ECK_CR_POD_NAME = "eck-custom-resources-operator";
  private static final String ES_POD_NAME_PATTERN = "%s-es";

  public static final ConcurrentHashMap<String, String> VARIABLES = new ConcurrentHashMap<>();

  @Given("Kubernetes cluster is available")
  public void kubernetesClusterAvailable() {
    withK8sClient()
        .run(
            client -> {
              client.pods().inNamespace("default").list();
              log.info("K8s available");
              return null;
            });
  }

  @And("ECK-CR operator is installed")
  public void eckCRInstalled() {
            withK8sClient().run(client -> {
                var pods = client.pods().inNamespace("default").list().getItems();
                if (pods.stream().noneMatch(pod ->
     pod.getMetadata().getName().contains(ECK_CR_POD_NAME))) {
                    fail("ECK-CR not installed");
                }
                log.info("ECK-CR present");
                return null;
            });
  }

    @And("Elasticsearch {string} is available")
    public void elasticsearchAvailable(String esName) {
    withK8sClient()
        .run(
            client -> {
              var pods = client.pods().inNamespace("default").list().getItems();
              var esPodName = ES_POD_NAME_PATTERN.formatted(esName);
              if (pods.stream().noneMatch(pod -> pod.getMetadata().getName().contains(esPodName))) {
                fail("ES %s not installed".formatted(esName));
              }
              log.info("ES {} present", esName);
              return null;
            });
  }

  @Given("the {string} is set to {string}")
  public void variableSetter(String key, String value) {
    if (value.equals("random number")) {
      value = String.valueOf(ThreadLocalRandom.current().nextInt());
    }
    log.info("Setting {} to {}", key, value);
    VARIABLES.put(key, value);
  }

  @After
  public void cleanupVariables() {
    VARIABLES.clear();
  }
}
