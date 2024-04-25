resource "aws_ecs_task_definition" "fargate_taskdefinition" {
  family                   = var.task_def_family
  task_role_arn            = var.task_role_arn
  execution_role_arn       = var.task_exec_role_arn
  network_mode             = "awsvpc"
  cpu                      = var.task_def_cpu
  memory                   = var.task_def_memory
  requires_compatibilities = ["FARGATE"]

  container_definitions = jsonencode([
    {
      "name" : "${var.app_container_name}",
      "image" : "${var.app_container_image}",
"portMappings" : [
    {
        "containerPort": var.port1,
        "hostPort": var.port1
    },
    {
        "containerPort": var.port2,
        "hostPort": var.port2
    }
],
      "environment" : "${var.env_variables}",
      "healthCheck" : "${var.cont_health_check_config}"
      "logConfiguration" : {
        "logDriver" : "awslogs",
        "options" : {
          "awslogs-group" : "${var.log_group}",
          "awslogs-region" : "${var.aws_region}",
          "awslogs-stream-prefix" : "${var.log_stream_prefix}"
        }
      }
    }
  ])
  tags = var.task_def_tags
}
