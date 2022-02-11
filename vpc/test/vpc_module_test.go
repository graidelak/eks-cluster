package test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestVpcModule(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		// Terraform path
		TerraformDir: "../",
		Vars: map[string]interface{}{
			"region":        "us-west-2",
			"business-unit": "qwiretestrun",
			"owner":         "flugeltestrun",
		},
	}
	// At the end run `terraform destroy`
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`
	terraform.InitAndApply(t, terraformOptions)
	actualVpcID := terraform.Output(t, terraformOptions, "vpc_id")

	// Extract the desired information via outputs
	actualPrivateSubnet := terraform.OutputList(t, terraformOptions, "private_subnets")
	actualPublicSubnet := terraform.OutputList(t, terraformOptions, "public_subnets")
	actualRegion := terraform.Output(t, terraformOptions, "region")
	actualAzs := terraform.OutputList(t, terraformOptions, "azs")

	fmt.Println(actualPrivateSubnet)
	fmt.Println(actualPublicSubnet)
	fmt.Println(actualAzs)

	// Get the expected subnets using the aws module
	expectedSubnets := aws.GetSubnetsForVpc(t, actualVpcID, actualRegion)
	public_subnet := aws.Subnet{Id: actualPublicSubnet[0], AvailabilityZone: actualAzs[0]}
	private_subnet := aws.Subnet{Id: actualPrivateSubnet[0], AvailabilityZone: actualAzs[0]}

	// check if expected and actual subnets match.
	assert.Contains(t, expectedSubnets, public_subnet)
	assert.Contains(t, expectedSubnets, private_subnet)
}
