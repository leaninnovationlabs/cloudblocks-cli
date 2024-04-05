module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 20.0"

  cluster_name                   = $CLUSTER_NAME
  cluster_version                = $CLUSTER_VERSION
  cluster_endpoint_public_access = $CLUSTER_PUBLIC_ACCESS

  cluster_addons = {
    coredns = {
      most_recent = true
    }
    kube-proxy = {
      most_recent = true
    }
    vpc-cni = {
      most_recent = true
    }
  }

  vpc_id     = $VPC_ID
  subnet_ids = $SUBNET_IDS

  control_plane_subnet_ids = $CONTROL_PLANE_SUBNETS

  # EKS Managed Node Group(s)
  eks_managed_node_group_defaults = {
    instance_types = [$INSTANCE_TYPE]
  }

  eks_managed_node_groups = {
    example = {
      min_size     = $MIN_SIZE
      max_size     = $MAX_SIZE
      desired_size = $DESIRED_SIZE
      instance_types = [$INSTANCE_TYPE]
      capacity_type  = $CAPACITY_TYPE
    }
  }

  # Cluster access entry
  # To add the current caller identity as an administrator
  enable_cluster_creator_admin_permissions = $ENABLE_CREATOR_ADMIN_PERM
  tags = $TAGS
}

