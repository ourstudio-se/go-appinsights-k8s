package appink8s

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type config_mockFileReader struct {
	token             string
	namespace         string
	cert              []byte
	container         string
	err               error
	callsForToken     int
	callsForNamespace int
	callsForCertFile  int
}

func (m *config_mockFileReader) ReadTokenFile() (string, error) {
	m.callsForToken = m.callsForToken + 1
	return m.token, m.err
}

func (m *config_mockFileReader) ReadNamespaceFile() (string, error) {
	m.callsForNamespace = m.callsForNamespace + 1
	return m.namespace, m.err
}

func (m *config_mockFileReader) ReadCertFile() ([]byte, error) {
	m.callsForCertFile = m.callsForCertFile + 1
	return m.cert, m.err
}

func (m *config_mockFileReader) ReadContainerID() (string, error) {
	return m.container, m.err
}

func Test_That_RunningInKubernetes_Is_Truthy_When_Token_Exist(t *testing.T) {
	cfg := &k8sconfig{
		filereader: &config_mockFileReader{
			token: "token",
		},
	}

	assert.True(t, cfg.RunningInKubernetes())
}

func Test_That_RunningInKubernetes_Is_Falsy_When_Token_Missing(t *testing.T) {
	cfg := &k8sconfig{
		filereader: &config_mockFileReader{
			token: "",
		},
	}

	assert.False(t, cfg.RunningInKubernetes())
}

func Test_That_RunningInKubernetes_Is_Falsy_On_Read_Error(t *testing.T) {
	cfg := &k8sconfig{
		filereader: &config_mockFileReader{
			token: "token",
			err:   errors.New("mock"),
		},
	}

	assert.False(t, cfg.RunningInKubernetes())
}

func Test_That_Token_Calls_Into_FileReader_Only_Once(t *testing.T) {
	fr := &config_mockFileReader{
		token: "token",
	}
	cfg := &k8sconfig{
		filereader: fr,
	}

	var err error
	_, err = cfg.Token()
	assert.NoError(t, err)
	_, err = cfg.Token()
	assert.NoError(t, err)
	_, err = cfg.Token()
	assert.NoError(t, err)

	assert.Equal(t, 1, fr.callsForToken)
}

func Test_That_CurrentNamespace_Calls_Into_FileReader_Only_Once(t *testing.T) {
	fr := &config_mockFileReader{
		namespace: "namespace",
	}
	cfg := &k8sconfig{
		filereader: fr,
	}

	var err error
	_, err = cfg.CurrentNamespace()
	assert.NoError(t, err)
	_, err = cfg.CurrentNamespace()
	assert.NoError(t, err)
	_, err = cfg.CurrentNamespace()
	assert.NoError(t, err)

	assert.Equal(t, 1, fr.callsForNamespace)
}

func Test_That_Certificate_Calls_Into_FileReader_Only_Once(t *testing.T) {
	fr := &config_mockFileReader{
		cert: []byte{},
	}
	cfg := &k8sconfig{
		filereader: fr,
	}

	var err error
	_, err = cfg.Certificate()
	assert.NoError(t, err)
	_, err = cfg.Certificate()
	assert.NoError(t, err)
	_, err = cfg.Certificate()
	assert.NoError(t, err)

	assert.Equal(t, 1, fr.callsForCertFile)
}
