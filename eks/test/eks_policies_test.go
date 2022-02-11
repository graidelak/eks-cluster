package test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

// checkRoleMatchesPolicies checkes if the policies are expected
func checkRoleMatchesPolicies(t *testing.T) {
	tfDir := "../"

	// Load Tf Options
	terraformOptions := test_structure.LoadTerraformOptions(t, tfDir)

	// Serialize & Deserialize data
	nodeOutputMap := EksNodeGroupOutput{}
	nodeOutput := terraform.OutputForKeys(t, terraformOptions, []string{"node_groups"})
	nodeData, _ := json.Marshal(nodeOutput["node_groups"].(map[string]interface{})["default"])
	json.Unmarshal([]byte(nodeData), &nodeOutputMap)

	// Grab Role Arn
	nodeRoleArn := nodeOutputMap.NodeRoleArn

	// Extract the role's name from its ARN (last element when splitting by /)
	roleArnSlice := strings.Split(nodeRoleArn, "/")
	roleName := string(roleArnSlice[len(roleArnSlice)-1])
	input := &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(roleName),
	}

	// Expected Policies
	expectedPoliciesNames := []string{
		"AmazonEKSWorkerNodePolicy",
		"AmazonEKS_CNI_Policy",
		"AmazonEC2ContainerRegistryReadOnly",
		"AmazonS3ReadOnlyAccess",
	}

	// Use the AWS SDK to pull data from its API
	iamClient := iam.New(session.New())
	rolePoliciesResponse, err := iamClient.ListAttachedRolePolicies(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				t.Log(iam.ErrCodeNoSuchEntityException, aerr.Error())
			case iam.ErrCodeServiceFailureException:
				t.Log(iam.ErrCodeServiceFailureException, aerr.Error())
			default:
				t.Log(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			t.Log(err.Error())
		}
		return
	}

	policyNames := make([]string, 0, len(rolePoliciesResponse.AttachedPolicies))
	// The SDK returns a role's policies as an array of pointers
	// We manually reconstructing by dereferencing each of these pointers...
	for _, v := range rolePoliciesResponse.AttachedPolicies {
		policyNames = append(policyNames, *v.PolicyName)
	}

	// And then compare them to the expected policies
	assert.ElementsMatch(t, policyNames, expectedPoliciesNames)
}
