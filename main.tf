terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = ">= 3.43.0"
  }
}

data "aws_availability_zones" "available" {
  state = "available"
}

data "aws_eks_cluster" "cluster" {
  name = module.eks.cluster_id
}

data "aws_eks_cluster_auth" "cluster" {
  name = module.eks.cluster_id
}

data "aws_iam_role" "kubernetes_worker_node" {
  name = module.eks.worker_iam_role_name
}

locals {
  resource_name = var.environment
  eks_cluster_name = var.eks_cluster_name
}

module "vpc" {
  source = "./vpc/modules/terraform-aws-vpc"

  name = "${local.resource_name}-vpc"

  cidr = var.cidr

  azs             = slice(data.aws_availability_zones.available.names, 0, var.azs_count)
  private_subnets = var.private_subnet
  public_subnets  = var.public_subnet

  tags = merge(
    {
      "Environment" = var.environment
    },
    {
      "kubernetes.io/cluster/${local.eks_cluster_name}" = "shared"
    }
  )

  vpc_tags = {
    Name = local.resource_name
  }

  public_subnet_tags = merge(
    {
      "Name" = "${local.resource_name}-public"
    },
    var.public_subnet_extra_tags,
  )

  private_subnet_tags = merge(
    {
      "Name" = "${local.resource_name}-private"
    },
    var.private_subnet_extra_tags,
  )
}

module "eks" {
  source = "./eks/modules/terraform-aws-eks"

  cluster_name    = local.eks_cluster_name
  cluster_version = var.cluster_version
  subnets         = module.vpc.private_subnets
  vpc_id          = module.vpc.vpc_id
  enable_irsa     = true
  region          = var.region
  tags = {
    Name        = local.eks_cluster_name
    Environment = var.environment
  }

  worker_groups = var.worker_groups

  node_groups_defaults = {
    ami_type      = var.node_group_ami_type
    disk_size     = var.node_group_disk_size
    instance_type = var.node_group_instance_type
    k8s_labels = {
      Environment = var.environment
    }
    tags = {
      Name        = local.eks_cluster_name
      Environment = var.environment
    }
  }

  node_groups = var.node_groups

  map_users = var.map_users

  map_roles = var.map_roles
}
