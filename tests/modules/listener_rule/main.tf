resource "aws_lb_listener_rule" "alb-listener-rule" {
  for_each     = { for key, value in var.listener_conditions : key => value if can(value.host_header) && can(value.path_pattern)}
  listener_arn = var.listener_arn
  action {
    type             = var.action_type
    target_group_arn = var.target_group_arn
  }

  condition {
    path_pattern {
      values = each.value.path_pattern
    }
  }

  condition {
    host_header {
      values = each.value.host_header
    }
  }
}

resource "aws_lb_listener_rule" "alb-listener-rule-with-path" {
  for_each     = { for key, value in var.listener_conditions : key => value if can(value.path_pattern)  && !can(value.host_header) }
  listener_arn = var.listener_arn
  action {
    type             = var.action_type
    target_group_arn = var.target_group_arn
  }

  condition {
    path_pattern {
      values = each.value.path_pattern
    }
  }
}

resource "aws_lb_listener_rule" "alb-listener-rule-with-host" {
  for_each     = { for key, value in var.listener_conditions : key => value if !can(value.path_pattern)  && can(value.host_header) }
  listener_arn = var.listener_arn
  action {
    type             = var.action_type
    target_group_arn = var.target_group_arn
  }

  condition {
    host_header {
      values = each.value.host_header
    }
  }
}




