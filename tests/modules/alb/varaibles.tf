variable "default_vpc_id" {
  type        = string
  description = "Default vpc id"
}

variable "subnet_ids" {
  type        = list(string)
  description = "Subnet ids"
}

variable "alb_name" {
  type        = string
  description = "Application Load Balancer Name"
}


variable "alb_sg_name" {
  type        = string
  description = "Security Group for Application Load Balancer"
}

variable "sg_ingress_rules" {
  description = "Ingress security group rules"
  type        = list(any)
  default = [{ from_port = 8443, to_port = 8443, protocol = "tcp", cidr_blocks = ["0.0.0.0/0"], description = "Allow all inbound from 8443" },
  { from_port = 8080, to_port = 8080, protocol = "tcp", cidr_blocks = ["0.0.0.0/0"], description = "Allow all inbound from 8080" }]
}

variable "sg_egress_rules" {
  description = "Egress security group rules"
  type        = list(any)
  default     = [{ from_port = 1, to_port = 65535, protocol = "tcp", cidr_blocks = ["0.0.0.0/0"], description = "Outbound Rules for the tcp range 1, 65535" }]
}

variable "ssl_policy" {
  type        = string
  description = "ALB Listener Rule's SSL Poilcy"
}

variable "cert_arn" {
  type        = string
  description = "ALB Listener Rule's certificate ARN"
}

variable "alb_tags" {
  type        = map(any)
  description = "ALB Tags"
  default     = {}
}
