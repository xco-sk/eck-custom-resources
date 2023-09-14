package sk.xco.eckcr.util;

import static org.junit.jupiter.api.Assertions.fail;

import io.fabric8.kubernetes.client.KubernetesClient;
import io.fabric8.kubernetes.client.KubernetesClientBuilder;
import io.fabric8.kubernetes.client.KubernetesClientException;
import java.util.Base64;
import java.util.Map;
import java.util.function.Function;
import lombok.AccessLevel;
import lombok.NoArgsConstructor;

@NoArgsConstructor(access = AccessLevel.PRIVATE)
public class K8sClient {

  public static final String NAMESPACE = "default";
  private static final K8sClient INSTANCE = new K8sClient();
  private static final String ES_USER_SECRET_PATTERN = "%s-es-elastic-user";
  private static final String ES_CERTS_SECRET_PATTERN = "%s-es-http-certs-internal";

  public static K8sClient withK8sClient() {
    return INSTANCE;
  }

  public <T> T run(Function<KubernetesClient, T> fn) {
    var builder = new KubernetesClientBuilder();
    try (KubernetesClient client = builder.build()) {
      return fn.apply(client);
    } catch (KubernetesClientException e) {
      fail("Kubernetes client exception", e);
    } catch (Exception e) {
      fail("Exception occurred during k8s client operation", e);
    }
    throw new IllegalStateException("Failed to execute");
  }

  protected static ElasticsearchUser getElasticsearchUserFromSecret(String esName) {
    return withK8sClient()
        .run(
            c ->
                c.secrets().inNamespace(NAMESPACE).list().getItems().stream()
                    .filter(
                        s ->
                            ES_USER_SECRET_PATTERN
                                .formatted(esName)
                                .equals(s.getMetadata().getName()))
                    .findFirst()
                    .map(secret -> secret.getData().entrySet())
                    .map(
                        data ->
                            data.stream()
                                .findFirst()
                                .orElseThrow(
                                    () ->
                                        new IllegalArgumentException(
                                            "Secret does not contain data")))
                    .map(
                        entry ->
                            new ElasticsearchUser(
                                entry.getKey(),
                                new String(Base64.getDecoder().decode(entry.getValue()))))
                    .orElseThrow(
                        () -> new IllegalArgumentException("No secret with ES user data found")));
  }

  protected static String getElasticsearchCACertFromSecret(String esName) {
    return withK8sClient()
        .run(
            c ->
                c.secrets().inNamespace(NAMESPACE).list().getItems().stream()
                    .filter(
                        s ->
                            ES_CERTS_SECRET_PATTERN
                                .formatted(esName)
                                .equals(s.getMetadata().getName()))
                    .findFirst()
                    .map(secret -> secret.getData().entrySet())
                    .map(
                        data ->
                            data.stream()
                                .filter(e -> e.getKey().equals("ca.crt"))
                                .findFirst()
                                .orElseThrow(
                                    () ->
                                        new IllegalArgumentException(
                                            "Secret does not contain data")))
                    .map(Map.Entry::getValue)
                    .map(encoded -> new String(Base64.getDecoder().decode(encoded)))
                    .orElseThrow(
                        () -> new IllegalArgumentException("No secret with ES user data found")));
  }
}
