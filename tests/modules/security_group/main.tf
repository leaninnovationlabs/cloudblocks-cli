locals {
  default_cidr_blocks = []
}
resource "aws_security_group" "aws_sec_group" {
  name        = var.sg_name
  description = var.sg_description
  vpc_id      = var.default_vpc_id
  dynamic "ingress" {
    for_each = var.sg_ingress_rules
    content {
      from_port       = ingress.value.from_port
      to_port         = ingress.value.to_port
      protocol        = ingress.value.protocol
      cidr_blocks     = lookup(ingress.value, "cidr_blocks", null)
      security_groups = lookup(ingress.value, "security_groups", null)
      description     = ingress.value.description
    }
  }
  dynamic "egress" {
    for_each = var.sg_egress_rules
    content {
      from_port       = egress.value.from_port
      to_port         = egress.value.to_port
      protocol        = egress.value.protocol
      cidr_blocks     = lookup(egress.value, "cidr_blocks", null)
      security_groups = lookup(egress.value, "security_groups", null)
      description     = egress.value.description
    }
  }
  tags = var.sg_tags
}