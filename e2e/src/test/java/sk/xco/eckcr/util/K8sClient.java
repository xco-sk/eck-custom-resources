package sk.xco.eckcr.util;

import io.fabric8.kubernetes.client.KubernetesClient;
import io.fabric8.kubernetes.client.KubernetesClientBuilder;
import io.fabric8.kubernetes.client.KubernetesClientException;
import lombok.AccessLevel;
import lombok.NoArgsConstructor;

import java.util.function.Function;

import static org.junit.jupiter.api.Assertions.fail;

@NoArgsConstructor(access = AccessLevel.PRIVATE)
public class K8sClient {

    public static final String NAMESPACE = "default";
    private static final K8sClient INSTANCE = new K8sClient();

    public static K8sClient withK8sClient() {
        return INSTANCE;
    }

    public void run(Function<KubernetesClient, Void> fn) {
        var builder = new KubernetesClientBuilder();
        try (KubernetesClient client = builder.build()) {
            fn.apply(client);
        } catch (KubernetesClientException e) {
            fail("Kubernetes client exception", e);
        } catch (Exception e) {
            fail("Exception occurred during k8s client operation", e);
        }
    }
}
