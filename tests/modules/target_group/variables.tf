variable "target_group_name" {
  type        = string
  description = "Target Group Name"
}

variable "vpc_id" {
  type        = string
  description = "VPC ID for the Target Group"
}

variable "target_group_type" {
  type        = string
  description = "Target Group Type"
  default     = "ip"
}

variable "target_group_protocol" {
  type        = string
  description = "Target Group protocol"
  default     = "HTTPS"
}

variable "tg_protocol_version" {
  type        = string
  description = "Target Group protocol"
  default     = "HTTP1"
}

variable "target_group_port" {
  type        = number
  description = "Target Group port"
  default     = 8443
}

variable "tg_health_check_path" {
  type        = string
  description = "Target Group Health Check Path"
}

variable "tg_health_check_protocol" {
  type        = string
  description = "Target Group Health Check Path"
  default     = "HTTPS"
}


variable "tg_health_check_port" {
  type        = string
  description = "Target Group Health Check Port"
  default     = "traffic-port"
}

variable "tg_healthy_threshold" {
  type        = number
  description = "Target Group Health Check Healthy Threshold"
  default     = 5
}

variable "tg_unhealthy_threshold" {
  type        = number
  description = "Target Group Health Check Unhealthy Threshold"
  default     = 2
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

variable "tg_health_matcher" {
  type        = string
  description = "Target Group Health Check matcher codes"
  default     = 200
}

variable "target_group_tags" {
  type        = map(any)
  description = "Target Group tags"
  default     = {}
}