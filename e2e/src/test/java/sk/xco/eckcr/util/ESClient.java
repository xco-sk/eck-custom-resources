package sk.xco.eckcr.util;

import co.elastic.clients.elasticsearch.ElasticsearchClient;
import co.elastic.clients.elasticsearch.ilm.GetLifecycleRequest;
import co.elastic.clients.elasticsearch.ilm.IlmPolicy;
import co.elastic.clients.elasticsearch.indices.*;
import co.elastic.clients.elasticsearch.indices.get_index_template.IndexTemplateItem;
import co.elastic.clients.json.jackson.JacksonJsonpMapper;
import co.elastic.clients.transport.ElasticsearchTransport;
import co.elastic.clients.transport.TransportUtils;
import co.elastic.clients.transport.rest_client.RestClientTransport;
import java.io.ByteArrayInputStream;
import java.io.IOException;
import javax.net.ssl.SSLContext;
import lombok.extern.slf4j.Slf4j;
import org.apache.http.HttpHost;
import org.apache.http.auth.AuthScope;
import org.apache.http.auth.UsernamePasswordCredentials;
import org.apache.http.impl.client.BasicCredentialsProvider;
import org.elasticsearch.client.RestClient;

@Slf4j
public class ESClient {

  public static final String DEFAULT_ES_NAME = "quickstart";

  public static IndexState getIndexState(String indexName) throws IOException {
    return getClient()
        .indices()
        .get(new GetIndexRequest.Builder().index(indexName).build())
        .get(indexName);
  }

  public static IndexTemplateItem getTemplate(String templateName) throws IOException {
    return getClient()
        .indices()
        .getIndexTemplate(new GetIndexTemplateRequest.Builder().name(templateName).build())
        .indexTemplates()
        .stream()
        .findFirst()
        .get();
  }

  public static IlmPolicy getIlmPolicy(String policyName) throws IOException {
    return getClient()
        .ilm()
        .getLifecycle(new GetLifecycleRequest.Builder().name(policyName).build())
        .get(policyName)
        .policy();
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
}
