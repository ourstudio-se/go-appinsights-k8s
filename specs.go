package appink8s

import (
	"fmt"
	"sort"
	"strings"
)

type podContainerStatusSpec struct {
	Name  string `json:"name"`
	Ready bool   `json:"ready"`
	ID    string `json:"containerID"`
}

type podStatusSpec struct {
	ContainerStatuses []podContainerStatusSpec `json:"containerStatuses"`
}

type podOwnerSpec struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
	ID   string `json:"uid"`
}

type metaDataSpec struct {
	Name   string            `json:"name"`
	ID     string            `json:"uid"`
	Labels map[string]string `json:"labels"`
	Owners []podOwnerSpec    `json:"ownerReferences"`
}

func (mds metaDataSpec) GetLabels() string {
	var labels []string

	for k, v := range mds.Labels {
		label := fmt.Sprintf("%s:%s", k, v)
		labels = append(labels, label)
	}

	sort.Strings(labels)
	return strings.Join(labels, ",")
}

type podNodeSpec struct {
	NodeName string `json:"nodeName"`
}

type podSpec struct {
	MetaData    metaDataSpec  `json:"metadata"`
	RuntimeSpec podNodeSpec   `json:"spec"`
	Status      podStatusSpec `json:"status"`
}

func (ps podSpec) FindReplicaSetName() string {
	for _, owner := range ps.MetaData.Owners {
		if strings.ToLower(owner.Kind) == "replicaset" {
			return owner.Name
		}
	}

	return ""
}

func (ps podSpec) FindDeploymentName() string {
	deploymentName := ""

	for _, owner := range ps.MetaData.Owners {
		if strings.ToLower(owner.Kind) == "deployment" {
			return owner.Name
		}
		if strings.ToLower(owner.Kind) == "replicaset" {
			if deploymentName == "" {
				p := strings.Split(owner.Name, "-")
				deploymentName = strings.Join(p[:len(p)-1], "-")
			}
		}
	}

	return deploymentName
}

type podListSpec struct {
	List []podSpec `json:"items"`
}

type nodeSpec struct {
	MetaData metaDataSpec `json:"metadata"`
}

type nodeListSpec struct {
	List []nodeSpec `json:"items"`
}

func (nls nodeListSpec) FindByName(name string) nodeSpec {
	for _, n := range nls.List {
		if n.MetaData.Name == name {
			return n
		}
	}

	return nodeSpec{
		MetaData: metaDataSpec{
			Name: name,
		},
	}
}

type runtimeSpec struct {
	PodID          string
	PodName        string
	PodLabels      string
	ReplicaSetName string
	DeploymentName string
	NodeID         string
	NodeName       string
	NodeLabels     string
	ContainerID    string
	ContainerName  string
}

func (r *runtimeSpec) ToPropertyMap() map[string]string {
	props := make(map[string]string)

	props["Kubernetes.Pod.ID"] = r.PodID
	props["Kubernetes.Pod.Name"] = r.PodName
	props["Kubernetes.Pod.Labels"] = r.PodLabels
	props["Kubernetes.ReplicaSet.Name"] = r.ReplicaSetName
	props["Kubernetes.Deployment.Name"] = r.DeploymentName
	props["Kubernetes.Container.ID"] = r.ContainerID
	props["Kubernetes.Container.Name"] = r.ContainerName
	props["Kubernetes.Node.ID"] = r.NodeID
	props["Kubernetes.Node.Name"] = r.NodeName
	props["Kubernetes.Node.Labels"] = r.NodeLabels

	return props
}
