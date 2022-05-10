package kibana

import (
	"context"
	"crypto/tls"
	configv2 "github.com/xco-sk/eck-custom-resources/apis/config/v2"
	"github.com/xco-sk/eck-custom-resources/utils"
	k8sv1 "k8s.io/api/core/v1"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

func DoGet(cli client.Client, ctx context.Context, kibanaSpec configv2.KibanaSpec, req ctrl.Request, url string) (*http.Response, error) {
	httpRequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return doRequest(cli, ctx, kibanaSpec, req, httpRequest)
}

func DoPost(cli client.Client, ctx context.Context, kibanaSpec configv2.KibanaSpec, req ctrl.Request, url string, body string) (*http.Response, error) {
	httpRequest, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	return doRequest(cli, ctx, kibanaSpec, req, httpRequest)
}

func DoPut(cli client.Client, ctx context.Context, kibanaSpec configv2.KibanaSpec, req ctrl.Request, url string, body string) (*http.Response, error) {
	httpRequest, err := http.NewRequest("PUT", url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	return doRequest(cli, ctx, kibanaSpec, req, httpRequest)
}

func DoDelete(cli client.Client, ctx context.Context, kibanaSpec configv2.KibanaSpec, req ctrl.Request, url string) (*http.Response, error) {
	httpRequest, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return doRequest(cli, ctx, kibanaSpec, req, httpRequest)
}

func getHttpClient(cli client.Client, ctx context.Context, kibanaSpec configv2.KibanaSpec, req ctrl.Request) (*http.Client, error) {

	tr := &http.Transport{}

	if kibanaSpec.Certificate != nil {
		var certificateSecret k8sv1.Secret
		if err := utils.GetCertificateSecret(cli, ctx, req.Namespace, kibanaSpec.Certificate, &certificateSecret); err != nil {
			return nil, err
		}
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // TODO: use real certificate
	}

	httpClient := &http.Client{
		Transport: tr,
	}

	return httpClient, nil
}

func doRequest(cli client.Client, ctx context.Context, kibanaSpec configv2.KibanaSpec, req ctrl.Request, httpRequest *http.Request) (*http.Response, error) {
	if kibanaSpec.Authentication != nil && kibanaSpec.Authentication.UsernamePassword != nil {
		var userSecret k8sv1.Secret
		if err := utils.GetUserSecret(cli, ctx, req.Namespace, kibanaSpec.Authentication.UsernamePassword, &userSecret); err != nil {
			return nil, err
		}
		httpRequest.SetBasicAuth(kibanaSpec.Authentication.UsernamePassword.UserName, string(userSecret.Data[kibanaSpec.Authentication.UsernamePassword.UserName]))
	}

	httpClient, err := getHttpClient(cli, ctx, kibanaSpec, req)
	if err != nil {
		return nil, err
	}

	response, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	return response, nil
}
