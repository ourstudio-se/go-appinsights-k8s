package appink8s

import (
	"errors"
	"fmt"
)

const k8sContainerInfoPath = "/proc/self/cgroup"

type k8sinitializer struct {
	client *k8sclient
}

func newK8sInitializer(c *k8sclient) *k8sinitializer {
	return &k8sinitializer{
		client: c,
	}
}

func (ki *k8sinitializer) ReadPropertySpec() (*runtimeSpec, error) {
	containerID, err := ki.client.ReadContainerID()
	if err != nil {
		return nil, err
	}

	pods, err := ki.client.GetPods()
	if err != nil {
		return nil, err
	}

	nodes, err := ki.client.GetNodes()
	if err != nil {
		return nil, err
	}

	result := &runtimeSpec{
		ContainerID: containerID,
	}
	for _, pod := range pods.List {
		found := false
		for _, statuses := range pod.Status.ContainerStatuses {
			if statuses.ID == fmt.Sprintf("docker://%s", containerID) {
				result.ContainerName = statuses.Name
				found = true
			}
		}

		if found {
			node := nodes.FindByName(pod.RuntimeSpec.NodeName)

			result.NodeID = node.MetaData.ID
			result.NodeName = node.MetaData.Name
			result.NodeLabels = node.MetaData.GetLabels()
			result.PodID = pod.MetaData.ID
			result.PodName = pod.MetaData.Name
			result.PodLabels = pod.MetaData.GetLabels()
			result.ReplicaSetName = pod.FindReplicaSetName()
			result.DeploymentName = pod.FindDeploymentName()
			return result, nil
		}
	}

	return nil, errors.New("no runtime spec could be found")
}
