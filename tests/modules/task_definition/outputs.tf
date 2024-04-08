output "task_definition_arn" {
  value = aws_ecs_task_definition.fargate_taskdefinition.arn
}

output "td_web_container_name" {
  value = var.app_container_name
}