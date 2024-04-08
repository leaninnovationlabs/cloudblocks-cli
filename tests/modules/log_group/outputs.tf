output "log_group_created" {
  value = var.create_log_group
}

output "log_group_created_arn" {
  value = one(aws_cloudwatch_log_group.log_group[*].arn)
}

output "log_group_name" {
  value = var.log_group_name
}