module "alb_sg" {
  source           = "../security_group"
  create_sg        = true
  sg_name          = var.alb_sg_name
  sg_description   = "Security Group for ${var.alb_sg_name}"
  default_vpc_id   = var.default_vpc_id
  sg_ingress_rules = var.sg_ingress_rules
  sg_egress_rules  = var.sg_egress_rules
  sg_tags          = { Name = "${var.alb_sg_name}" }
}

resource "aws_lb" "app_alb" {
  //for_each           = { for s in var.alb_properties : s.name => s }
  internal           = true
  name               = var.alb_name
  load_balancer_type = "application"
  ip_address_type    = "ipv4"
  security_groups    = [module.alb_sg.security_group_id]
  subnets            = var.subnet_ids
  tags               = var.alb_tags
}

resource "aws_lb_listener" "alb_listener" {
  //for_each          = { for s in var.alb_properties : s.name => s }
  load_balancer_arn = aws_lb.app_alb.arn
  ssl_policy        = var.ssl_policy
  certificate_arn   = var.cert_arn
  port              = 443
  protocol          = "HTTPS"
  default_action {
    type = "fixed-response"
    fixed_response {
      content_type = "text/plain"
      status_code  = 503
      message_body = "You shouldn't reach this page, something is configured wrong."
    }
  }

}
