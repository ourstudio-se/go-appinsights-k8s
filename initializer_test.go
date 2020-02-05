package appink8s

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type initializer_mockFileReader struct {
	token     string
	namespace string
	cert      []byte
	container string
	err       error
}

func (m *initializer_mockFileReader) ReadTokenFile() (string, error) {
	return m.token, m.err
}

func (m *initializer_mockFileReader) ReadNamespaceFile() (string, error) {
	return m.namespace, m.err
}

func (m *initializer_mockFileReader) ReadCertFile() ([]byte, error) {
	return m.cert, m.err
}

func (m *initializer_mockFileReader) ReadContainerID() (string, error) {
	return m.container, m.err
}

type initializer_mockHTTPClient struct{}

func (m *initializer_mockHTTPClient) Do(r *http.Request) (*http.Response, error) {
	reqpath := strings.Split(r.URL.String(), "/")

	if reqpath[len(reqpath)-1] == "nodes" {
		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(k8sNodeResponse))),
			StatusCode: 200,
		}, nil
	}

	return &http.Response{
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(k8sPodResponse))),
		StatusCode: 200,
	}, nil
}

func Test_That_ReadPropertySpec_Parses_Kubernetes_Response(t *testing.T) {
	cfg := &k8sconfig{
		filereader: &initializer_mockFileReader{
			container: "TEST-CONTAINER-ID",
		},
	}
	c := &k8sclient{
		httpclient: &initializer_mockHTTPClient{},
		k8sconfig:  cfg,
	}
	i := &k8sinitializer{
		client: c,
	}

	spec, err := i.ReadPropertySpec()

	assert.NoError(t, err)
	assert.NotNil(t, spec)
	assert.Equal(t, "TEST-POD-ID", spec.PodID)
	assert.Equal(t, "TEST-POD-NAME", spec.PodName)
	assert.Equal(t, "TEST-CONTAINER-ID", spec.ContainerID)
	assert.Equal(t, "TEST-CONTAINER-NAME", spec.ContainerName)
	assert.Equal(t, "TEST-DEPLOYMENT-NAME-REPLICASETID", spec.ReplicaSetName)
	assert.Equal(t, "TEST-DEPLOYMENT-NAME", spec.DeploymentName)
	assert.Equal(t, "TEST-NODE-ID", spec.NodeID)
	assert.Equal(t, "TEST-NODE-NAME", spec.NodeName)
}

const k8sNodeResponse = `{
	"kind": "NodeList",
	"apiVersion": "v1",
	"metadata": {
	  "selfLink": "/api/v1/nodes",
	  "resourceVersion": "20109945"
	},
	"items": [
	  {
		"metadata": {
		  "name": "TEST-NODE-NAME",
		  "selfLink": "/api/v1/nodes/agentpool-0",
		  "uid": "TEST-NODE-ID",
		  "resourceVersion": "20109871",
		  "creationTimestamp": "2019-09-10T08:06:20Z",
		  "labels": {
			"agentpool": "agentpool",
			"beta.kubernetes.io/arch": "amd64",
			"beta.kubernetes.io/instance-type": "Standard_D2s_v3",
			"beta.kubernetes.io/os": "linux",
			"failure-domain.beta.kubernetes.io/region": "westeurope",
			"failure-domain.beta.kubernetes.io/zone": "1",
			"kubernetes.azure.com/cluster": "MC_cluster",
			"kubernetes.azure.com/role": "agent",
			"kubernetes.io/arch": "amd64",
			"kubernetes.io/hostname": "agentpool-0",
			"kubernetes.io/os": "linux",
			"kubernetes.io/role": "agent",
			"node-role.kubernetes.io/agent": "",
			"storageprofile": "managed",
			"storagetier": "Premium_LRS"
		  },
		  "annotations": {
			"node.alpha.kubernetes.io/ttl": "0",
			"volumes.kubernetes.io/controller-managed-attach-detach": "true"
		  }
		},
		"spec": {
		  "podCIDR": "10.244.1.0/24",
		  "providerID": "azure:///subscriptions/vm/agentpool-0"
		},
		"status": {
		  "capacity": {
			"attachable-volumes-azure-disk": "4",
			"cpu": "2",
			"ephemeral-storage": "50758760Ki",
			"hugepages-1Gi": "0",
			"hugepages-2Mi": "0",
			"memory": "8145348Ki",
			"pods": "110"
		  },
		  "allocatable": {
			"attachable-volumes-azure-disk": "4",
			"cpu": "1234m",
			"ephemeral-storage": "46779273139",
			"hugepages-1Gi": "0",
			"hugepages-2Mi": "0",
			"memory": "5490116Ki",
			"pods": "110"
		  },
		  "conditions": [
			{
			  "type": "NetworkUnavailable",
			  "status": "False",
			  "lastHeartbeatTime": "2019-09-10T08:08:58Z",
			  "lastTransitionTime": "2019-09-10T08:08:58Z",
			  "reason": "RouteCreated",
			  "message": "RouteController created a route"
			},
			{
			  "type": "MemoryPressure",
			  "status": "False",
			  "lastHeartbeatTime": "2020-02-04T09:01:40Z",
			  "lastTransitionTime": "2020-01-04T02:12:27Z",
			  "reason": "KubeletHasSufficientMemory",
			  "message": "kubelet has sufficient memory available"
			},
			{
			  "type": "DiskPressure",
			  "status": "False",
			  "lastHeartbeatTime": "2020-02-04T09:01:40Z",
			  "lastTransitionTime": "2019-10-08T23:48:04Z",
			  "reason": "KubeletHasNoDiskPressure",
			  "message": "kubelet has no disk pressure"
			},
			{
			  "type": "PIDPressure",
			  "status": "False",
			  "lastHeartbeatTime": "2020-02-04T09:01:40Z",
			  "lastTransitionTime": "2019-10-08T23:48:04Z",
			  "reason": "KubeletHasSufficientPID",
			  "message": "kubelet has sufficient PID available"
			},
			{
			  "type": "Ready",
			  "status": "True",
			  "lastHeartbeatTime": "2020-02-04T09:01:40Z",
			  "lastTransitionTime": "2019-11-27T07:26:19Z",
			  "reason": "KubeletReady",
			  "message": "kubelet is posting ready status. AppArmor enabled"
			}
		  ],
		  "addresses": [
			{
			  "type": "Hostname",
			  "address": "agentpool-0"
			},
			{
			  "type": "InternalIP",
			  "address": "10.240.0.5"
			}
		  ],
		  "daemonEndpoints": {
			"kubeletEndpoint": {
			  "Port": 10250
			}
		  },
		  "nodeInfo": {
			"machineID": "MACHINE-ID",
			"systemUUID": "B2B5000D-33F1-5D43-91C1-B481B9F6C1B9",
			"bootID": "1c2d2ccc-d748-48f3-ccab-5ac8db3c9d7f",
			"kernelVersion": "4.15.0-1052-azure",
			"osImage": "Ubuntu 16.04.6 LTS",
			"containerRuntimeVersion": "docker://3.0.6",
			"kubeletVersion": "v1.14.6",
			"kubeProxyVersion": "v1.14.6",
			"operatingSystem": "linux",
			"architecture": "amd64"
		  },
		  "images": []
		}
	  }
	]
  }`

const k8sPodResponse = `{
	"kind": "PodList",
	"apiVersion": "v1",
	"metadata": {
	  "selfLink": "/api/v1/namespaces/default/pods",
	  "resourceVersion": "6434072"
	},
	"items": [
	  {
		"metadata": {
		  "name": "TEST-POD-NAME",
		  "generateName": "TEST-POD-NAME-86b784d44c-",
		  "namespace": "default",
		  "selfLink": "/api/v1/namespaces/default/pods/TEST-POD-NAME-86b784d44c-xxvpw",
		  "uid": "TEST-POD-ID",
		  "resourceVersion": "5610449",
		  "creationTimestamp": "2020-01-23T08:18:17Z",
		  "labels": {
			"test/label": "test-label"
		  },
		  "annotations": {
			"test/annotation": "test-annotation"
		  },
		  "ownerReferences": [
			{
			  "apiVersion": "apps/v1",
			  "kind": "ReplicaSet",
			  "name": "TEST-DEPLOYMENT-NAME-REPLICASETID",
			  "uid": "e8987468-3db8-11ea-a877-22acad587db4",
			  "controller": true,
			  "blockOwnerDeletion": true
			}
		  ]
		},
		"spec": {
		  "volumes": [
			{
			  "name": "default-token-lp9pd",
			  "secret": {
				"secretName": "default-token-lp9pd",
				"defaultMode": 420
			  }
			}
		  ],
		  "initContainers": [],
		  "containers": [
			{
			  "name": "test-application",
			  "image": "registry.docker.io/images/test-application:ddbdd5056b530d6dac910f2ab49e219fcaf46dae",
			  "ports": [
				{
				  "name": "http",
				  "containerPort": 8080,
				  "protocol": "TCP"
				}
			  ],
			  "envFrom": [
				{
				  "configMapRef": {
					"name": "TEST-CONFIG-MAP"
				  }
				}
			  ],
			  "env": [],
			  "resources": {},
			  "volumeMounts": [
				{
				  "name": "default-token-lp9pd",
				  "readOnly": true,
				  "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
				}
			  ],
			  "terminationMessagePath": "/dev/termination-log",
			  "terminationMessagePolicy": "File",
			  "imagePullPolicy": "Always"
			}
		  ],
		  "restartPolicy": "Always",
		  "terminationGracePeriodSeconds": 30,
		  "dnsPolicy": "ClusterFirst",
		  "serviceAccountName": "default",
		  "serviceAccount": "default",
		  "nodeName": "TEST-NODE-NAME",
		  "securityContext": {},
		  "imagePullSecrets": [],
		  "affinity": {},
		  "schedulerName": "default-scheduler",
		  "tolerations": [
			{
			  "key": "node.kubernetes.io/not-ready",
			  "operator": "Exists",
			  "effect": "NoExecute",
			  "tolerationSeconds": 300
			},
			{
			  "key": "node.kubernetes.io/unreachable",
			  "operator": "Exists",
			  "effect": "NoExecute",
			  "tolerationSeconds": 300
			}
		  ],
		  "priority": 0,
		  "enableServiceLinks": true
		},
		"status": {
		  "phase": "Running",
		  "conditions": [
			{
			  "type": "Initialized",
			  "status": "True",
			  "lastProbeTime": null,
			  "lastTransitionTime": "2020-01-23T08:18:20Z"
			},
			{
			  "type": "Ready",
			  "status": "True",
			  "lastProbeTime": null,
			  "lastTransitionTime": "2020-01-23T08:18:31Z"
			},
			{
			  "type": "ContainersReady",
			  "status": "True",
			  "lastProbeTime": null,
			  "lastTransitionTime": "2020-01-23T08:18:31Z"
			},
			{
			  "type": "PodScheduled",
			  "status": "True",
			  "lastProbeTime": null,
			  "lastTransitionTime": "2020-01-23T08:18:17Z"
			}
		  ],
		  "hostIP": "10.240.0.6",
		  "podIP": "10.244.0.86",
		  "startTime": "2020-01-23T08:18:17Z",
		  "initContainerStatuses": [],
		  "containerStatuses": [
			{
			  "name": "TEST-CONTAINER-NAME",
			  "state": {
				"running": {
				  "startedAt": "2020-01-23T08:18:23Z"
				}
			  },
			  "lastState": {},
			  "ready": true,
			  "restartCount": 0,
			  "image": "registry.docker.io/images/test-application:ddbdd5056b530d6dac910f2ab49e219fcaf46dae",
			  "imageID": "docker-pullable://registry.docker.io/images/test-application@sha256:16c6180ebe5e7338a541c8da9fd0e0573b340ce2e041fe6026654893516c912b",
			  "containerID": "docker://TEST-CONTAINER-ID"
			}
		  ],
		  "qosClass": "BestEffort"
		}
	  }
	]
  }`
