resource "aws_cloudwatch_log_group" "log_group" {
  count           = var.create_log_group ? 1 : 0
  name            = var.log_group_name
  log_group_class = var.log_group_class
  tags            = var.log_group_tags
}