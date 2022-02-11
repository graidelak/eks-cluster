# EKS and VPC Module

This module stands up an EKS cluster, with its own VPC, a workers
 group, and a node group cluster

## Requirements

* Terraform >= 1.0.0

## Worker nodes

Definition of worker groups and the nodes that use them is fully handled via variables. A detailed definiton of the available values for `worker_groups` can be found [here](https://github.com/terraform-aws-modules/terraform-aws-eks/blob/master/local.tf#L26), and for `node_groups` [here](https://github.com/terraform-aws-modules/terraform-aws-eks/blob/master/modules/node_groups/README.md)

Example:

```
worker_groups = [
  {
    name                 = "test-eks-worker-group-1"
    instance_type        = "t2.small"
    additional_user_data = "echo hello world"
    asg_desired_capacity = 1
  }
]

node_groups = {
  test-cluster-node-group = {
    desired_capacity = 1
    min_capacity     = 1
    max_capacity     = 2
  }
}

```

## VPC

The module definition expects a VPC ID and its private subnets as input
 variables. It's important to note that the VPC that is to be used with
 this cluster has to be tagged with some extra considerations:

* The VPC and all its resources should be tagged with
 `kubernetes.io/cluster/<cluster name>: shared`
* The VPC's public subnets should be tagged with
 `kubernetes.io/role/elb: 1`
* The VPC's private subnets should be tagged with
 `kubernetes.io/role/internal-elb: 1`

## How to 

### PreRequirements

* Terraform v1.0 https://learn.hashicorp.com/tutorials/terraform/install-cli
* AWS cli https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html
* An AWS account with admin privileges

    Create access key



  
## Steps

First we need to configure our AWS credencials for that , we need to run:

```bash
$ aws configure
AWS Access Key ID [None]: xxxxxxxxxxxxxxxxx
AWS Secret Access Key [None]: xxxxxxxxxxxxxx
Default region name [None]: us-west-1
```

To validate your aws user , run this command:

```bash
aws sts get-caller-identity
```

Once your AWS profile/user is configure , the next step is to install the EKS cluster.

To install the EKS we need to initialize Terraform by running `terraform init`. Terraform will generate a directory named `.terraform` and download each module source declared in `main.tf`. Initialization will pull in any providers required by these modules.

```bash
$ terraform init
Initializing modules...
- eks in eks/modules/terraform-aws-eks
- eks.node_groups in eks/modules/terraform-aws-eks/modules/node_groups
- vpc in vpc/modules/terraform-aws-vpc

Initializing the backend...

Initializing provider plugins...
- Reusing previous version of terraform-aws-modules/http from the dependency lock file
- Reusing previous version of hashicorp/local from the dependency lock file
- Reusing previous version of hashicorp/cloudinit from the dependency lock file
- Reusing previous version of hashicorp/aws from the dependency lock file
- Reusing previous version of hashicorp/kubernetes from the dependency lock file
- Installing hashicorp/cloudinit v2.2.0...
- Installed hashicorp/cloudinit v2.2.0 (signed by HashiCorp)
- Installing hashicorp/aws v4.0.0...
- Installed hashicorp/aws v4.0.0 (signed by HashiCorp)
- Installing hashicorp/kubernetes v2.8.0...
- Installed hashicorp/kubernetes v2.8.0 (signed by HashiCorp)
- Installing terraform-aws-modules/http v2.4.1...
- Installed terraform-aws-modules/http v2.4.1 (self-signed, key ID B2C1C0641B6B0EB7)
- Installing hashicorp/local v2.1.0...
- Installed hashicorp/local v2.1.0 (signed by HashiCorp)

Partner and community providers are signed by their developers.
If you'd like to know more about provider signing, you can read about it here:
https://www.terraform.io/docs/cli/plugins/signing.html

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```

After Terraform has been successfully initialized, run `terraform plan` to review what will be created:

After the plan is validated, apply the changes by running `terraform apply`. For one last validation step, Terraform will output the plan again and prompt for confirmation before applying. This step will take around 15-20 minutes to complete.

To connect to your cluster, run this command:

```bash
aws eks --region us-east-1 update-kubeconfig --name <eksclustername>
```

Next, run `kubectl get no` and you will see worker nodes from your cluster.

To destroy the cluster, run:

```bash
terraform destroy
```