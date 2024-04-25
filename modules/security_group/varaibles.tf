variable "create_sg" {
  type        = bool
  description = "Create the security group"
  default     = false
}

variable "sg_name" {
  type        = string
  description = "Security Group Name"
}

variable "sg_description" {
  type        = string
  description = "Security Group Description"
  default     = "Security Group to control access"
}

variable "default_vpc_id" {
  type        = string
  description = "Default VPC id"
}

variable "sg_ingress_rules" {
  description = "Ingress security group rules"
  type        = list(any)
}

variable "sg_egress_rules" {
  description = "Egress security group rules"
  type        = list(any)
}

variable "sg_tags" {
  description = "Tags for Security Groups"
  type        = map(any)
  default     = {}
}