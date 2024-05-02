output "aws_alb_created_arn" {
  value = aws_lb.app_alb.arn
}

output "alb_listener_arn" {
  value = aws_lb_listener.alb_listener.arn
}