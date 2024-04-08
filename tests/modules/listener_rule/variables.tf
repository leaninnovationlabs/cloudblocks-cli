variable "listener_arn" {
  type        = string
  description = "Listener Rule ARN"
}

variable "target_group_arn" {
  type        = string
  description = "Target Group ARN"
}

variable "action_type" {
  type        = string
  default     = "forward"
  description = "Listener Action type"
}

variable "listener_conditions" {
  type        = list(any)
  description = "Listener Conditions"
}

