output "task_definition_arn" {
  value = var.fluent_bit_image == null ? aws_ecs_task_definition.fargate_taskdefinition[0].arn : null
}

output "task_definition_with_fluentbit_arn" {
  value = var.fluent_bit_image == null ? null : aws_ecs_task_definition.fargate_taskdefinition_with_fluentbit[0].arn
}

output "td_web_container_name" {
  value = var.app_container_name
}