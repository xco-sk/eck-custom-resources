package sk.xco.eckcr.util;

import co.elastic.clients.elasticsearch.ElasticsearchClient;
import co.elastic.clients.elasticsearch.indices.GetIndexRequest;
import co.elastic.clients.elasticsearch.indices.IndexState;
import co.elastic.clients.json.jackson.JacksonJsonpMapper;
import co.elastic.clients.transport.ElasticsearchTransport;
import co.elastic.clients.transport.rest_client.RestClientTransport;
import org.apache.http.HttpHost;
import org.elasticsearch.client.RestClient;

import java.io.IOException;

public class ESClient {

    public static IndexState getIndexState(String indexName) throws IOException {
        return getClient().indices().get(new GetIndexRequest.Builder().index(indexName).build()).get(indexName);
    }

    private static ElasticsearchClient getClient() {

        RestClient restClient = RestClient.builder(
                new HttpHost("localhost", 9200)).build();

        ElasticsearchTransport transport = new RestClientTransport(
                restClient, new JacksonJsonpMapper());

        return new ElasticsearchClient(transport);
    }
}
