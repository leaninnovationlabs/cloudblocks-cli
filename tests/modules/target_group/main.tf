resource "aws_alb_target_group" "alb-target-group" {
  name             = var.target_group_name
  vpc_id           = var.vpc_id
  target_type      = var.target_group_type
  port             = var.target_group_port
  protocol         = var.target_group_protocol
  protocol_version = var.tg_protocol_version

  health_check {
    path                = var.tg_health_check_path
    protocol            = var.tg_health_check_protocol
    port                = var.tg_health_check_port
    healthy_threshold   = var.tg_healthy_threshold
    unhealthy_threshold = var.tg_unhealthy_threshold
    timeout             = var.tg_healthy_threshold
    interval            = var.tg_health_interval
    matcher             = var.tg_health_matcher
  }

  tags = var.target_group_tags
}