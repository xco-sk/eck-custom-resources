package sk.xco.eckcr.step;

import io.cucumber.java.After;
import io.cucumber.java.ParameterType;
import io.cucumber.java.en.Given;
import io.fabric8.kubernetes.client.KubernetesClientException;
import lombok.extern.slf4j.Slf4j;
import sk.xco.eckcr.util.ApiType;

import java.util.HashSet;
import java.util.Set;

import static sk.xco.eckcr.util.K8sClient.withK8sClient;

@Slf4j
public class K8sResource {
    private final Set<String> toCleanup = new HashSet<>();

    @Given("the {ApiType} {string} is applied")
    public void applyResource(ApiType apiType, String resourceName) {
        withK8sClient().run(client -> {
            client.load(K8sResource.class.getResourceAsStream(getResourcePath(resourceName)))
                    .inNamespace("default")
                    .create();
            return null;
        });
        toCleanup.add(resourceName);
    }

    @After
    public void afterScenario() {
        withK8sClient().run(c -> {
            toCleanup.forEach(r -> {
                log.info("Cleaning-up {}", r);
                try {
                    c.load(K8sResource.class.getResourceAsStream(getResourcePath(r)))
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
