variable "aws_region" {
  type        = string
  description = "AWS region"
  default     = "us-gov-west-1"
}

variable "default_vpc_id" {
  type        = string
  description = "Default vpc id"
}

variable "subnet_ids" {
  type        = list(string)
  description = "The subnet to launch the ECS into"
}

//LOG GROUP
variable "log_group_name" {
  type        = string
  description = "Log Group Name"
}

variable "log_group_tags" {
  type    = map(any)
  default = {}
}

//TASK DEFINITION
variable "task_def_family" {
  type        = string
  description = "Task Definition Name"
}

variable "task_role_arn" {
  type        = string
  description = "Task Definition Role ARN"
}

variable "task_exec_role_arn" {
  type        = string
  description = "Task Definition Execution Role ARN"
}

variable "task_def_cpu" {
  type        = string
  description = "Task Definition CPU"
}

variable "task_def_memory" {
  type        = string
  description = "Task Definition Memory"
}

variable "container_definitions_path" {
  type        = string
  description = "Task Definition container defintion paths"
  default     = ""
}

variable "app_container_name" {
  type        = string
  description = "Task Definition Application Container Name"
}

variable "app_container_image" {
  type        = string
  description = "Application Container image"
}

variable "cont_health_check_config" {
  type        = object({ command = list(string), interval = number, timeout = number, retries = number })
  description = "Container Health Check Configurations"
  default     = null
}

variable "app_port_mappings" {
  type        = list(any)
  description = "Port Mappings"
}

variable "env_variables" {
  type        = list(any)
  description = "Environment Variables"
}

variable "log_router_name" {
  type        = string
  description = "Task Definition Log Container Router Name"
  default     = null
}

variable "fluent_bit_image" {
  type        = string
  description = "Log fluent bit image"
  default     = null
}

variable "log_stream_prefix" {
  type        = string
  description = "Log Stream prefix"
}

variable "mount_points" {
  type        = map(string)
  description = "Container Mount Points for volume"
  default     = null
}

variable "log_delivery_stream" {
  type        = string
  description = "Kinesis Delivery stream for log"
  default     = null
}

variable "volume_name" {
  type        = string
  description = "Log Stream prefix"
  default     = null
}

variable "efs_vol_config" {
  type    = map(any)
  default = null
}

variable "task_def_tags" {
  type        = map(any)
  description = "Task Definition tags"
  default     = {}
}

//TAREGET GROUP
variable "target_group_name" {
  type        = string
  description = "Target Group Name"
}

variable "tg_health_check_path" {
  type        = string
  description = "Target Group Health Check Path"
}

variable "tg_health_timeout" {
  type        = number
  description = "Target Group Health Check timout"
  default     = 50
}

variable "tg_health_interval" {
  type        = number
  description = "Target Group Health Check interval"
  default     = 59
}

variable "target_group_protocol" {
  type        = string
  description = "Target Group protocol"
  default     = "HTTPS"
}

variable "target_group_port" {
  type        = number
  description = "Target Group port"
  default     = 8443
}

variable "tg_health_check_port" {
  type        = string
  description = "Target Group Health Check Port"
  default     = "traffic-port"
}

variable "tg_health_check_protocol" {
  type        = string
  description = "Target Group Health Check Path"
  default     = "HTTPS"
}

variable "target_group_tags" {
  type        = map(any)
  description = "Target Group tags"
  default     = {}
}

//LISTENER RULES
variable "listener_arn" {
  type        = string
  description = "Listener Rule ARN"
}

variable "listener_conditions" {
  type        = list(any)
  description = "Listener Conditions"
}

//ECS SERVICE
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
}

variable "assign_public_ip" {
  type        = bool
  description = "Assign public id for ECS service"
  default     = false
}

variable "cluster_id" {
  type        = string
  description = "Cluster id"
}

variable "desired_count" {
  type        = number
  default     = 1
  description = "Desired count for ECS service"
}

variable "alb_security_group" {
  type        = string
  description = "Security group id for ALB"
}

variable "container_port" {
  type        = string
  description = "Please provide the contianer port"
  default     = 8443
}

variable "ecs_tags" {
  type    = map(any)
  default = {}
}

variable "efs_security_group" {
  type        = string
  description = "EFS security group id"
  default     = null
}
