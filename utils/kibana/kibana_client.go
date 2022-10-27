package kibana

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"strings"

	configv2 "github.com/xco-sk/eck-custom-resources/apis/config/v2"
	"github.com/xco-sk/eck-custom-resources/utils"
	k8sv1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Client struct {
	Cli        client.Client
	Ctx        context.Context
	KibanaSpec configv2.KibanaSpec
	Req        ctrl.Request
}

func (kClient Client) DoGet(path string) (*http.Response, error) {
	httpRequest, err := http.NewRequest("GET", kClient.KibanaSpec.Url+path, nil)
	if err != nil {
		return nil, err
	}

	return kClient.doRequest(httpRequest)
}

func (kClient Client) DoPost(path string, body string) (*http.Response, error) {
	httpRequest, err := http.NewRequest("POST", kClient.KibanaSpec.Url+path, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	return kClient.doRequest(httpRequest)
}

func (kClient Client) DoPut(path string, body string) (*http.Response, error) {
	httpRequest, err := http.NewRequest("PUT", kClient.KibanaSpec.Url+path, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	return kClient.doRequest(httpRequest)
}

func (kClient Client) DoDelete(path string) (*http.Response, error) {
	httpRequest, err := http.NewRequest("DELETE", kClient.KibanaSpec.Url+path, nil)
	if err != nil {
		return nil, err
	}

	return kClient.doRequest(httpRequest)
}

func (kClient Client) getHttpClient() (*http.Client, error) {

	tr := &http.Transport{}

	if kClient.KibanaSpec.Certificate != nil {
		var certificateSecret k8sv1.Secret
		if err := utils.GetCertificateSecret(kClient.Cli, kClient.Ctx, kClient.Req.Namespace, kClient.KibanaSpec.Certificate, &certificateSecret); err != nil {
			return nil, err
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(certificateSecret.Data[kClient.KibanaSpec.Certificate.CertificateKey])

		tr.TLSClientConfig = &tls.Config{
			RootCAs:            caCertPool,
			InsecureSkipVerify: false,
		}
	} else if strings.HasPrefix(kClient.KibanaSpec.Url, "https://") {
		return nil, errors.New("Failed to configure http client, certificate not configured (kibana.certificate)")
	}

	httpClient := &http.Client{
		Transport: tr,
	}

	return httpClient, nil
}

func (kClient Client) doRequest(httpRequest *http.Request) (*http.Response, error) {
	if kClient.KibanaSpec.Authentication != nil && kClient.KibanaSpec.Authentication.UsernamePassword != nil {
		var userSecret k8sv1.Secret
		if err := utils.GetUserSecret(kClient.Cli, kClient.Ctx, kClient.Req.Namespace, kClient.KibanaSpec.Authentication.UsernamePassword, &userSecret); err != nil {
			return nil, err
		}
		httpRequest.SetBasicAuth(kClient.KibanaSpec.Authentication.UsernamePassword.UserName, string(userSecret.Data[kClient.KibanaSpec.Authentication.UsernamePassword.UserName]))
	}

	httpClient, err := kClient.getHttpClient()
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set("kbn-xsrf", "true")
	response, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	return response, nil
}
