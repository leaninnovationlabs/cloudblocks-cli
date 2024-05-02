variable "aws_region" {
  type        = string
  description = "AWS Region"
}

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

variable "log_group" {
  type        = string
  description = "Log Group"
}

variable "log_stream_prefix" {
  type        = string
  description = "Log Stream prefix"
}

variable "log_delivery_stream" {
  type        = string
  description = "Delivery stream for log"
  default     = null
}

variable "mount_points" {
  type        = map(string)
  description = "Container Mount Points for volume"
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