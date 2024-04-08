provider "aws" {
  region = "us-gov-east-1"
}

data "aws_caller_identity" "current" {}

locals {
  account_id = data.aws_caller_identity.current.account_id
}

output "account_selected" {
  value = "account: ${local.account_id}"
}

# ##################################################################
# # LOG GROUP
# ##################################################################
module "cont_log_group" {
  source           = "../log_group"
  create_log_group = true
  log_group_name   = var.log_group_name
  log_group_tags   = var.log_group_tags
}

# output "cont_log_group_output" {
#   value = module.cont_log_group
# }

##################################################################
# TASK DEFINITION
##################################################################
module "cont_task_definitions" {
  source                   = "../task_definition"
  aws_region               = var.aws_region
  task_def_family          = var.task_def_family
  task_role_arn            = var.task_role_arn
  task_exec_role_arn       = var.task_exec_role_arn
  task_def_cpu             = var.task_def_cpu
  task_def_memory          = var.task_def_memory
  app_container_name       = var.app_container_name
  app_container_image      = var.app_container_image
  cont_health_check_config = var.cont_health_check_config
  app_port_mappings        = var.app_port_mappings
  env_variables            = var.env_variables
  log_group                = module.cont_log_group.log_group_name
  log_stream_prefix        = var.log_stream_prefix
  task_def_tags            = var.task_def_tags
}

output "cont_task_definitions_output" {
  value = module.cont_task_definitions
}

##################################################################
# TARGET GROUP
##################################################################
module "cont_target_group" {
  source                   = "../target_group"
  vpc_id                   = var.default_vpc_id
  target_group_name        = var.target_group_name
  tg_health_check_path     = var.tg_health_check_path
  target_group_protocol    = var.target_group_protocol
  target_group_port        = var.target_group_port
  tg_health_check_protocol = var.tg_health_check_protocol
  tg_health_check_port     = var.tg_health_check_port
  target_group_tags        = var.target_group_tags
}

output "cont_target_group_output" {
  value = module.cont_target_group
}

##################################################################
# ADD LISTENER RULE
##################################################################
module "cont_listener_rule" {
  source              = "../listener_rule"
  target_group_arn    = module.cont_target_group.cont_target_group_arn
  listener_arn        = var.listener_arn
  listener_conditions = var.listener_conditions
}

output "cont_cont_listener_rule" {
  value = module.cont_listener_rule
}

##################################################################
# ECS SERVICE
##################################################################
module "cont_ecs_service" {
  source               = "../ecs"
  ecs_service_name     = var.ecs_service_name
  default_vpc_id       = var.default_vpc_id
  ecs_sg_name          = var.ecs_sg_name
  assign_public_ip     = var.assign_public_ip
  ecs_sg_ingress_rules = var.ecs_sg_ingress_rules
  ecs_sg_egress_rules  = var.ecs_sg_egress_rules
  cluster_id           = var.cluster_id
  task_defn_arn        = module.cont_task_definitions.task_definition_arn
  desired_count        = var.desired_count
  alb_security_group   = var.alb_security_group
  ecs_subnet_ids       = var.subnet_ids
  target_group_arn     = module.cont_target_group.cont_target_group_arn
  container_name       = module.cont_task_definitions.td_web_container_name
  container_port       = var.container_port
  depends_on           = [module.cont_listener_rule.alb-listener-rule]
}

output "cont_cont_ecs_service" {
  value = module.cont_ecs_service
}
