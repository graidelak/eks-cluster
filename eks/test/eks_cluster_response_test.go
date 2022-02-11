package test

import (
	"crypto/tls"
	"strings"
	"testing"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

// checkClusterRespondsImmediately checks if the eks is returning something
func checkClusterRespondsImmediately(t *testing.T) {
	tfDir := "../"

	// Load Tf Options
	terraformOptions := test_structure.LoadTerraformOptions(t, tfDir)

	// Grab Tf output
	clusterURL := terraform.Output(t, terraformOptions, "cluster_endpoint")
	clusterURL = strings.Replace(clusterURL, "\"", "", 2)

	t.Log("Starting test checkClusterRespondsImmediately")

	// Expected status and body
	status := 403
	body := `{"kind":"Status","apiVersion":"v1","metadata":{},"status":"Failure","message":"forbidden: User \"system:anonymous\" cannot get path \"/\"","reason":"Forbidden","details":{},"code":403}`
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	http_helper.HttpGetWithValidation(t, clusterURL, tlsConfig, status, body)

	t.Log("Ending test checkClusterRespondsImmediately")
}
