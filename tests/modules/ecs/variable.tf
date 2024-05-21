variable "default_vpc_id" {
  type        = string
  description = "VPC id for the security group"
}

variable "ecs_service_name" {
  type        = string
  description = "ECS Fargate Service name"
}

variable "ecs_sg_name" {
  type        = string
  description = "ECS Fargate Service Security Group name"
}

variable "ecs_sg_ingress_rules" {
  description = "Ingress security group rules"
  type        = list(any)
}

variable "ecs_sg_egress_rules" {
  description = "Egress security group rules"
  type        = list(any)
  default     = [{ from_port = 1, to_port = 65535, protocol = "tcp", cidr_blocks = ["0.0.0.0/0"], description = "Outbound Rules for the tcp range 1, 65535" }]
}

variable "cluster_id" {
  type        = string
  description = "Cluster id"
}

variable "launch_type" {
  type        = string
  default     = "FARGATE"
  description = "Launch type for the ECS service"
}

variable "service_type" {
  type        = string
  default     = "REPLICA"
  description = "Service type for the ECS service"
}

variable "task_defn_arn" {
  type        = string
  description = "Task Definition ARN for ECS service"
}

variable "desired_count" {
  type        = number
  default     = 1
  description = "Desired count for ECS service"
}

variable "min_healthy_percent" {
  type        = number
  default     = 100
  description = "Minimum healthy Percentage for ECS service"
}

variable "max_healthy_percent" {
  type        = number
  default     = 200
  description = "Maximum healthy Percentage for ECS service"
}


variable "alb_security_group" {
  type        = string
  description = "Security group id for ALB"
}

variable "ecs_subnet_ids" {
  type        = list(any)
  description = "Subnet ids for ECS service"
}

variable "assign_public_ip" {
  type        = bool
  default     = false
  description = "Assign public id for ECS service"
}

variable "target_group_arn" {
  type        = string
  description = "Please provide the Target Group ARN"
}

variable "container_name" {
  type        = string
  description = "Please provide the contianer name from Task Defintion"
}

variable "container_port" {
  type        = number
  description = "Please provide the contianer port"
  default     = 8443
}

variable "efs_security_group" {
  type        = string
  description = "EFS security group id"
  default     = null
}

variable "ecs_tags" {
  type    = map(any)
  default = {}
}



