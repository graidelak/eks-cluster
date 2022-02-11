package test

import (
	"encoding/json"
	"fmt"
	"testing"

	terraform "github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

// checkClusterNodeGroupIsActive checks if the node group is enabled in the cluster
func checkClusterNodeGroupIsActive(t *testing.T) {
	tfDir := "../"

	// Load Tf Options
	terraformOptions := test_structure.LoadTerraformOptions(t, tfDir)

	// Load Eks Node Group Struct
	nodeOutput := EksNodeGroupOutput{}

	// Grab Tf Output
	eksOutput := terraform.OutputForKeys(t, terraformOptions, []string{"node_groups"})

	// Grab Cluster Name
	clusterName := test_structure.LoadString(t, tfDir, "ClusterName")

	// Serialize data
	nodeData, _ := json.Marshal(eksOutput["node_groups"].(map[string]interface{})[fmt.Sprintf("%s-node-group", clusterName)])
	json.Unmarshal([]byte(nodeData), &nodeOutput)

	assert.Equal(t, "ACTIVE", nodeOutput.Status)
}
