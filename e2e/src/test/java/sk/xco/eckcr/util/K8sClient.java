package sk.xco.eckcr.util;

import io.fabric8.kubernetes.client.KubernetesClient;
import io.fabric8.kubernetes.client.KubernetesClientBuilder;
import io.fabric8.kubernetes.client.KubernetesClientException;
import lombok.AccessLevel;
import lombok.NoArgsConstructor;
import lombok.RequiredArgsConstructor;

import java.util.function.Function;

import static org.junit.jupiter.api.Assertions.fail;

@NoArgsConstructor(access = AccessLevel.PRIVATE)
public class K8sClient {

    private static final K8sClient INSTANCE = new K8sClient();

    public static K8sClient withK8sClient() {
        return INSTANCE;
    }

    public void run(Function<KubernetesClient, Void> fn) {
        var builder = new KubernetesClientBuilder();
        try (KubernetesClient client = builder.build()) {
            fn.apply(client);
        } catch (KubernetesClientException e) {
            fail("Kubernetes not available");
            throw new IllegalStateException("Kubernetes not available");
        }
    }
}
