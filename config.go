package appink8s

import (
	"crypto/x509"
	"errors"
	"fmt"
)

type filereader interface {
	ReadTokenFile() (string, error)
	ReadNamespaceFile() (string, error)
	ReadCertFile() ([]byte, error)
	ReadContainerID() (string, error)
}

type k8sconfig struct {
	token       string
	namespace   string
	certificate []byte
	filereader
}

func newK8sConfig() *k8sconfig {
	return &k8sconfig{
		filereader: newK8sFileReader(),
	}
}

func (c *k8sconfig) RunningInKubernetes() bool {
	t, err := c.Token()
	return err == nil && t != ""
}

func (c *k8sconfig) Token() (string, error) {
	if c.token != "" {
		return c.token, nil
	}

	token, err := c.ReadTokenFile()
	if err != nil {
		return "", fmt.Errorf("error retrieving service account token: %w", err)
	}

	c.token = token
	return c.token, nil
}

func (c *k8sconfig) CurrentNamespace() (string, error) {
	if c.namespace != "" {
		return c.namespace, nil
	}

	namespace, err := c.ReadNamespaceFile()
	if err != nil {
		return "", fmt.Errorf("error retrieving current namespace: %w", err)
	}

	c.namespace = namespace
	return c.namespace, nil
}

func (c *k8sconfig) Certificate() ([]byte, error) {
	if c.certificate != nil {
		return c.certificate, nil
	}

	certificate, err := c.ReadCertFile()
	if err != nil {
		return nil, fmt.Errorf("error retrieving certificate: %w", err)
	}

	c.certificate = certificate
	return c.certificate, nil
}

func (c *k8sconfig) CertPool() (*x509.CertPool, error) {
	ca, err := x509.SystemCertPool()
	if err != nil || ca == nil {
		ca = x509.NewCertPool()
	}

	cert, err := c.Certificate()
	if err != nil {
		return nil, err
	}

	ok := ca.AppendCertsFromPEM(cert)
	if !ok {
		return nil, errors.New("could not create certificate pool with dedicated certificate")
	}

	return ca, nil
}
