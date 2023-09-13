package sk.xco.eckcr.util;

import static java.util.Objects.nonNull;
import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.fail;

import co.elastic.clients.elasticsearch.ElasticsearchClient;
import co.elastic.clients.elasticsearch._types.ElasticsearchException;
import co.elastic.clients.elasticsearch.ilm.GetLifecycleRequest;
import co.elastic.clients.elasticsearch.ilm.IlmPolicy;
import co.elastic.clients.elasticsearch.indices.*;
import co.elastic.clients.elasticsearch.indices.get_index_template.IndexTemplateItem;
import co.elastic.clients.elasticsearch.ingest.GetPipelineRequest;
import co.elastic.clients.elasticsearch.ingest.Pipeline;
import co.elastic.clients.elasticsearch.slm.SnapshotLifecycle;
import co.elastic.clients.elasticsearch.snapshot.GetRepositoryRequest;
import co.elastic.clients.elasticsearch.snapshot.Repository;
import co.elastic.clients.json.jackson.JacksonJsonpMapper;
import co.elastic.clients.transport.ElasticsearchTransport;
import co.elastic.clients.transport.TransportUtils;
import co.elastic.clients.transport.rest_client.RestClientTransport;
import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.util.concurrent.TimeUnit;
import java.util.function.Function;
import javax.net.ssl.SSLContext;
import lombok.extern.slf4j.Slf4j;
import org.apache.http.HttpHost;
import org.apache.http.auth.AuthScope;
import org.apache.http.auth.UsernamePasswordCredentials;
import org.apache.http.impl.client.BasicCredentialsProvider;
import org.awaitility.Awaitility;
import org.elasticsearch.client.RestClient;

@Slf4j
public class ESClient {

  public static final String DEFAULT_ES_NAME = "quickstart";

  public static IndexState getIndexState(String indexName) {
    try {
      return getClient()
          .indices()
          .get(new GetIndexRequest.Builder().index(indexName).build())
          .get(indexName);
    } catch (IOException e) {
      throw new RuntimeException(e);
    }
  }

  public static IndexTemplateItem getTemplate(String templateName) {
    try {
      return getClient()
          .indices()
          .getIndexTemplate(new GetIndexTemplateRequest.Builder().name(templateName).build())
          .indexTemplates()
          .stream()
          .findFirst()
          .get();
    } catch (IOException e) {
      throw new RuntimeException(e);
    }
  }

  public static IlmPolicy getIlmPolicy(String policyName) {
    try {
      return getClient()
          .ilm()
          .getLifecycle(new GetLifecycleRequest.Builder().name(policyName).build())
          .get(policyName)
          .policy();
    } catch (IOException e) {
      throw new RuntimeException(e);
    }
  }

  public static Pipeline getIngestPipeline(String pipelineName) {
    try {
      return getClient()
          .ingest()
          .getPipeline(new GetPipelineRequest.Builder().id(pipelineName).build())
          .get(pipelineName);
    } catch (IOException e) {
      throw new RuntimeException(e);
    }
  }

  public static Repository getSnapshotRepo(String repoName) {
    try {
      return getClient()
          .snapshot()
          .getRepository(new GetRepositoryRequest.Builder().name(repoName).build())
          .get(repoName);
    } catch (IOException e) {
      throw new RuntimeException(e);
    }
  }

  public static SnapshotLifecycle getSnapshotLifecyclePolicy(String policyName) {
    try {
      return getClient()
          .slm()
          .getLifecycle(
              new co.elastic.clients.elasticsearch.slm.GetLifecycleRequest.Builder()
                  .policyId(policyName)
                  .build())
          .get(policyName);
    } catch (IOException e) {
      throw new RuntimeException(e);
    }
  }

  private static ElasticsearchClient getClient() {
    var user = K8sClient.getElasticsearchUserFromSecret(DEFAULT_ES_NAME);
    var caCrt = K8sClient.getElasticsearchCACertFromSecret(DEFAULT_ES_NAME);

    BasicCredentialsProvider credsProv = new BasicCredentialsProvider();
    credsProv.setCredentials(
        AuthScope.ANY, new UsernamePasswordCredentials(user.username(), user.password()));

    SSLContext sslContext =
        TransportUtils.sslContextFromHttpCaCrt(new ByteArrayInputStream(caCrt.getBytes()));

    RestClient restClient =
        RestClient.builder(new HttpHost("quickstart-es-http", 9200, "https"))
            .setHttpClientConfigCallback(
                hc -> hc.setSSLContext(sslContext).setDefaultCredentialsProvider(credsProv))
            .build();

    ElasticsearchTransport transport =
        new RestClientTransport(restClient, new JacksonJsonpMapper());

    return new ElasticsearchClient(transport);
  }

  public static <T> void awaitResourceNotPresent(
      String resourceName, Function<String, T> getResourceFunction) {
    Await.untilAsserted(
        () -> {
          try {
            T resource = getResourceFunction.apply(resourceName);
            if (nonNull(resource)) {
              fail("Resource %s present in Elasticsearch: %s".formatted(resourceName, resource));
            }
          } catch (ElasticsearchException e) {
            assertThat(e.status()).isEqualTo(404);
          }
        });
  }

  public static <T> void waitForResource(
      String resourceName, Function<String, T> getResourceFunction) {
    Awaitility.await()
        .atMost(10, TimeUnit.SECONDS)
        .until(
            () -> {
              try {
                getResourceFunction.apply(resourceName);
                return true;
              } catch (ElasticsearchException e) {
                return false;
              }
            });
  }
}
