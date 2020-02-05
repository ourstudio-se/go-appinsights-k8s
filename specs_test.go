package appink8s

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_That_MetaData_GetLabels_Returns_Labels_Map_As_String(t *testing.T) {
	labels := make(map[string]string)
	labels["prop1"] = "value1"
	labels["prop2"] = "value2"

	spec := &metaDataSpec{
		Labels: labels,
	}

	var list []string
	for k, v := range labels {
		list = append(list, fmt.Sprintf("%s:%s", k, v))
	}
	sort.Strings(list)
	expected := strings.Join(list, ",")

	assert.Equal(t, expected, spec.GetLabels())
}

func Test_That_PodSpec_FindDeploymentName_Returns_Deployment_Owner_Case_Insensitive(t *testing.T) {
	deployment := podOwnerSpec{
		Name: "deployment-1",
		Kind: "DePlOyMeNt",
	}
	spec := &podSpec{
		MetaData: metaDataSpec{
			Owners: []podOwnerSpec{
				podOwnerSpec{
					Name: "replicaset-1",
					Kind: "replicaset",
				},
				deployment,
				podOwnerSpec{
					Name: "other-1",
					Kind: "other",
				},
			},
		},
	}

	result := spec.FindDeploymentName()
	assert.Equal(t, deployment.Name, result)
}

func Test_That_PodSpec_FindReplicaSetName_Returns_ReplicaSet_Owner_Case_Insensitive(t *testing.T) {
	replicaset := podOwnerSpec{
		Name: "replicaset-1",
		Kind: "RePlIcAsEt",
	}
	spec := &podSpec{
		MetaData: metaDataSpec{
			Owners: []podOwnerSpec{
				podOwnerSpec{
					Name: "deployment-1",
					Kind: "deployment",
				},
				replicaset,
				podOwnerSpec{
					Name: "other-1",
					Kind: "other",
				},
			},
		},
	}

	result := spec.FindReplicaSetName()
	assert.Equal(t, replicaset.Name, result)
}

func Test_That_PodSpec_FindDeploymentName_Returns_ReplicaSet_Prefix_When_Missing_Deployment_Owner(t *testing.T) {
	deploymentName := "deployment"
	replicaset := podOwnerSpec{
		Name: fmt.Sprintf("%s-1", deploymentName),
		Kind: "RePlIcAsEt",
	}
	spec := &podSpec{
		MetaData: metaDataSpec{
			Owners: []podOwnerSpec{
				podOwnerSpec{
					Name: "non-deployment-1",
					Kind: "non-deployment",
				},
				replicaset,
				podOwnerSpec{
					Name: "other-1",
					Kind: "other",
				},
			},
		},
	}

	result := spec.FindDeploymentName()
	assert.Equal(t, deploymentName, result)
}

func Test_That_NodeListSpec_FindByName_Returns_Matching_NodeSpec(t *testing.T) {
	node := nodeSpec{
		MetaData: metaDataSpec{
			ID:   "node-id",
			Name: "node-name",
		},
	}
	spec := nodeListSpec{
		List: []nodeSpec{
			nodeSpec{
				MetaData: metaDataSpec{
					ID:   "1",
					Name: "1",
				},
			},
			node,
			nodeSpec{
				MetaData: metaDataSpec{
					ID:   "2",
					Name: "2",
				},
			},
		},
	}

	result := spec.FindByName(node.MetaData.Name)
	assert.Equal(t, node, result)
}
