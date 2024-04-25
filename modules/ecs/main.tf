# module "ecs_service_sg" {
#   source           = "../security_group"
#   create_sg        = true
#   sg_name          = var.ecs_sg_name
#   sg_description   = "Security Group for ${var.ecs_sg_name}"
#   default_vpc_id   = var.default_vpc_id
#   sg_ingress_rules = var.ecs_sg_ingress_rules
#   sg_egress_rules  = var.ecs_sg_egress_rules
#   sg_tags          = { Name = "${var.ecs_sg_name}" }
# }

#ecs_sg_ingress_rules = [{ from_port = 8443, to_port = 8443, protocol = "tcp", security_groups = ["sg-0c3d633ae57ce2328"], description = "Allow inbound from 8443" },
#{ from_port = 8080, to_port = 8080, protocol = "tcp", security_groups = ["sg-0c3d633ae57ce2328"], description = "Allow inbound from 8080" }]


module "ecs_service_sg" {
  source           = "../security_group"
  create_sg        = true
  sg_name          = var.ecs_sg_name
  sg_description   = "Security Group for ${var.ecs_sg_name}"
  default_vpc_id   = var.default_vpc_id
  sg_ingress_rules = var.ecs_sg_ingress_rules
  sg_egress_rules  = var.ecs_sg_egress_rules
  sg_tags          = { Name = "${var.ecs_sg_name}" }
}

resource "aws_ecs_service" "alb_ecs_service" {
  name                               = var.ecs_service_name
  cluster                            = var.cluster_id
  task_definition                    = var.task_defn_arn
  desired_count                      = var.desired_count
  deployment_minimum_healthy_percent = var.min_healthy_percent
  deployment_maximum_percent         = var.max_healthy_percent
  launch_type                        = var.launch_type
  scheduling_strategy                = var.service_type

  network_configuration {
    security_groups  = [module.ecs_service_sg.security_group_id]
    subnets          = var.ecs_subnet_ids
    assign_public_ip = var.assign_public_ip
  }

  load_balancer {
    target_group_arn = var.target_group_arn
    container_name   = var.container_name
    container_port   = var.container_port
  }

}