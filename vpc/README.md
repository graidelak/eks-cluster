# Simple VPC

Configuration in this directory creates set of VPC resources.

There is a public and private subnet created per availability zone in addition to single NAT Gateway shared between all 3 availability zones.


## Environment

This module expects you to set the AWS_PROFILE environment variable, in order to read the correct AWS  credentials to operate with.

```bash
export AWS_PROFILE=<your aws profile name>
```


## Pre-requisites
- terraform 1.0.0


## Outputs

| Name | Description |
|------|-------------|
| nat_public_ips | NAT gateways |
| private_subnets | Subnets |
| public_subnets | List of IDs of public subnets |
| vpc_cidr_block | CIDR blocks |
| vpc_id | VPC |



## Reference

If you want to have more details information about the module, check the follow links:

https://github.com/gruntwork-io/terragrunt


https://github.com/terraform-aws-modules/terraform-aws-vpc

https://registry.terraform.io/modules/terraform-aws-modules/vpc/aws/1.46.0
