package appink8s

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const k8sHostAddress = "https://kubernetes.default.svc"
const k8sNodeURI = "api/v1/nodes"
const k8sPodURI = "api/v1/namespaces/%s/pods"

type httpclient interface {
	Do(*http.Request) (*http.Response, error)
}

type k8sclient struct {
	httpclient
	*k8sconfig
}

func newK8sClient(cfg *k8sconfig) (*k8sclient, error) {
	ca, err := cfg.CertPool()
	if err != nil {
		return nil, err
	}

	tls := &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            ca,
	}
	tr := &http.Transport{
		TLSClientConfig: tls,
	}

	c := &http.Client{
		Transport: tr,
	}

	return &k8sclient{
		httpclient: c,
		k8sconfig:  cfg,
	}, nil
}

func (c *k8sclient) PodListURI() (*url.URL, error) {
	namespace, err := c.CurrentNamespace()
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf(k8sPodURI, namespace)
	u := fmt.Sprintf("%s/%s", k8sHostAddress, path)
	return url.Parse(u)
}

func (c *k8sclient) NodeListURI() (*url.URL, error) {
	u := fmt.Sprintf("%s/%s", k8sHostAddress, k8sNodeURI)
	return url.Parse(u)
}

func (c *k8sclient) GetPods() (*podListSpec, error) {
	u, err := c.PodListURI()
	if err != nil {
		return nil, fmt.Errorf("error parsing pod list URI: %w", err)
	}

	b, err := c.request(u)
	if err != nil {
		return nil, fmt.Errorf("error reading pod list spec: %w", err)
	}

	var specs podListSpec
	if err = json.Unmarshal(b, &specs); err != nil {
		return nil, fmt.Errorf("error parsing pod list spec: %w", err)
	}

	return &specs, nil
}

func (c *k8sclient) GetNodes() (*nodeListSpec, error) {
	u, err := c.NodeListURI()
	if err != nil {
		return nil, fmt.Errorf("error parsing node list URI: %w", err)
	}

	b, err := c.request(u)
	if err != nil {
		return nil, fmt.Errorf("error reading node list spec: %w", err)
	}

	var specs nodeListSpec
	if err = json.Unmarshal(b, &specs); err != nil {
		return nil, fmt.Errorf("error parsing node list spec: %w", err)
	}

	return &specs, nil
}

func (c *k8sclient) request(u *url.URL) ([]byte, error) {
	token, err := c.Token()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token))

	if err != nil {
		return nil, fmt.Errorf("unable to create request URI: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("unable to request Kubernetes: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("unable to read Kubernetes API, received status code: %d", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read Kubernetes data: %v", err)
	}

	return b, nil
}
