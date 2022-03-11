package utils

import (
	"context"
	eseckv1 "github.com/xco-sk/eck-custom-resources/api/v1"
	k8sv1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const EckUserSecretSuffix = "-es-elastic-user"
const EckHttpServiceSuffix = "-es-http"

func GetElasticSecret(cli client.Client, ctx context.Context, namespace string, index *eseckv1.Index, secret *k8sv1.Secret) error {
	if err := cli.Get(ctx, client.ObjectKey{Namespace: namespace, Name: generateElasticSecretName(index.Spec.TargetCluster.EckCluster.ClusterName)}, secret); err != nil {
		return err
	}
	return nil
}

func GenerateElasticEndpoint(clusterName string, namespace string) string {
	return "https://" + clusterName + EckHttpServiceSuffix + "." + namespace + ":9200"
}

func generateElasticSecretName(clusterName string) string {
	return clusterName + EckUserSecretSuffix
}
