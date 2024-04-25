variable "create_log_group" {
  type        = bool
  default     = true
  description = "Create the log group"
}

variable "log_group_name" {
  type        = string
  description = "Log Group Name"
}

variable "log_group_class" {
  type        = string
  default     = "STANDARD"
  description = "The log class of the log group"
}


variable "log_group_tags" {
  type        = map(any)
  description = "Log Group tags"
  default     = {}
}