variable "cidr" {
  description = "CIDR block"
  default     = "10.0.0.0/16"
}

variable "region" {
  description = "AWS region to deploy to"
  type        = string
  default     = "us-east-1"
}

variable "azs_count" {
  description = "Number of AZs to use."
  type        = number
  default     = 2
}

variable "private_subnet" {
  default = ["10.0.1.0/24", "10.0.2.0/24"]
  type    = list(string)
}

variable "public_subnet" {
  default = ["10.0.101.0/24", "10.0.102.0/24"]
  type    = list(string)
}


variable "environment" {
  default = "dev"
}


variable "extra_tags" {
  description = "Map of extra tags to append to the resulting VPC"
  type        = map(string)
  default     = {}
}

variable "public_subnet_extra_tags" {
  description = "Map of extra tags to append to the resulting public subnet"
  type        = map(string)
  default     = {"kubernetes.io/role/elb" = 1}
}

variable "private_subnet_extra_tags" {
  description = "Map of extra tags to append to the resulting private subnet"
  type        = map(string)
  default     = {"kubernetes.io/role/internal-elb" = 1}
}

variable "eks_cluster_name" {
  description = "EKS cluster name"
  default = "eks-cluster"
  type = string
}

## EKS variables

variable "owner" {
  type    = string
  default = ""
}

variable "business_unit" {
  type    = string
  default = ""
}

variable "cluster_version" {
  type        = string
  description = "Version of the EKS Cluster to be deployed"
  default     = "1.18"
}

variable "node_group_instance_type" {
  type        = string
  description = "Instance type for the node group's instances."
  default     = "m5.2xlarge"
}

variable "node_group_disk_size" {
  type        = string
  description = "Size in GB to allocate for the node group's instances."
  default     = "50"
}

variable "node_group_ami_type" {
  type        = string
  description = "AMI Identifier for the node group's instances."
  default     = "AL2_x86_64"
}

variable "enable_irsa" {
  type        = bool
  description = "Enable IAM Roles for EKS Service-Accounts (IRSA)."
  default     = false
}

variable "worker_groups" {
  type        = any
  description = "List of group definitions for worker nodes."
  default     = []
}

variable "node_groups" {
  type        = map(any)
  description = "Definitios for worker nodes, with a worker group's name as its key."
  default     = {}
}

variable "map_users" {
  type = list(object({
    userarn  = string
    username = string
    groups   = list(string)
  }))
  description = "Additional IAM users to add to the aws-auth configmap."
  default     = []
}

variable "map_roles" {
  description = "Additional IAM roles to add to the aws-auth configmap. See examples/basic/variables.tf for example format."
  type = list(object({
    rolearn  = string
    username = string
    groups   = list(string)
  }))
  default = []
}
variable "workers_role_name" {
  description = "User defined workers role name."
  type        = string
  default     = ""
}

variable "route53_domain" {
  type    = string
  default = ""
}

variable "certificate_arn" {
  type    = string
  default = ""
}
