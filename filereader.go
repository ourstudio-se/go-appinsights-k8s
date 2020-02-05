package appink8s

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

const k8sTokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
const k8sNamespacePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
const k8sCertPath = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"

type k8sfiles struct{}

func newK8sFileReader() *k8sfiles {
	return &k8sfiles{}
}

func (kf *k8sfiles) ReadTokenFile() (string, error) {
	token, err := ioutil.ReadFile(k8sTokenPath)
	if err != nil {
		return "", nil
	}

	return string(token), nil
}

func (kf *k8sfiles) ReadNamespaceFile() (string, error) {
	namespace, err := ioutil.ReadFile(k8sNamespacePath)
	if err != nil {
		return "", nil
	}

	return string(namespace), nil
}

func (kf *k8sfiles) ReadCertFile() ([]byte, error) {
	return ioutil.ReadFile(k8sCertPath)
}

func (kf *k8sfiles) ReadContainerID() (string, error) {
	raw, err := ioutil.ReadFile(k8sContainerInfoPath)
	if err != nil {
		return "", fmt.Errorf("could not read container ID: %w", err)
	}

	cg := string(raw)
	id, err := parseContainerIDFromCGroupInfo(cg)

	if err != nil {
		return "", fmt.Errorf("could not parse container ID: %w", err)
	}

	return id, nil
}

func parseContainerIDFromCGroupInfo(raw string) (string, error) {
	re := regexp.MustCompile(`(?im)cpu.+/([^/]*)$`)
	path := re.FindString(raw)

	if path == "" {
		return "", errors.New("could not find container ID")
	}

	parts := strings.Split(path, "/")
	return parts[len(parts)-1], nil
}
