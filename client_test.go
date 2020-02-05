package appink8s

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type client_mockHTTPClient struct {
	lastRequest *http.Request
	response    *http.Response
	err         error
}

func (m *client_mockHTTPClient) Do(r *http.Request) (*http.Response, error) {
	m.lastRequest = r
	return m.response, m.err
}

func Test_That_PodListURI_Returns_Correct_URI(t *testing.T) {
	namespace := "default"
	c := &k8sclient{
		k8sconfig: &k8sconfig{
			namespace: namespace,
		},
	}

	expected := fmt.Sprintf("https://kubernetes.default.svc/api/v1/namespaces/%s/pods", namespace)

	u, err := c.PodListURI()
	assert.NoError(t, err)
	assert.Equal(t, expected, u.String())
}

func Test_That_NodeListURI_Returns_Correct_URI(t *testing.T) {
	c := &k8sclient{
		k8sconfig: &k8sconfig{},
	}

	expected := "https://kubernetes.default.svc/api/v1/nodes"

	u, err := c.NodeListURI()
	assert.NoError(t, err)
	assert.Equal(t, expected, u.String())
}

func Test_That_GetPods_Requests_Correct_URI(t *testing.T) {
	namespace := "default"
	m := &client_mockHTTPClient{
		response: &http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"items": []}`))),
			StatusCode: 200,
		},
	}
	c := &k8sclient{
		httpclient: m,
		k8sconfig: &k8sconfig{
			token:     "token",
			namespace: namespace,
		},
	}

	_, err := c.GetPods()
	assert.NoError(t, err)

	expected, _ := url.Parse(fmt.Sprintf("https://kubernetes.default.svc/api/v1/namespaces/%s/pods", namespace))
	assert.Equal(t, expected, m.lastRequest.URL)
}

func Test_That_GetPods_Makes_Valid_Request(t *testing.T) {
	m := &client_mockHTTPClient{
		response: &http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"items": []}`))),
			StatusCode: 200,
		},
	}
	c := &k8sclient{
		httpclient: m,
		k8sconfig: &k8sconfig{
			token:     "token",
			namespace: "default",
		},
	}

	_, err := c.GetPods()
	assert.NoError(t, err)

	assert.Equal(t, "GET", m.lastRequest.Method)
	assert.Equal(t, "application/json", m.lastRequest.Header.Get("accept"))
}

func Test_That_GetPods_Uses_Bearer_Token_For_Request(t *testing.T) {
	token := "test-token"
	m := &client_mockHTTPClient{
		response: &http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"items": []}`))),
			StatusCode: 200,
		},
	}
	c := &k8sclient{
		httpclient: m,
		k8sconfig: &k8sconfig{
			token:     token,
			namespace: "default",
		},
	}

	_, err := c.GetPods()
	assert.NoError(t, err)

	expected := fmt.Sprintf("Bearer %s", token)
	assert.Equal(t, expected, m.lastRequest.Header.Get("authorization"))
}

func Test_That_GetNodes_Requests_Correct_URI(t *testing.T) {
	m := &client_mockHTTPClient{
		response: &http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"items": []}`))),
			StatusCode: 200,
		},
	}
	c := &k8sclient{
		httpclient: m,
		k8sconfig: &k8sconfig{
			token:     "token",
			namespace: "namespace",
		},
	}

	_, err := c.GetNodes()
	assert.NoError(t, err)

	expected, _ := url.Parse("https://kubernetes.default.svc/api/v1/nodes")
	assert.Equal(t, expected, m.lastRequest.URL)
}

func Test_That_GetNodes_Makes_Valid_Request(t *testing.T) {
	m := &client_mockHTTPClient{
		response: &http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"items": []}`))),
			StatusCode: 200,
		},
	}
	c := &k8sclient{
		httpclient: m,
		k8sconfig: &k8sconfig{
			token:     "token",
			namespace: "default",
		},
	}

	_, err := c.GetNodes()
	assert.NoError(t, err)

	assert.Equal(t, "GET", m.lastRequest.Method)
	assert.Equal(t, "application/json", m.lastRequest.Header.Get("accept"))
}

func Test_That_GetNodes_Uses_Bearer_Token_For_Request(t *testing.T) {
	token := "test-token"
	m := &client_mockHTTPClient{
		response: &http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"items": []}`))),
			StatusCode: 200,
		},
	}
	c := &k8sclient{
		httpclient: m,
		k8sconfig: &k8sconfig{
			token:     token,
			namespace: "default",
		},
	}

	_, err := c.GetNodes()
	assert.NoError(t, err)

	expected := fmt.Sprintf("Bearer %s", token)
	assert.Equal(t, expected, m.lastRequest.Header.Get("authorization"))
}
