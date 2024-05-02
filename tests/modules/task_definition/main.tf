resource "aws_ecs_task_definition" "fargate_taskdefinition" {
  count                    = var.fluent_bit_image == null ? 1 : 0
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
      "portMappings" : "${var.app_port_mappings}",
      "environment" : "${var.env_variables}",
      "healthCheck" : "${var.cont_health_check_config}",
      "mountPoints" : ["${var.mount_points}"],
      "logConfiguration" : {
        "logDriver" : "awslogs",
        "options" : {
          "awslogs-group" : "${var.log_group}",
          "awslogs-region" : "${var.aws_region}",
          "awslogs-stream-prefix" : "${var.log_stream_prefix}",
          "awslogs-create-group" : "true"
        }
      }
    }
  ])

  dynamic "volume" {
    for_each = var.volume_name == null ? [] : [1]
    content {
      name = var.volume_name
      dynamic "efs_volume_configuration" {
        for_each = var.efs_vol_config == null ? [] : [1]
        content {
          file_system_id     = var.efs_vol_config.file_system_id
          root_directory     = var.efs_vol_config.root_directory
          transit_encryption = var.efs_vol_config.transit_encryption
        }
      }
    }
  }
  tags = var.task_def_tags
}

resource "aws_ecs_task_definition" "fargate_taskdefinition_with_fluentbit" {
  count                    = var.fluent_bit_image == null ? 0 : 1
  family                   = var.task_def_family
  task_role_arn            = var.task_role_arn
  execution_role_arn       = var.task_exec_role_arn
  network_mode             = "awsvpc"
  cpu                      = var.task_def_cpu
  memory                   = var.task_def_memory
  requires_compatibilities = ["FARGATE"]

  container_definitions = jsonencode([
    {
      "name" : "${var.log_router_name}",
      "image" : "${var.fluent_bit_image}",
      "logConfiguration" : {
        "logDriver" : "awslogs",
        "options" : {
          "awslogs-group" : "${var.log_group}",
          "awslogs-region" : "us-gov-west-1",
          "awslogs-stream-prefix" : "${var.log_stream_prefix}",
          "awslogs-create-group" : "true"
        }
      },
      "firelensConfiguration" : {
        "type" : "fluentbit"
      }
    },
    {
      "name" : "${var.app_container_name}",
      "image" : "${var.app_container_image}",
      "portMappings" : "${var.app_port_mappings}",
      "environment" : "${var.env_variables}",
      "healthCheck" : "${var.cont_health_check_config}"
      "mountPoints" : ["${var.mount_points}"],
      "logConfiguration" : {
        "logDriver" : "awsfirelens",
        "options" : {
          "Match" : "*",
          "Name" : "kinesis_firehose",
          "Time_key" : "time",
          "delivery_stream" : "${var.log_delivery_stream}",
          "region" : "us-gov-west-1",
          "time_key_format" : "%Y-%m-%dT%H:%M:%S"
        }
      }
    }
  ])

  dynamic "volume" {
    for_each = var.volume_name == null ? [] : [1]
    content {
      name = var.volume_name
      dynamic "efs_volume_configuration" {
        for_each = var.efs_vol_config == null ? [] : [1]
        content {
          file_system_id     = var.efs_vol_config.file_system_id
          root_directory     = var.efs_vol_config.root_directory
          transit_encryption = var.efs_vol_config.transit_encryption
        }
      }
    }
  }
  tags = var.task_def_tags
}