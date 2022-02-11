package test

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

type EksNodeGroupOutput struct {
	AmiType        string                   `json:"ami_type"`
	Arn            string                   `json:"arn"`
	ClusterName    string                   `json:"cluster_name"`
	DiskSize       int                      `json:"disk_size"`
	ID             string                   `json:"id"`
	InstanceTypes  []string                 `json:"instance_types"`
	Labels         map[string]interface{}   `json:"labels"`
	NodeGroupName  string                   `json:"node_group_name"`
	NodeRoleArn    string                   `json:"node_role_arn"`
	ReleaseVersion string                   `json:"release_version"`
	RemoteAccess   []string                 `json:"remote_access"`
	Resources      map[string]interface{}   `json:"resources"`
	ScalingConfig  []map[string]interface{} `json:"scaling_config"`
	Status         string                   `json:"status"`
	SubnetIds      []string                 `json:"subnet_ids"`
	Tags           map[string]interface{}   `json:"tags"`
	Version        string                   `json:"version"`
}

func TestMain(t *testing.T) {

	// The folder where we have our Terraform code
	tfDir := "../"

	// Sets up a Unique ID for this test
	uniqueID := strings.ToLower(random.UniqueId())
	reg, err := regexp.Compile("[^a-zA-Z]+")
	if err != nil {
		log.Fatal(err)
	}
	uniqueID = reg.ReplaceAllString(uniqueID, "")
	test_structure.SaveString(t, tfDir, "UniqueID", uniqueID)

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	// awsRegion := aws.GetRandomStableRegion(t, nil, nil)
	awsRegion := "us-west-1"
	test_structure.SaveString(t, tfDir, "AwsRegion", awsRegion)

	// At the end of the test, destroy the AMI Created and Tf
	defer test_structure.RunTestStage(t, "cleanup", func() {
		// Clean up Tf
		destroyTerraform(t, tfDir)
	})

	// Stage to Deploy Terraform
	test_structure.RunTestStage(t, "terraform-create", func() {
		deployUsingTerraform(t, tfDir, uniqueID, awsRegion)
	})

	// Enable testing
	t.Parallel()

	// Test Monitoring
	test_structure.RunTestStage(t, "testing", func() {
		checkClusterRespondsImmediately(t)
		checkClusterNodeGroupIsActive(t)
		// checkRoleMatchesPolicies(t)
	})

}

// Undeploy the app using Terraform
func destroyTerraform(t *testing.T, tfDir string) {
	// Load the Terraform Options saved by the earlier deploy_terraform stage
	terraformOptions := test_structure.LoadTerraformOptions(t, tfDir)

	// Destroy Tf
	terraform.Destroy(t, terraformOptions)
}

// Deploy the app using Terraform
func deployUsingTerraform(t *testing.T, tfDir string, uniqueID string, awsRegion string) {

	// Get the Default VPC
	defaultVpc := aws.GetDefaultVpc(t, awsRegion)
	subnetsList := []string{defaultVpc.Subnets[0].Id, defaultVpc.Subnets[1].Id}

	// Define Cluster Name for Testing in String for test
	clusterName := fmt.Sprintf("eks-%s", uniqueID)
	test_structure.SaveString(t, tfDir, "ClusterName", clusterName)

	// Define Worker Groups Config
	workerGroups := []map[string]interface{}{
		{
			"name":                 fmt.Sprintf("node-%s", uniqueID),
			"instance_type":        "t3.large",
			"asg_desired_capacity": 1,
		},
	}

	// Define Node Groups Config
	nodeGroups := map[string]interface{}{
		fmt.Sprintf("%s-node-group", clusterName): map[string]interface{}{
			"desired_capacity": 1,
			"min_capacity":     1,
			"max_capacity":     2,
		},
	}

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: tfDir,

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"region":                   awsRegion,
			"cluster_name":             clusterName,
			"environment":              "test",
			"vpc_id":                   defaultVpc.Id,
			"private_subnets":          subnetsList,
			"node_group_instance_type": "t2.micro",
			"node_group_ami_type":      "AL2_x86_64",
			"node_group_disk_size":     6,
			"cluster_version":          "1.18",
			"worker_groups":            workerGroups,
			"node_groups":              nodeGroups,
		},

		// // Variables to pass to our Terraform code using -var-file options // In case you want to use in the future :)
		// VarFiles: []string{"./test/terratest.tfvars"},
	}

	// Save the Terraform Options struct text so future test stages can use it
	test_structure.SaveTerraformOptions(t, tfDir, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)
}
